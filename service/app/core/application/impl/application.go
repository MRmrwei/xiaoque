package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"suki/xqyun/app/core/application"
)

// Application 是 application.Application 的默认实现。
type Application struct {
	db          *gorm.DB
	redis       *redis.Client
	jwt         application.JWTConfig
	renderCache application.RenderCacheConfig
	upload      application.UploadConfig
}

var _ application.Application = (*Application)(nil)

// New 创建 Application，并完成 DB / Redis 的初始化。
// Redis 初始化失败时会先关闭已经建立的 DB 连接，避免资源泄漏。
func New(ctx context.Context, cfg application.Config) (*Application, error) {
	db, err := NewDB(ctx, cfg.DB)
	if err != nil {
		return nil, err
	}

	rdb, err := NewRedis(ctx, cfg.Redis)
	if err != nil {
		if sqlDB, dbErr := db.DB(); dbErr == nil {
			_ = sqlDB.Close()
		}
		return nil, err
	}

	return &Application{
		db:          db,
		redis:       rdb,
		jwt:         cfg.JWT,
		renderCache: cfg.RenderCache,
		upload:      cfg.Upload,
	}, nil
}

// DB 返回 MySQL 连接。
func (a *Application) DB() *gorm.DB {
	return a.db
}

// Redis 返回 Redis 客户端；未启用 Redis 时为 nil。
func (a *Application) Redis() *redis.Client {
	return a.redis
}

// JWT 返回令牌签名配置；复制 map 与密钥，避免调用方修改内部状态。
func (a *Application) JWT() application.JWTConfig {
	config := a.jwt
	config.Clients = make(map[string]application.JWTClientConfig, len(a.jwt.Clients))
	for name, client := range a.jwt.Clients {
		client.Secret = append([]byte(nil), client.Secret...)
		config.Clients[name] = client
	}
	return config
}

// RenderCache 返回前台渲染缓存配置。
func (a *Application) RenderCache() application.RenderCacheConfig {
	return a.renderCache
}

// Upload 返回媒体上传与访问配置。
func (a *Application) Upload() application.UploadConfig {
	return a.upload
}

// Close 释放 DB 与 Redis 连接，并合并返回所有关闭错误。
func (a *Application) Close() error {
	var errs []error

	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil {
			errs = append(errs, fmt.Errorf("application: get sql db: %w", err))
		} else if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("application: close db: %w", err))
		}
		a.db = nil
	}

	if a.redis != nil {
		if err := a.redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("application: close redis: %w", err))
		}
		a.redis = nil
	}

	return errors.Join(errs...)
}
