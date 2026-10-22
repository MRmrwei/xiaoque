package media

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
)

// Detail 查询素材详情。
// GET /api/v1/admin/media/:id
func (m *Controller) Detail(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "素材 ID 无效")
	if !ok {
		return
	}

	service, err := newService(m.App())
	if err != nil {
		fail(ctx, c, err)
		return
	}

	item, err := service.Detail(ctx, id)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, toMediaItem(item))
}
