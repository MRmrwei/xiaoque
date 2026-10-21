package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"suki/xqyun/app/core/application"
)

// NewRedis 初始化 Redis 连接。
//
// TODO: 按需补充哨兵 / 集群模式、TLS 与键前缀等配置。
func NewRedis(ctx context.Context, cfg application.RedisConfig) (*redis.Client, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	if cfg.Addr == "" {
		return nil, errors.New("application: redis addr is empty")
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("application: ping redis: %w", err)
	}

	return client, nil
}
