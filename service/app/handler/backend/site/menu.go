package site

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
	"suki/xqyun/app/core/response"
	"suki/xqyun/app/dto"
)

// MenuController 处理「站点编排 → 顶部导航 / 底部导航」相关接口。
// 路由前缀：/api/v1/admin/menus
//
// 菜单用 code 区分容器：header（顶部导航）、footer（底部导航）、footer_bottom（页脚底部链接）。
type MenuController struct {
	*controller.BaseController
}

// NewMenuController 创建导航菜单控制器。
func NewMenuController(app application.Application) *MenuController {
	return &MenuController{BaseController: controller.New(app)}
}

// List 菜单列表。
// GET /api/v1/admin/menus
func (m *MenuController) List(ctx context.Context, c *app.RequestContext) {
	result, err := newMenuService(m.App()).List(ctx)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// Detail 菜单详情与菜单项。
// GET /api/v1/admin/menus/:code
func (m *MenuController) Detail(ctx context.Context, c *app.RequestContext) {
	result, err := newMenuService(m.App()).Detail(ctx, c.Param("code"))
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// Save 整体保存菜单项（排序、显隐、文案、链接）。
// PUT /api/v1/admin/menus/:code，菜单容器不存在时自动创建。
func (m *MenuController) Save(ctx context.Context, c *app.RequestContext) {
	var request dto.AdminMenuSaveRequest
	if err := c.Bind(&request); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	result, err := newMenuService(m.App()).Save(
		ctx,
		c.Param("code"),
		request,
		time.Now().Unix(),
	)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}
