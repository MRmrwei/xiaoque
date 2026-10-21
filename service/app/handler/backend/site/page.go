package site

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
	"suki/xqyun/app/core/response"
	"suki/xqyun/app/dto"
	siteservice "suki/xqyun/app/services/site"
)

// PageController 处理「站点编排 → 页面内容」相关接口。
// 路由前缀：/api/v1/admin/pages
type PageController struct {
	*controller.BaseController
}

// NewPageController 创建页面内容控制器。
func NewPageController(app application.Application) *PageController {
	return &PageController{BaseController: controller.New(app)}
}

// List 页面列表。
// GET /api/v1/admin/pages，支持 locale / status / keyword / 分页。
func (p *PageController) List(ctx context.Context, c *app.RequestContext) {
	page, pageSize, ok := pagination(c)
	if !ok {
		return
	}

	result, err := newPageService(p.App()).List(ctx, siteservice.PageListInput{
		Locale:   c.Query("locale"),
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// Create 新建页面。
// POST /api/v1/admin/pages，CopyFromID 非空时复制来源页面结构。
func (p *PageController) Create(ctx context.Context, c *app.RequestContext) {
	var request dto.AdminPageCreateRequest
	if err := c.Bind(&request); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	result, err := newPageService(p.App()).Create(ctx, siteservice.PageCreateInput{
		Locale:     request.Locale,
		Path:       request.Path,
		Name:       request.Name,
		SEO:        request.SEO,
		Blocks:     request.Blocks,
		CopyFromID: request.CopyFromID,
		OperatorID: adminID(c),
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// Detail 页面详情（当前草稿）。
// GET /api/v1/admin/pages/:id
func (p *PageController) Detail(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "页面 ID 无效")
	if !ok {
		return
	}

	result, err := newPageService(p.App()).Detail(ctx, id)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// SaveDraft 保存页面草稿。
// PUT /api/v1/admin/pages/:id，只写草稿，不影响线上内容。
func (p *PageController) SaveDraft(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "页面 ID 无效")
	if !ok {
		return
	}

	var request dto.AdminPageSaveDraftRequest
	if err := c.Bind(&request); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	result, err := newPageService(p.App()).SaveDraft(ctx, id, siteservice.PageDraftInput{
		Path:       request.Path,
		Name:       request.Name,
		SEO:        request.SEO,
		Blocks:     request.Blocks,
		OperatorID: adminID(c),
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// Delete 软删除页面。
// DELETE /api/v1/admin/pages/:id
func (p *PageController) Delete(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "页面 ID 无效")
	if !ok {
		return
	}

	if err := newPageService(p.App()).Delete(ctx, id); err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, nil)
}

// Translate 创建 / 获取另一语言配对草稿。
// POST /api/v1/admin/pages/:id/translate
func (p *PageController) Translate(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "页面 ID 无效")
	if !ok {
		return
	}

	var request dto.AdminPageTranslateRequest
	if err := c.Bind(&request); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	result, err := newPageService(p.App()).Translate(ctx, id, siteservice.PageTranslateInput{
		FromLocale: request.FromLocale,
		ToLocale:   request.ToLocale,
		Path:       request.Path,
		Name:       request.Name,
		SEO:        request.SEO,
		Blocks:     request.Blocks,
		OperatorID: adminID(c),
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// Publish 发布当前草稿，覆盖线上内容。
// POST /api/v1/admin/pages/:id/publish
func (p *PageController) Publish(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "页面 ID 无效")
	if !ok {
		return
	}

	var request dto.AdminPagePublishRequest
	body, err := c.Body()
	if err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}
	if len(bytes.TrimSpace(body)) > 0 {
		if err := json.Unmarshal(body, &request); err != nil {
			response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
			return
		}
	}

	result, err := newPageService(p.App()).Publish(ctx, id, siteservice.PagePublishInput{
		Note:       request.Note,
		OperatorID: adminID(c),
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}

// Unpublish 下架页面：状态改为未发布，清空线上内容，草稿保留。
// POST /api/v1/admin/pages/:id/unpublish
func (p *PageController) Unpublish(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "页面 ID 无效")
	if !ok {
		return
	}

	var request dto.AdminPagePublishRequest
	body, err := c.Body()
	if err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}
	if len(bytes.TrimSpace(body)) > 0 {
		if err := json.Unmarshal(body, &request); err != nil {
			response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
			return
		}
	}

	err = newPageService(p.App()).Unpublish(ctx, id, siteservice.PagePublishInput{
		Note:       request.Note,
		OperatorID: adminID(c),
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, nil)
}

// CreatePreviewToken 生成草稿预览令牌。
// POST /api/v1/admin/pages/:id/preview-token
func (p *PageController) CreatePreviewToken(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "页面 ID 无效")
	if !ok {
		return
	}

	result, err := newPageService(p.App()).CreatePreviewToken(ctx, id)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}
