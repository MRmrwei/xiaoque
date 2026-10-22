package site

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/dto"
	"suki/xqyun/app/models"
	pagerepository "suki/xqyun/app/repository/page"
	renderrepository "suki/xqyun/app/repository/render"
	userrepository "suki/xqyun/app/repository/user"
	"suki/xqyun/app/utils/ulid"
)

const (
	// pageStatusDraft 未发布。
	pageStatusDraft = "draft"
	// pageStatusPublished 已发布。
	pageStatusPublished = "published"
)

// PageListInput 是页面列表查询条件。
type PageListInput struct {
	Locale   string
	Status   string
	Keyword  string
	Page     int
	PageSize int
}

// PageCreateInput 是新建页面入参。
type PageCreateInput struct {
	Locale     string
	Path       string
	Name       string
	SEO        *dto.SEO
	Blocks     []dto.PageBlock
	CopyFromID string
	OperatorID *uint64
}

// PageDraftInput 是保存草稿入参。路径 / 名称为空时保留原值，SEO / Blocks 为 nil 时保留原内容。
type PageDraftInput struct {
	Path       string
	Name       string
	SEO        *dto.SEO
	Blocks     []dto.PageBlock
	OperatorID *uint64
}

// PageTranslateInput 是创建翻译配对入参。
type PageTranslateInput struct {
	FromLocale string
	ToLocale   string
	Path       string
	Name       string
	SEO        *dto.SEO
	Blocks     []dto.PageBlock
	OperatorID *uint64
}

// PagePublishInput 是发布 / 下架入参。
type PagePublishInput struct {
	Note       string
	OperatorID *uint64
}

// PageService 是页面内容用例服务。
type PageService struct {
	pages *pagerepository.Repository
	users *userrepository.Repository
	cache *renderrepository.Cache
}

// NewPageService 创建页面内容用例服务。
// cache 只在前台渲染缓存启用时生效，发布 / 下架 / 删除后主动失效。
func NewPageService(
	pages *pagerepository.Repository,
	users *userrepository.Repository,
	cache *renderrepository.Cache,
) *PageService {
	return &PageService{
		pages: pages,
		users: users,
		cache: cache,
	}
}

