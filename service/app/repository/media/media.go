// Package media 提供媒体素材（当前为视频）的数据访问与本地文件存储。
package media

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"suki/xqyun/app/models"
)

// ErrNotFound 表示素材不存在或已删除。
var ErrNotFound = errors.New("素材不存在")

// ListFilter 是素材列表查询条件。
type ListFilter struct {
	Kind     string
	Keyword  string
	Page     int
	PageSize int
}

// Repository 是媒体素材的 GORM 数据访问对象。
type Repository struct {
	db *gorm.DB
}

// New 创建媒体素材数据访问对象。
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create 写入一条素材记录。
func (r *Repository) Create(ctx context.Context, item *models.Media) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("保存素材记录失败：%w", err)
	}
	return nil
}

// FindByID 按主键查询未删除的素材。
func (r *Repository) FindByID(ctx context.Context, id uint64) (*models.Media, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var item models.Media
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at = 0", id).
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询素材失败：%w", err)
	}
	return &item, nil
}

// UpdateOriginalName 修改素材显示名称。
func (r *Repository) UpdateOriginalName(
	ctx context.Context,
	id uint64,
	name string,
) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Media{}).
		Where("id = ? AND deleted_at = 0", id).
		Update("original_name", name)
	if result.Error != nil {
		return fmt.Errorf("修改素材名称失败：%w", result.Error)
	}
	if result.RowsAffected == 0 {
		// 名称未变化时 MySQL 也可能返回 0，这里再用一次查询区分「不存在」。
		if _, err := r.FindByID(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// SoftDelete 软删除素材；文件保留，避免已引用该 URL 的页面失效。
func (r *Repository) SoftDelete(ctx context.Context, id uint64, at int64) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Media{}).
		Where("id = ? AND deleted_at = 0", id).
		Update("deleted_at", at)
	if result.Error != nil {
		return fmt.Errorf("删除素材失败：%w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// SoftDeleteMany 批量软删除素材；不存在的 ID 会被忽略。
func (r *Repository) SoftDeleteMany(ctx context.Context, ids []uint64, at int64) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}
	if len(ids) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Media{}).
		Where("id IN ? AND deleted_at = 0", ids).
		Update("deleted_at", at).Error; err != nil {
		return fmt.Errorf("批量删除素材失败：%w", err)
	}
	return nil
}

// List 按类型 / 关键词分页查询素材，按上传时间倒序。
func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.Media, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("数据访问层：数据库未初始化")
	}

	query := r.db.WithContext(ctx).Model(&models.Media{}).Where("deleted_at = 0")
	if kind := strings.TrimSpace(filter.Kind); kind != "" {
		query = query.Where("kind = ?", kind)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		query = query.Where("original_name LIKE ?", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询素材总数失败：%w", err)
	}

	var items []models.Media
	err := query.
		Order("created_at DESC, id DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&items).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询素材列表失败：%w", err)
	}
	return items, total, nil
}
