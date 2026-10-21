package controller

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"suki/xqyun/app/core/application"
)

// BaseController 提供各业务控制器共享的依赖访问。
// 业务控制器通过嵌入 *BaseController 复用，不需要重复持有 Application。
type BaseController struct {
	app application.Application
}

// New 创建 BaseController。
func New(app application.Application) *BaseController {
	return &BaseController{app: app}
}

// App 返回应用基础设施。
func (c *BaseController) App() application.Application {
	return c.app
}

// DB 返回 MySQL 连接。
func (c *BaseController) DB() *gorm.DB {
	return c.app.DB()
}

// Redis 返回 Redis 客户端；未启用 Redis 时为 nil。
func (c *BaseController) Redis() *redis.Client {
	return c.app.Redis()
}
