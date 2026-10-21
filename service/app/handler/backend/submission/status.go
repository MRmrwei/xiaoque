package submission

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
	"suki/xqyun/app/dto"
)

// UpdateStatus 标记留言处理状态。
// PUT /api/v1/admin/submissions/:id/status
func (s *Controller) UpdateStatus(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	var request dto.AdminSubmissionStatusRequest
	if err := c.Bind(&request); err != nil {
		response.Fail(c, response.CodeBadRequest, "请求参数格式不正确")
		return
	}

	result, err := newService(s.App()).UpdateStatus(
		ctx,
		id,
		request.Status,
		time.Now().Unix(),
	)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}
