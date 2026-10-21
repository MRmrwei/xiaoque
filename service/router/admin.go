package router

import (
	"github.com/cloudwego/hertz/pkg/route"
	"suki/xqyun/middleware"

	"suki/xqyun/app/core/application"
	sitehandler "suki/xqyun/app/handler/backend/site"
	submissionhandler "suki/xqyun/app/handler/backend/submission"
	userhandler "suki/xqyun/app/handler/backend/user"
)

// AdminRouters 注册管理后台接口。
// 前缀：/api/v1/admin
func AdminRouters(
	r *route.RouterGroup,
	app application.Application,
) {
	login := userhandler.NewLoginController(app)
	refresh := userhandler.NewRefreshController(app)
	logout := userhandler.NewLogoutController(app)
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

	// 站点编排：页面内容
	pages := sitehandler.NewPageController(app)
	authorized.GET("/pages", pages.List)
	authorized.POST("/pages", pages.Create)
	authorized.GET("/pages/:id", pages.Detail)
	authorized.PUT("/pages/:id", pages.SaveDraft)
	authorized.DELETE("/pages/:id", pages.Delete)
	authorized.POST("/pages/:id/translate", pages.Translate)
	authorized.POST("/pages/:id/publish", pages.Publish)
	authorized.POST("/pages/:id/unpublish", pages.Unpublish)
	authorized.POST("/pages/:id/preview-token", pages.CreatePreviewToken)

	// 站点编排：顶部导航 / 底部导航 / 页脚底部链接
	menus := sitehandler.NewMenuController(app)
	authorized.GET("/menus", menus.List)
	authorized.GET("/menus/:code", menus.Detail)
	authorized.PUT("/menus/:code", menus.Save)

	// 站点编排：站点配置
	settings := sitehandler.NewSettingController(app)
	authorized.GET("/settings", settings.Get)
	authorized.PUT("/settings", settings.Save)

	// 线索：表单留言
	submissions := submissionhandler.NewController(app)
	authorized.GET("/submissions", submissions.List)
	// 静态路径放在 :id 之前注册，避免 export 被当成留言 ID。
	authorized.GET("/submissions/export", submissions.Export)
	authorized.GET("/submissions/:id", submissions.Detail)
	authorized.PUT("/submissions/:id/status", submissions.UpdateStatus)
}
