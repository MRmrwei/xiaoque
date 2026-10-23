// Package impl 提供素材用例的默认实现：上传（经 FileUploader）与素材库管理。
package impl

import (
	"context"
	"errors"
	"strings"
	"time"

	"suki/xqyun/app/models"
	mediarepository "suki/xqyun/app/repository/media"
	"suki/xqyun/app/services/media"
)

// Service 是素材用例的默认实现：
// 嵌入 FileUploader 负责上传文件，自身负责写素材记录与管理逻辑。
type Service struct {
	repo *mediarepository.Repository
	// 嵌入纯上传能力：由装配方用 NewFileUploader 创建后注入。
	*FileUploader
	now func() time.Time
}

// 上传文件与落库两种逻辑通过 interface 解耦：
// 上传见 media.FileUploader（由 FileUploader 实现），落库见 media.Uploader（本类型实现）。
var _ media.Uploader = (*Service)(nil)

// New 创建素材管理实现；uploader 由装配方用 NewFileUploader 创建。
func New(repo *mediarepository.Repository, uploader *FileUploader) *Service {
	return &Service{repo: repo, FileUploader: uploader, now: time.Now}
}

// Upload 先用 FileUploader 校验并写入存储，拿到文件信息后再落库；
// 落库失败时删除已写入的文件，避免产生孤儿文件。
func (s *Service) Upload(
	ctx context.Context,
	input media.UploadInput,
	operatorID *uint64,
) (*models.Media, error) {
	stored, err := s.FileUploader.Upload(ctx, input)
	if err != nil {
		return nil, err
	}

	item := &models.Media{
		Kind:         stored.Kind,
		OriginalName: stored.OriginalName,
		Filename:     stored.Filename,
		URL:          stored.URL,
		Mime:         stored.Mime,
		Size:         uint64(stored.Size),
		CreatedBy:    operatorID,
		CreatedAt:    stored.CreatedAt,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		_ = s.FileUploader.Remove(stored.Filename)
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

// Delete 删除素材：先同步删除物理文件，再软删除记录。
// 文件删除失败时不删记录，便于重试；已不存在的文件视为删除成功。
func (s *Service) Delete(ctx context.Context, id uint64) error {
	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mediarepository.ErrNotFound) {
			return media.ErrMediaNotFound
		}
		return err
	}

	if err := s.FileUploader.Remove(item.Filename); err != nil {
		return err
	}

	err = s.repo.SoftDelete(ctx, id, s.now().Unix())
	if errors.Is(err, mediarepository.ErrNotFound) {
		return media.ErrMediaNotFound
	}
	return err
}

// DeleteMany 批量删除素材：逐个同步删除物理文件，再批量软删除记录；忽略不存在的 ID。
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

	items := make([]models.Media, 0, len(cleaned))
	for _, id := range cleaned {
		item, err := s.repo.FindByID(ctx, id)
		if errors.Is(err, mediarepository.ErrNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		items = append(items, *item)
	}

	for index := range items {
		if err := s.FileUploader.Remove(items[index].Filename); err != nil {
			return err
		}
	}

	return s.repo.SoftDeleteMany(ctx, cleaned, s.now().Unix())
}