// List 查询页面列表，并补充翻译配对与最后修改人名称。
func (s *PageService) List(ctx context.Context, input PageListInput) (dto.AdminPageListResponse, error) {
	locale := ""
	if strings.TrimSpace(input.Locale) != "" {
		normalized, err := normalizeLocale(input.Locale, true)
		if err != nil {
			return dto.AdminPageListResponse{}, err
		}
		locale = normalized
	}

	status := strings.TrimSpace(input.Status)
	if status != "" && status != pageStatusDraft && status != pageStatusPublished {
		return dto.AdminPageListResponse{}, badRequest("页面状态参数不正确")
	}

	page, pageSize := input.Page, input.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	pages, total, err := s.pages.List(ctx, pagerepository.ListFilter{
		Locale:   locale,
		Status:   status,
		Keyword:  input.Keyword,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return dto.AdminPageListResponse{}, err
	}

	groupIDs := make([]string, 0, len(pages))
	groupSeen := make(map[string]struct{}, len(pages))
	operatorIDs := make([]uint64, 0, len(pages))
	operatorSeen := make(map[uint64]struct{}, len(pages))
	for _, item := range pages {
		if _, ok := groupSeen[item.TranslationGroupID]; !ok {
			groupSeen[item.TranslationGroupID] = struct{}{}
			groupIDs = append(groupIDs, item.TranslationGroupID)
		}
		if item.UpdatedBy != nil {
			if _, ok := operatorSeen[*item.UpdatedBy]; !ok {
				operatorSeen[*item.UpdatedBy] = struct{}{}
				operatorIDs = append(operatorIDs, *item.UpdatedBy)
			}
		}
	}

	groupPages, err := s.pages.FindByGroupIDs(ctx, groupIDs)
	if err != nil {
		return dto.AdminPageListResponse{}, err
	}
	groupMap := make(map[string][]models.Page, len(groupIDs))
	for _, item := range groupPages {
		groupMap[item.TranslationGroupID] = append(
			groupMap[item.TranslationGroupID],
			item,
		)
	}

	operatorNames, err := s.users.FindNamesByIDs(ctx, operatorIDs)
	if err != nil {
		return dto.AdminPageListResponse{}, err
	}

	items := make([]dto.AdminPageSummary, 0, len(pages))
	for _, item := range pages {
		updatedBy := ""
		if item.UpdatedBy != nil {
			updatedBy = operatorNames[*item.UpdatedBy]
		}

		translations := make([]dto.AdminPageTranslation, 0, 1)
		for _, sibling := range groupMap[item.TranslationGroupID] {
			if sibling.ID == item.ID {
				continue
			}
			translations = append(translations, dto.AdminPageTranslation{
				PageID: strconv.FormatUint(sibling.ID, 10),
				Locale: sibling.Locale,
				Path:   sibling.Path,
				Name:   sibling.Name,
				Status: sibling.Status,
			})
		}

		items = append(items, dto.AdminPageSummary{
			ID:           strconv.FormatUint(item.ID, 10),
			GroupID:      item.TranslationGroupID,
			Locale:       item.Locale,
			Path:         item.Path,
			Name:         item.Name,
			Status:       item.Status,
			Version:      int(item.Version),
			UpdatedAt:    item.UpdatedAt,
			UpdatedBy:    updatedBy,
			Translations: translations,
		})
	}

	return dto.AdminPageListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Detail 查询页面详情（当前草稿）。
func (s *PageService) Detail(ctx context.Context, id uint64) (dto.AdminPageDraft, error) {
	page, err := s.findPage(ctx, id)
	if err != nil {
		return dto.AdminPageDraft{}, err
	}
	return toPageDraft(page)
}

// Create 新建页面，可复制来源页面的草稿结构。
func (s *PageService) Create(ctx context.Context, input PageCreateInput) (dto.AdminPageDraft, error) {
	locale, err := normalizeLocale(input.Locale, true)
	if err != nil {
		return dto.AdminPageDraft{}, err
	}
	path, err := validatePath(input.Path)
	if err != nil {
		return dto.AdminPageDraft{}, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return dto.AdminPageDraft{}, badRequest("请填写页面名称")
	}
	if len(name) > 200 {
		return dto.AdminPageDraft{}, badRequest("页面名称不能超过 200 个字符")
	}

	content := emptyPageContent(name)
	if copyFromID := strings.TrimSpace(input.CopyFromID); copyFromID != "" {
		sourceID, parseErr := strconv.ParseUint(copyFromID, 10, 64)
		if parseErr != nil {
			return dto.AdminPageDraft{}, badRequest("复制来源页面 ID 无效")
		}
		source, findErr := s.findPage(ctx, sourceID)
		if findErr != nil {
			return dto.AdminPageDraft{}, findErr
		}
		sourceContent, decodeErr := pageContentOf(source.DraftContent)
		if decodeErr != nil {
			return dto.AdminPageDraft{}, decodeErr
		}
		content.SEO = sourceContent.SEO
		content.Blocks = regenerateBlockIDs(sourceContent.Blocks)
	}
	if input.SEO != nil {
		content.SEO = *input.SEO
	}
	if input.Blocks != nil {
		content.Blocks = input.Blocks
	}
	normalizeSEO(&content.SEO)

	if err := s.ensurePathAvailable(ctx, locale, path, 0); err != nil {
		return dto.AdminPageDraft{}, err
	}

	groupID, err := ulid.New()
	if err != nil {
		return dto.AdminPageDraft{}, err
	}
	rawContent, err := json.Marshal(content)
	if err != nil {
		return dto.AdminPageDraft{}, fmt.Errorf("编码页面内容失败：%w", err)
	}

	now := time.Now().Unix()
	page := models.Page{
		TranslationGroupID: groupID,
		Locale:             locale,
		Path:               path,
		Name:               name,
		Status:             pageStatusDraft,
		DraftContent:       rawContent,
		CreatedBy:          input.OperatorID,
		UpdatedBy:          input.OperatorID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := s.pages.Create(ctx, &page); err != nil {
		return dto.AdminPageDraft{}, pageRepositoryError(err)
	}
	return toPageDraft(&page)
}

// SaveDraft 保存页面草稿，只写草稿字段，不影响线上内容。
func (s *PageService) SaveDraft(
	ctx context.Context,
	id uint64,
	input PageDraftInput,
) (dto.AdminPageDraft, error) {
	page, err := s.findPage(ctx, id)
	if err != nil {
		return dto.AdminPageDraft{}, err
	}

	content, err := pageContentOf(page.DraftContent)
	if err != nil {
		return dto.AdminPageDraft{}, err
	}

	oldPath := page.Path
	path := page.Path
	if strings.TrimSpace(input.Path) != "" {
		path, err = validatePath(input.Path)
		if err != nil {
			return dto.AdminPageDraft{}, err
		}
	}
	name := page.Name
	if trimmed := strings.TrimSpace(input.Name); trimmed != "" {
		if len(trimmed) > 200 {
			return dto.AdminPageDraft{}, badRequest("页面名称不能超过 200 个字符")
		}
		name = trimmed
	}
	if input.SEO != nil {
		content.SEO = *input.SEO
	}
	if input.Blocks != nil {
		content.Blocks = input.Blocks
	}
	normalizeSEO(&content.SEO)

	if path != page.Path {
		if err := s.ensurePathAvailable(ctx, page.Locale, path, page.ID); err != nil {
			return dto.AdminPageDraft{}, err
		}
	}

	rawContent, err := json.Marshal(content)
	if err != nil {
		return dto.AdminPageDraft{}, fmt.Errorf("编码页面内容失败：%w", err)
	}

	page.Path = path
	page.Name = name
	page.DraftContent = rawContent
	page.UpdatedBy = input.OperatorID
	page.UpdatedAt = time.Now().Unix()
	if err := s.pages.UpdateDraft(ctx, page); err != nil {
		return dto.AdminPageDraft{}, pageRepositoryError(err)
	}

	// 草稿内容不影响线上，只有路径变化会改变前台地址与 alternates，需要失效缓存。
	if path != oldPath {
		s.invalidateGroupCache(ctx, page)
		s.invalidatePageCache(ctx, &models.Page{Locale: page.Locale, Path: oldPath})
	}
	return toPageDraft(page)
}

// Delete 软删除页面。
func (s *PageService) Delete(ctx context.Context, id uint64) error {
	page, err := s.findPage(ctx, id)
	if err != nil {
		return err
	}
	if err := s.pages.SoftDelete(ctx, id, time.Now().Unix()); err != nil {
		return pageRepositoryError(err)
	}
	s.invalidateGroupCache(ctx, page)
	return nil
}

// Translate 创建或获取另一语言配对草稿。
func (s *PageService) Translate(
	ctx context.Context,
	id uint64,
	input PageTranslateInput,
) (dto.AdminPageDraft, error) {
	source, err := s.findPage(ctx, id)
	if err != nil {
		return dto.AdminPageDraft{}, err
	}

	if strings.TrimSpace(input.FromLocale) != "" {
		fromLocale, localeErr := normalizeLocale(input.FromLocale, false)
		if localeErr != nil {
			return dto.AdminPageDraft{}, localeErr
		}
		if fromLocale != source.Locale {
			return dto.AdminPageDraft{}, badRequest("源语言与页面语言不一致")
		}
	}

	toLocale := otherLocale(source.Locale)
	if strings.TrimSpace(input.ToLocale) != "" {
		toLocale, err = normalizeLocale(input.ToLocale, false)
		if err != nil {
			return dto.AdminPageDraft{}, err
		}
	}
	if toLocale == source.Locale {
		return dto.AdminPageDraft{}, badRequest("目标语言不能与源语言相同")
	}

	existing, findErr := s.pages.FindByGroupAndLocale(ctx, source.TranslationGroupID, toLocale)
	if findErr == nil {
		return toPageDraft(existing)
	}
	if !errors.Is(findErr, pagerepository.ErrNotFound) {
		return dto.AdminPageDraft{}, findErr
	}

	path, err := validatePath(input.Path)
	if err != nil {
		return dto.AdminPageDraft{}, err
	}
	if err := s.ensurePathAvailable(ctx, toLocale, path, 0); err != nil {
		return dto.AdminPageDraft{}, err
	}

	sourceContent, err := pageContentOf(source.DraftContent)
	if err != nil {
		return dto.AdminPageDraft{}, err
	}

	content := pageContent{
		SEO:    dto.SEO{Keywords: []string{}},
		Blocks: clearTranslatableBlocks(sourceContent.Blocks),
	}
	if input.SEO != nil {
		content.SEO = *input.SEO
	}
	if input.Blocks != nil {
		content.Blocks = input.Blocks
	}
	normalizeSEO(&content.SEO)

	name := strings.TrimSpace(input.Name)
	if name == "" {
		if toLocale == localeEN {
			name = source.Name + " (EN)"
		} else {
			name = source.Name + "（中文）"
		}
	}
	if len(name) > 200 {
		return dto.AdminPageDraft{}, badRequest("页面名称不能超过 200 个字符")
	}

	rawContent, err := json.Marshal(content)
	if err != nil {
		return dto.AdminPageDraft{}, fmt.Errorf("编码页面内容失败：%w", err)
	}

	now := time.Now().Unix()
	translated := models.Page{
		TranslationGroupID: source.TranslationGroupID,
		Locale:             toLocale,
		Path:               path,
		Name:               name,
		Status:             pageStatusDraft,
		DraftContent:       rawContent,
		CreatedBy:          input.OperatorID,
		UpdatedBy:          input.OperatorID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := s.pages.Create(ctx, &translated); err != nil {
		return dto.AdminPageDraft{}, pageRepositoryError(err)
	}
	return toPageDraft(&translated)
}

// Publish 发布当前草稿，覆盖线上内容。
func (s *PageService) Publish(
	ctx context.Context,
	id uint64,
	input PagePublishInput,
) (dto.AdminPagePublishResponse, error) {
	page, err := s.findPage(ctx, id)
	if err != nil {
		return dto.AdminPagePublishResponse{}, err
	}

	now := time.Now().Unix()
	version, err := s.pages.Publish(ctx, id, input.OperatorID, optionalString(input.Note), now)
	if err != nil {
		return dto.AdminPagePublishResponse{}, pageRepositoryError(err)
	}
	s.invalidateGroupCache(ctx, page)

	return dto.AdminPagePublishResponse{
		Version:     version,
		PublishedAt: now,
	}, nil
}

// Unpublish 下架页面：状态回到未发布，清空线上内容，草稿保留。
func (s *PageService) Unpublish(
	ctx context.Context,
	id uint64,
	input PagePublishInput,
) error {
	page, err := s.findPage(ctx, id)
	if err != nil {
		return err
	}

	note := strings.TrimSpace(input.Note)
	if note == "" {
		note = "下架"
	}
	noteValue := note
	if err := s.pages.Unpublish(
		ctx,
		id,
		input.OperatorID,
		&noteValue,
		time.Now().Unix(),
	); err != nil {
		return pageRepositoryError(err)
	}
	s.invalidateGroupCache(ctx, page)
	return nil
}

// invalidateGroupCache 失效页面及其同翻译组所有语言的缓存。
// alternates 只包含已发布语言，发布 / 下架 / 删除会改变同组其他页面的 alt 地址。
func (s *PageService) invalidateGroupCache(ctx context.Context, page *models.Page) {
	targets := []*models.Page{page}
	if page.TranslationGroupID != "" {
		siblings, err := s.pages.FindByGroupIDs(ctx, []string{page.TranslationGroupID})
		if err != nil {
			hlog.CtxWarnf(ctx, "render cache: query translation group: %v", err)
		} else {
			for index := range siblings {
				targets = append(targets, &siblings[index])
			}
		}
	}
	s.invalidatePageCache(ctx, targets...)
}

// invalidatePageCache 删除页面分片；缓存删除失败只记日志，不影响发布结果。
func (s *PageService) invalidatePageCache(ctx context.Context, pages ...*models.Page) {
	if s.cache == nil || !s.cache.Enabled() || len(pages) == 0 {
		return
	}

	keys := make([]renderrepository.PageKey, 0, len(pages))
	for _, page := range pages {
		keys = append(keys, renderrepository.PageKey{
			Locale: page.Locale,
			Path:   page.Path,
		})
	}
	if err := s.cache.InvalidatePage(ctx, keys...); err != nil {
		hlog.CtxWarnf(ctx, "render cache: invalidate page: %v", err)
	}
}

func (s *PageService) findPage(ctx context.Context, id uint64) (*models.Page, error) {
	page, err := s.pages.FindByID(ctx, id)
	if err != nil {
		return nil, pageRepositoryError(err)
	}
	return page, nil
}

func (s *PageService) ensurePathAvailable(
	ctx context.Context,
	locale string,
	path string,
	excludeID uint64,
) error {
	exists, err := s.pages.ExistsPath(ctx, locale, path, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return ErrPathConflict
	}
	return nil
}

func validatePath(raw string) (string, error) {
	path := normalizePath(raw)
	if path == "" {
		return "", badRequest("请填写页面路径")
	}
	if len(path) > 255 {
		return "", badRequest("页面路径不能超过 255 个字符")
	}
	return path, nil
}

func pageRepositoryError(err error) error {
	if errors.Is(err, pagerepository.ErrNotFound) {
		return ErrPageNotFound
	}
	if errors.Is(err, pagerepository.ErrPathConflict) {
		return ErrPathConflict
	}
	return err
}

func optionalString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func toPageDraft(page *models.Page) (dto.AdminPageDraft, error) {
	content, err := pageContentOf(page.DraftContent)
	if err != nil {
		return dto.AdminPageDraft{}, err
	}
	return dto.AdminPageDraft{
		ID:        strconv.FormatUint(page.ID, 10),
		GroupID:   page.TranslationGroupID,
		Locale:    page.Locale,
		Path:      page.Path,
		Name:      page.Name,
		Status:    page.Status,
		Version:   int(page.Version),
		SEO:       content.SEO,
		Blocks:    content.Blocks,
		UpdatedAt: page.UpdatedAt,
	}, nil
}
