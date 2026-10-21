package bootstrap

import (
	"context"
	"errors"
	"fmt"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/application/impl"
	"suki/xqyun/app/core/config"
	"suki/xqyun/app/core/logger"
)

// Bootstrap 聚合启动阶段创建好的配置、应用实例与可关闭资源。
type Bootstrap struct {
	Config      *config.Config
	App         application.Application
	closeLogger func() error
}

// New 加载配置并创建 Application。
// 配置加载、DB / Redis 初始化失败时会直接返回错误。
func New(ctx context.Context) (*Bootstrap, error) {
	cfg, err := config.Load("")
	if err != nil {
		return nil, fmt.Errorf("bootstrap: load config: %w", err)
	}

	closeLogger, err := logger.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: init logger: %w", err)
	}

	appCfg, err := cfg.ApplicationConfig()
	if err != nil {
		return nil, fmt.Errorf("bootstrap: build application config: %w", err)
	}

	app, err := impl.New(ctx, appCfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: init application: %w", err)
	}

	return &Bootstrap{
		Config:      cfg,
		App:         app,
		closeLogger: closeLogger,
	}, nil
}

// Close 释放 Application 持有的 DB / Redis 连接，并关闭日志文件。
func (b *Bootstrap) Close() error {
	if b == nil {
		return nil
	}

	var errs []error
	if b.App != nil {
		if err := b.App.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if b.closeLogger != nil {
		if err := b.closeLogger(); err != nil {
			errs = append(errs, fmt.Errorf("bootstrap: close logger: %w", err))
		}
	}
	return errors.Join(errs...)
}
