// Package impl 提供 media.Service 的默认实现：校验素材、写入存储并落库。
package impl

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"suki/xqyun/app/models"
	mediarepository "suki/xqyun/app/repository/media"
	"suki/xqyun/app/services/media"
)

// allowedVideoExtensions 是允许上传的视频扩展名。
var allowedVideoExtensions = map[string]struct{}{
	".mp4":  {},
	".webm": {},
	".ogv":  {},
	".ogg":  {},
	".mov":  {},
	".m4v":  {},
}

// allowedImageExtensions 是允许上传的图片扩展名。
// 暂不接收 SVG：同源直开会执行脚本，后续需要时再做清洗。
var allowedImageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".webp": {},
	".gif":  {},
	".avif": {},
}

// Service 是 media.Service 的默认实现。
type Service struct {
	repo          *mediarepository.Repository
	storage       mediarepository.Storage
	maxVideoBytes int64
	maxImageBytes int64
	now           func() time.Time
}

var _ media.Service = (*Service)(nil)

// New 创建默认上传实现。
// storage 通过接口注入，后续可切换本地磁盘或对象存储实现。
func New(
	repo *mediarepository.Repository,
	storage mediarepository.Storage,
	maxVideoBytes int64,
	maxImageBytes int64,
) *Service {
	return &Service{
		repo:          repo,
		storage:       storage,
		maxVideoBytes: maxVideoBytes,
		maxImageBytes: maxImageBytes,
		now:           time.Now,
	}
}

// Upload 校验素材类型与大小，保存文件并写入素材记录。
// Kind 为空时按扩展名推断；显式传入时必须是 video / image 且与扩展名匹配。
func (s *Service) Upload(
	ctx context.Context,
	input media.UploadInput,
	operatorID *uint64,
) (*models.Media, error) {
	originalName := filepath.Base(strings.TrimSpace(input.OriginalName))
	if originalName == "" || originalName == "." || input.Size <= 0 {
		return nil, media.ErrEmptyFile
	}
	if len(originalName) > 255 {
		return nil, media.ErrOriginalNameTooLong
	}

	ext := strings.ToLower(filepath.Ext(originalName))
	kind, ok := detectKind(ext)
	if !ok {
		return nil, media.ErrUnsupportedFile
	}
	if requested := strings.TrimSpace(input.Kind); requested != "" && requested != kind {
		return nil, media.ErrUnsupportedFile
	}

	maxBytes := s.maxImageBytes
	expectedMimePrefix := "image/"
	if kind == media.KindVideo {
		maxBytes = s.maxVideoBytes
		expectedMimePrefix = "video/"
	}
	if maxBytes > 0 && input.Size > maxBytes {
		return nil, media.ErrFileTooLarge
	}

	mime := strings.ToLower(strings.TrimSpace(input.Mime))
	if mime != "" && mime != "application/octet-stream" &&
		!strings.HasPrefix(mime, expectedMimePrefix) {
		return nil, media.ErrUnsupportedFile
	}

	now := s.now()
	filename, url, err := s.storage.Save(input.Reader, ext, now)
	if err != nil {
		return nil, err
	}

	item := &models.Media{
		Kind:         kind,
		OriginalName: originalName,
		Filename:     filename,
		URL:          url,
		Mime:         mime,
		Size:         uint64(input.Size),
		CreatedBy:    operatorID,
		CreatedAt:    now.Unix(),
	}
	if err := s.repo.Create(ctx, item); err != nil {
		// 数据库写入失败时删除已落盘的文件，避免产生孤儿文件。
		_ = s.storage.Remove(filename)
		return nil, err
	}
	return item, nil
}

// List 分页查询素材；Kind 为空表示全部类型。
func (s *Service) List(
	ctx context.Context,
	input media.ListInput,
) ([]models.Media, int64, error) {
	kind := strings.TrimSpace(input.Kind)
	if kind != "" && kind != media.KindVideo && kind != media.KindImage {
		return nil, 0, media.ErrUnsupportedKind
	}

	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 30
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return s.repo.List(ctx, mediarepository.ListFilter{
		Kind:     kind,
		Keyword:  input.Keyword,
		Page:     page,
		PageSize: pageSize,
	})
}

// Detail 查询素材详情。
func (s *Service) Detail(ctx context.Context, id uint64) (*models.Media, error) {
	item, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, mediarepository.ErrNotFound) {
		return nil, media.ErrMediaNotFound
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Update 修改素材显示名称。
func (s *Service) Update(
	ctx context.Context,
	id uint64,
	input media.UpdateInput,
) (*models.Media, error) {
	name := strings.TrimSpace(input.OriginalName)
	if name == "" {
		return nil, media.ErrOriginalNameRequired
	}
	if len(name) > 255 {
		return nil, media.ErrOriginalNameTooLong
	}

	if err := s.repo.UpdateOriginalName(ctx, id, name); err != nil {
		if errors.Is(err, mediarepository.ErrNotFound) {
			return nil, media.ErrMediaNotFound
		}
		return nil, err
	}
	return s.Detail(ctx, id)
}

// Delete 软删除素材；物理文件保留，避免已引用该 URL 的页面失效。
func (s *Service) Delete(ctx context.Context, id uint64) error {
	err := s.repo.SoftDelete(ctx, id, s.now().Unix())
	if errors.Is(err, mediarepository.ErrNotFound) {
		return media.ErrMediaNotFound
	}
	return err
}

// DeleteMany 批量软删除素材；忽略不存在的 ID。
func (s *Service) DeleteMany(ctx context.Context, ids []uint64) error {
	cleaned := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id > 0 {
			cleaned = append(cleaned, id)
		}
	}
	if len(cleaned) == 0 {
		return nil
	}
	return s.repo.SoftDeleteMany(ctx, cleaned, s.now().Unix())
}

// detectKind 根据扩展名判断素材类型，为后续扩展其他类型保留统一入口。
func detectKind(ext string) (string, bool) {
	if _, ok := allowedVideoExtensions[ext]; ok {
		return media.KindVideo, true
	}
	if _, ok := allowedImageExtensions[ext]; ok {
		return media.KindImage, true
	}
	return "", false
}
