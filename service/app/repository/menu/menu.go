// Package menu 提供顶部 / 底部导航菜单的数据访问。
package menu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"suki/xqyun/app/models"
)

// ErrNotFound 表示菜单容器不存在。
var ErrNotFound = errors.New("菜单不存在")

// Item 是保存菜单项时的树形输入。
type Item struct {
	Label     json.RawMessage
	LinkType  string
	LinkValue *string
	Target    string
	Sort      int
	Visible   int8
	Children  []Item
}

// Repository 是菜单与菜单项的 GORM 数据访问对象。
type Repository struct {
	db *gorm.DB
}

// New 创建菜单数据访问对象。
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 查询全部菜单容器。
func (r *Repository) List(ctx context.Context) ([]models.Menu, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var menus []models.Menu
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("查询菜单列表失败：%w", err)
	}
	return menus, nil
}

// FindByCode 按编码查询菜单容器。
func (r *Repository) FindByCode(ctx context.Context, code string) (*models.Menu, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var menu models.Menu
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询菜单失败：%w", err)
	}
	return &menu, nil
}

// FindItems 查询菜单下的全部菜单项。
func (r *Repository) FindItems(ctx context.Context, menuID uint64) ([]models.MenuItem, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var items []models.MenuItem
	err := r.db.WithContext(ctx).
		Where("menu_id = ?", menuID).
		Order("sort ASC, id ASC").
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("查询菜单项失败：%w", err)
	}
	return items, nil
}

// ReplaceItems 整体替换菜单项：删除旧项后按树形顺序写入新项。
// 菜单容器不存在时会自动创建；items 为空表示清空菜单。
//
// 整体替换会先删除再插入，多个菜单并发保存时可能触发 InnoDB 死锁，
// 这里对死锁 / 锁等待超时做最多 3 次重试。
func (r *Repository) ReplaceItems(
	ctx context.Context,
	code string,
	name string,
	items []Item,
	now int64,
) (*models.Menu, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		menu, err := r.replaceItems(ctx, code, name, items, now)
		if err == nil {
			return menu, nil
		}
		lastErr = err
		if !isRetryableLockError(err) {
			return nil, err
		}
		time.Sleep(time.Duration(attempt+1) * 10 * time.Millisecond)
	}
	return nil, lastErr
}

func (r *Repository) replaceItems(
	ctx context.Context,
	code string,
	name string,
	items []Item,
	now int64,
) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("code = ?", code).First(&menu).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			menu = models.Menu{
				Code:      code,
				Name:      name,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := tx.Create(&menu).Error; err != nil {
				return fmt.Errorf("创建菜单失败：%w", err)
			}
		case err != nil:
			return fmt.Errorf("查询菜单失败：%w", err)
		default:
			updates := map[string]any{"updated_at": now}
			if name != "" {
				updates["name"] = name
				menu.Name = name
			}
			menu.UpdatedAt = now
			if err := tx.Model(&models.Menu{}).
				Where("id = ?", menu.ID).
				Updates(updates).Error; err != nil {
				return fmt.Errorf("更新菜单失败：%w", err)
			}
		}

		if err := tx.Where("menu_id = ?", menu.ID).Delete(&models.MenuItem{}).Error; err != nil {
			return fmt.Errorf("清理菜单项失败：%w", err)
		}

		for index := range items {
			parent := models.MenuItem{
				MenuID:    menu.ID,
				Label:     items[index].Label,
				LinkType:  items[index].LinkType,
				LinkValue: items[index].LinkValue,
				Target:    items[index].Target,
				Sort:      items[index].Sort,
				Visible:   items[index].Visible,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := tx.Create(&parent).Error; err != nil {
				return fmt.Errorf("写入菜单项失败：%w", err)
			}

			for childIndex := range items[index].Children {
				child := models.MenuItem{
					MenuID:    menu.ID,
					ParentID:  &parent.ID,
					Label:     items[index].Children[childIndex].Label,
					LinkType:  items[index].Children[childIndex].LinkType,
					LinkValue: items[index].Children[childIndex].LinkValue,
					Target:    items[index].Children[childIndex].Target,
					Sort:      items[index].Children[childIndex].Sort,
					Visible:   items[index].Children[childIndex].Visible,
					CreatedAt: now,
					UpdatedAt: now,
				}
				if err := tx.Create(&child).Error; err != nil {
					return fmt.Errorf("写入子菜单项失败：%w", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// isRetryableLockError 判断是否为 MySQL 死锁（1213）或锁等待超时（1205）。
func isRetryableLockError(err error) bool {
	var mysqlError *mysqldriver.MySQLError
	if !errors.As(err, &mysqlError) {
		return false
	}
	return mysqlError.Number == 1213 || mysqlError.Number == 1205
}
