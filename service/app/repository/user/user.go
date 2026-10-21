// Package user 提供后台账号的数据访问。
package user

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"suki/xqyun/app/models"
)

// ErrNotFound 表示按条件没有找到用户。
var ErrNotFound = errors.New("用户不存在")

// UserRepository 定义后台账号的数据访问能力。
// 只负责用户资料本身；登录会话与刷新令牌见 app/repository/session。
type UserRepository interface {
	// FindByUsername 按登录名查询未删除的用户。
	FindByUsername(ctx context.Context, username string) (*models.User, error)

	// FindByID 按主键查询未删除的用户。
	FindByID(ctx context.Context, id uint64) (*models.User, error)

	// UpdateLoginState 更新用户最近登录时间与 IP。
	// 登录会话存放在 Redis，这里只维护 MySQL 中的登录信息。
	UpdateLoginState(ctx context.Context, userID uint64, lastLoginAt int64, ip *string) error
}

// Repository 是 UserRepository 的 GORM 实现。
type Repository struct {
	db *gorm.DB
}

var _ UserRepository = (*Repository)(nil)

// New 创建后台账号数据访问对象。
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindByUsername 按登录名查询未删除的用户。
func (r *Repository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var user models.User
	err := r.db.WithContext(ctx).
		Where("username = ? AND deleted_at = 0", username).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户失败：%w", err)
	}
	return &user, nil
}

// FindByID 按主键查询未删除的用户。
func (r *Repository) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	var user models.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at = 0", id).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户失败：%w", err)
	}
	return &user, nil
}

// FindNamesByIDs 批量查询用户显示名称，供页面列表展示最后修改人。
func (r *Repository) FindNamesByIDs(
	ctx context.Context,
	ids []uint64,
) (map[uint64]string, error) {
	if r.db == nil {
		return nil, fmt.Errorf("数据访问层：数据库未初始化")
	}

	names := make(map[uint64]string, len(ids))
	if len(ids) == 0 {
		return names, nil
	}

	var users []models.User
	err := r.db.WithContext(ctx).
		Select("id", "name").
		Where("id IN ? AND deleted_at = 0", ids).
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("查询用户名称失败：%w", err)
	}
	for _, user := range users {
		names[user.ID] = user.Name
	}
	return names, nil
}

// UpdateLoginState 更新用户最近登录时间与 IP。
func (r *Repository) UpdateLoginState(
	ctx context.Context,
	userID uint64,
	lastLoginAt int64,
	ip *string,
) error {
	if r.db == nil {
		return fmt.Errorf("数据访问层：数据库未初始化")
	}

	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND deleted_at = 0", userID).
		Updates(map[string]any{
			"last_login_at": lastLoginAt,
			"last_login_ip": ip,
			"updated_at":    lastLoginAt,
		})
	if result.Error != nil {
		return fmt.Errorf("更新最近登录信息失败：%w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
