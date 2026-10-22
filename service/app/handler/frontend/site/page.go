package site

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
	"suki/xqyun/app/core/response"
	siteservice "suki/xqyun/app/services/site"
)

// PageController 处理官网前台页面渲染接口。
// 路由前缀：/api/v1/frontend/render
type PageController struct {
	*controller.BaseController
}

// NewPageController 创建前台页面渲染控制器。
func NewPageController(app application.Application) *PageController {
	return &PageController{BaseController: controller.New(app)}
}

// RenderPage 获取页面渲染所需的全部数据。
// GET /api/v1/frontend/render/page?path=/about&locale=zh-CN
//
// 只返回已发布页面；页面不存在或未发布返回 40400，前台据此输出 404。
func (p *PageController) RenderPage(ctx context.Context, c *app.RequestContext) {
	result, err := newRenderer(p.App()).RenderPage(ctx, siteservice.RenderPageInput{
		Path:   c.Query("path"),
		Locale: c.Query("locale"),
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}
