// Package setting 提供全站设置（站点配置）的数据访问。
package setting

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"suki/xqyun/app/models"
)

// Repository 是全站设置的 GORM 数据访问对象。
type Repository struct {
	db *gorm.DB
}

// New 创建全站设置数据访问对象。
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 查询指定语言的设置，同时返回全局设置（locale 为空字符串）。
// 全局项排在前面，同 key 时语言级设置由上层覆盖。
func (r *Repository) List(ctx context.Context, locale string) ([]models.Setting, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	query := r.db.WithContext(ctx).Model(&models.Setting{})
	if locale == "" {
		query = query.Where("locale = ''")
	} else {
		query = query.Where("locale = ? OR locale = ''", locale)
	}

	var settings []models.Setting
	if err := query.Order("locale ASC, id ASC").Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("查询站点设置失败：%w", err)
	}
	return settings, nil
}

// UpsertMany 按 key + locale 写入设置，已存在时覆盖 value。
func (r *Repository) UpsertMany(
	ctx context.Context,
	locale string,
	values map[string]json.RawMessage,
	operatorID *uint64,
	now int64,
) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}
	if len(values) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for key, value := range values {
			setting := models.Setting{
				Key:       key,
				Locale:    locale,
				Value:     value,
				UpdatedBy: operatorID,
				UpdatedAt: now,
			}
			err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "key"},
					{Name: "locale"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"value",
					"updated_by",
					"updated_at",
				}),
			}).Create(&setting).Error
			if err != nil {
				return fmt.Errorf("保存设置 %s 失败：%w", key, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
