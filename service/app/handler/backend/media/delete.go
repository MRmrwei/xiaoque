package media

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
)

// Delete 软删除素材；文件保留，已引用该 URL 的页面不受影响。
// DELETE /api/v1/admin/media/:id
func (m *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	id, ok := pathID(c, "id", "素材 ID 无效")
	if !ok {
		return
	}

	service, err := newService(m.App())
	if err != nil {
		fail(ctx, c, err)
		return
	}

	if err := service.Delete(ctx, id); err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, nil)
}
