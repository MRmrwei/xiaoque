package router

import (
	"github.com/cloudwego/hertz/pkg/route"
	"suki/xqyun/middleware"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/handler/backend"
)

// AdminRouters 注册管理后台接口。
// 前缀：/api/v1/admin
func AdminRouters(
	r *route.RouterGroup,
	app application.Application,
) {
	login := backend.NewLoginController(app)
	refresh := backend.NewRefreshController(app)
	logout := backend.NewLogoutController(app)
	{
		// 登录、刷新令牌与退出登录都不需要 Bearer Token：
		// Access Token 过期后正是靠刷新令牌续期，退出时也可能已无有效 Access Token。
		r.POST("/auth/login", login.Login)
		r.POST("/auth/refresh", refresh.Refresh)
		r.POST("/auth/logout", logout.Logout)
	}

	// 需要 Bearer Token 的后台接口统一注册到 authorized 子组。
	authorized := r.Group("")
	authorized.Use(middleware.NewAdminAuth(app).Handle)
	// TODO: authorized.GET("/me", ...) 等后台接口。
}
