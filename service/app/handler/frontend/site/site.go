package site

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
	"suki/xqyun/app/core/response"
)

// SiteController 处理官网前台的站点级数据（站点配置、导航）。
// 路由前缀：/api/v1/frontend/site
type SiteController struct {
	*controller.BaseController
}

// NewSiteController 创建前台站点配置控制器。
func NewSiteController(app application.Application) *SiteController {
	return &SiteController{BaseController: controller.New(app)}
}

// Config 读取脱敏后的站点配置，供头部 / 页脚使用。
// GET /api/v1/frontend/site/config?locale=zh-CN
func (s *SiteController) Config(ctx context.Context, c *app.RequestContext) {
	result, err := newSettingService(s.App()).Get(ctx, c.Query("locale"))
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// Nav 读取顶部菜单与页脚栏目。
// GET /api/v1/frontend/site/nav?locale=zh-CN
func (s *SiteController) Nav(ctx context.Context, c *app.RequestContext) {
	result, err := newMenuService(s.App()).Nav(ctx, c.Query("locale"))
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}
