package media

import (
	"context"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
	"suki/xqyun/app/dto"
	mediaservice "suki/xqyun/app/services/media"
)

// List 查询素材列表。
// GET /api/v1/admin/media?kind=video&page=1&page_size=30
func (m *Controller) List(ctx context.Context, c *app.RequestContext) {
	page, pageSize := 1, 30
	if raw := strings.TrimSpace(c.Query("page")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			response.Fail(c, response.CodeBadRequest, "分页参数不正确")
			return
		}
		page = value
	}
	if raw := strings.TrimSpace(c.Query("page_size")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			response.Fail(c, response.CodeBadRequest, "分页参数不正确")
			return
		}
		if value > 100 {
			value = 100
		}
		pageSize = value
	}

	service, err := newService(m.App())
	if err != nil {
		fail(ctx, c, err)
		return
	}

	items, total, err := service.List(ctx, mediaservice.ListInput{
		Kind:     c.Query("kind"),
		Keyword:  c.Query("keyword"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}

	result := make([]dto.AdminMediaItem, 0, len(items))
	for index := range items {
		result = append(result, toMediaItem(&items[index]))
	}
	response.Success(c, dto.AdminMediaListResponse{
		Items:    result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
