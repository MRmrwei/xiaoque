package submission

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
)

// List 留言列表。
// GET /api/v1/admin/submissions
// 支持 status / form_key / locale / keyword / created_from / created_to / 分页。
func (s *Controller) List(ctx context.Context, c *app.RequestContext) {
	input, ok := listInput(c)
	if !ok {
		return
	}

	result, err := newService(s.App()).List(ctx, input)
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}
