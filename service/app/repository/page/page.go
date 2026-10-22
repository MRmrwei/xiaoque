// Package page 提供页面（站点编排 → 页面内容）的数据访问。
package page

import (
	"context"
	"errors"
	"fmt"
	"strings"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"suki/xqyun/app/models"
)

var (
	// ErrNotFound 表示页面不存在或已删除。
	ErrNotFound = errors.New("页面不存在")
	// ErrPathConflict 表示同一语言下路径已被占用（唯一索引冲突）。
	ErrPathConflict = errors.New("页面路径已存在")
)

// ListFilter 是页面列表查询条件。
type ListFilter struct {
	Locale   string
	Status   string
	Keyword  string
	Page     int
	PageSize int
}

// Repository 是页面的 GORM 数据访问对象。
type Repository struct {
	db *gorm.DB
}

// New 创建页面数据访问对象。
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 按语言、状态、关键词分页查询未删除页面，按更新时间倒序。
func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.Page, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("数据访问层：数据库未初始化")
	}

	query := r.db.WithContext(ctx).Model(&models.Page{}).Where("deleted_at = 0")
	if filter.Locale != "" {
		query = query.Where("locale = ?", filter.Locale)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(name LIKE ? OR path LIKE ?)", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询页面总数失败：%w", err)
	}

	var pages []models.Page
	err := query.
		Order("updated_at DESC, id DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&pages).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询页面列表失败：%w", err)
	}
	return pages, total, nil
}

// FindByID 按主键查询未删除页面。
func (r *Repository) FindByID(ctx context.Context, id uint64) (*models.Page, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var page models.Page
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at = 0", id).
		First(&page).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询页面失败：%w", err)
	}
	return &page, nil
}

// FindByPath 按语言与路径查询未删除页面，不区分发布状态。
// 同一语言下路径唯一，查询结果最多一条。
func (r *Repository) FindByPath(ctx context.Context, locale string, path string) (*models.Page, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var page models.Page
	err := r.db.WithContext(ctx).
		Where("locale = ? AND path = ? AND deleted_at = 0", locale, path).
		First(&page).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询页面失败：%w", err)
	}
	return &page, nil
}

// FindByGroupAndLocale 按翻译组与语言查询页面，用于判断是否已有配对版本。
func (r *Repository) FindByGroupAndLocale(
	ctx context.Context,
	groupID string,
	locale string,
) (*models.Page, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var page models.Page
	err := r.db.WithContext(ctx).
		Where("translation_group_id = ? AND locale = ? AND deleted_at = 0", groupID, locale).
		First(&page).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询配对页面失败：%w", err)
	}
	return &page, nil
}

// FindByGroupIDs 批量查询多个翻译组下的页面。
func (r *Repository) FindByGroupIDs(
	ctx context.Context,
	groupIDs []string,
) ([]models.Page, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}
	if len(groupIDs) == 0 {
		return nil, nil
	}

	var pages []models.Page
	err := r.db.WithContext(ctx).
		Where("translation_group_id IN ? AND deleted_at = 0", groupIDs).
		Find(&pages).Error
	if err != nil {
		return nil, fmt.Errorf("查询翻译组页面失败：%w", err)
	}
	return pages, nil
}

// ExistsPath 判断同一语言下路径是否已被其他未删除页面占用。
func (r *Repository) ExistsPath(
	ctx context.Context,
	locale string,
	path string,
	excludeID uint64,
) (bool, error) {
	if r.db == nil {
		return false, fmt.Errorf("数据访问层：数据库未初始化")
	}

	query := r.db.WithContext(ctx).
		Model(&models.Page{}).
		Where("locale = ? AND path = ? AND deleted_at = 0", locale, path)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查页面路径失败：%w", err)
	}
	return count > 0, nil
}

// Create 新建页面。
func (r *Repository) Create(ctx context.Context, page *models.Page) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}
	if err := r.db.WithContext(ctx).Create(page).Error; err != nil {
		if isDuplicateKey(err) {
			return ErrPathConflict
		}
		return fmt.Errorf("创建页面失败：%w", err)
	}
	return nil
}

// UpdateDraft 保存页面路径、名称与草稿内容，不修改线上内容。
func (r *Repository) UpdateDraft(ctx context.Context, page *models.Page) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Page{}).
		Where("id = ? AND deleted_at = 0", page.ID).
		Updates(map[string]any{
			"path":          page.Path,
			"name":          page.Name,
			"draft_content": page.DraftContent,
			"updated_by":    page.UpdatedBy,
			"updated_at":    page.UpdatedAt,
		})
	if result.Error != nil {
		if isDuplicateKey(result.Error) {
			return ErrPathConflict
		}
		return fmt.Errorf("保存页面草稿失败：%w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// SoftDelete 软删除页面。
func (r *Repository) SoftDelete(ctx context.Context, id uint64, at int64) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Page{}).
		Where("id = ? AND deleted_at = 0", id).
		Updates(map[string]any{"deleted_at": at, "updated_at": at})
	if result.Error != nil {
		return fmt.Errorf("删除页面失败：%w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// isDuplicateKey 判断是否为 MySQL 唯一索引冲突（错误码 1062）。
func isDuplicateKey(err error) bool {
	var mysqlError *mysqldriver.MySQLError
	return errors.As(err, &mysqlError) && mysqlError.Number == 1062
}

// Publish 用当前草稿覆盖线上内容，版本号 +1，并写入发布记录。
// 草稿与发布记录在同一个事务内完成。
func (r *Repository) Publish(
	ctx context.Context,
	id uint64,
	operatorID *uint64,
	note *string,
	now int64,
) (int, error) {
	if r.db == nil {
		return 0, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var version int
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var page models.Page
		err := tx.Where("id = ? AND deleted_at = 0", id).First(&page).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("查询页面失败：%w", err)
		}

		version = int(page.Version) + 1
		result := tx.Model(&models.Page{}).
			Where("id = ?", page.ID).
			Updates(map[string]any{
				"published_content": page.DraftContent,
				"status":            "published",
				"version":           version,
				"updated_by":        operatorID,
				"updated_at":        now,
			})
		if result.Error != nil {
			return fmt.Errorf("发布页面失败：%w", result.Error)
		}

		record := models.PublishRecord{
			TargetType: "page",
			TargetID:   page.ID,
			Locale:     page.Locale,
			Status:     "success",
			Message:    note,
			OperatorID: operatorID,
			CreatedAt:  now,
		}
		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("写入发布记录失败：%w", err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return version, nil
}

// Unpublish 下架页面：清空线上内容并回到未发布，草稿保留。
func (r *Repository) Unpublish(
	ctx context.Context,
	id uint64,
	operatorID *uint64,
	note *string,
	now int64,
) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var page models.Page
		err := tx.Where("id = ? AND deleted_at = 0", id).First(&page).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("查询页面失败：%w", err)
		}

		result := tx.Model(&models.Page{}).
			Where("id = ?", page.ID).
			Updates(map[string]any{
				"published_content": nil,
				"status":            "draft",
				"updated_by":        operatorID,
				"updated_at":        now,
			})
		if result.Error != nil {
			return fmt.Errorf("下架页面失败：%w", result.Error)
		}

		record := models.PublishRecord{
			TargetType: "page",
			TargetID:   page.ID,
			Locale:     page.Locale,
			Status:     "success",
			Message:    note,
			OperatorID: operatorID,
			CreatedAt:  now,
		}
		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("写入下架记录失败：%w", err)
		}
		return nil
	})
}
