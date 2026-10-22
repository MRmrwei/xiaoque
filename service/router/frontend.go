package router

import (
	"github.com/cloudwego/hertz/pkg/route"

	"suki/xqyun/app/core/application"
	sitehandler "suki/xqyun/app/handler/frontend/site"
	submissionhandler "suki/xqyun/app/handler/frontend/submission"
)

// FrontendRouters 注册官网前台公开接口。
// 前缀：/api/v1/frontend
func FrontendRouters(r *route.RouterGroup, app application.Application) {
	// 页面渲染：一次返回页面、导航、站点配置与内容集合。
	pages := sitehandler.NewPageController(app)
	r.GET("/render/page", pages.RenderPage)

	// 站点级数据：头部 / 页脚单独取配置与导航。
	site := sitehandler.NewSiteController(app)
	r.GET("/site/config", site.Config)
	r.GET("/site/nav", site.Nav)

	// 表单提交：联系表单等公开表单写入后台「表单留言」。
	submissions := submissionhandler.NewController(app)
	r.POST("/forms/:key/submissions", submissions.Submit)
}
