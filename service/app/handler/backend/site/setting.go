package site

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
	"suki/xqyun/app/core/response"
)

// SettingController 处理「站点编排 → 站点配置」相关接口。
// 路由前缀：/api/v1/admin/settings
type SettingController struct {
	*controller.BaseController
}

// NewSettingController 创建站点配置控制器。
func NewSettingController(app application.Application) *SettingController {
	return &SettingController{BaseController: controller.New(app)}
}

// Get 读取站点设置。
// GET /api/v1/admin/settings?locale=zh-CN，全局项与语言级设置合并返回。
func (s *SettingController) Get(ctx context.Context, c *app.RequestContext) {
	result, err := newSettingService(s.App()).Get(ctx, c.Query("locale"))
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// Save 保存站点设置。
// PUT /api/v1/admin/settings，请求体为 {locale, key: value, ...}；
// 品牌、收录开关等全局项不区分语言，其余按 locale 保存。
func (s *SettingController) Save(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	var values map[string]json.RawMessage
	if err := json.Unmarshal(body, &values); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	locale := ""
	if raw, ok := values["locale"]; ok {
		if err := json.Unmarshal(raw, &locale); err != nil {
			response.Fail(c, response.CodeBadRequest, "语言参数不正确")
			return
		}
		delete(values, "locale")
	}

	result, err := newSettingService(s.App()).Save(
		ctx,
		locale,
		values,
		adminID(c),
		time.Now().Unix(),
	)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}
