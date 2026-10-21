package submission

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
)

// Detail 留言详情。
// GET /api/v1/admin/submissions/:id
func (s *Controller) Detail(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	result, err := newService(s.App()).Detail(ctx, id)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}
