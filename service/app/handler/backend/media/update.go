package media

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
	"suki/xqyun/app/dto"
	mediaservice "suki/xqyun/app/services/media"
)

// Update 修改素材信息（当前支持显示名称）。
// PUT /api/v1/admin/media/:id
func (m *Controller) Update(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "素材 ID 无效")
	if !ok {
		return
	}

	var request dto.AdminMediaUpdateRequest
	if err := c.Bind(&request); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	service, err := newService(m.App())
	if err != nil {
		fail(ctx, c, err)
		return
	}

	item, err := service.Update(ctx, id, mediaservice.UpdateInput{
		OriginalName: request.OriginalName,
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, toMediaItem(item))
}
