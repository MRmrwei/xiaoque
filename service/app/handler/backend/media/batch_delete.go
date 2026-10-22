package media

import (
	"context"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
)

// BatchDelete 批量软删除素材。
// DELETE /api/v1/admin/media?ids=1,2,3
func (m *Controller) BatchDelete(ctx context.Context, c *app.RequestContext) {
	raw := strings.TrimSpace(c.Query("ids"))
	if raw == "" {
		response.Fail(c, response.CodeBadRequest, "请选择要删除的素材")
		return
	}

	parts := strings.Split(raw, ",")
	ids := make([]uint64, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil || id == 0 {
			response.Fail(c, response.CodeBadRequest, "素材 ID 无效")
			return
		}
		ids = append(ids, id)
	}

	service, err := newService(m.App())
	if err != nil {
		fail(ctx, c, err)
		return
	}

	if err := service.DeleteMany(ctx, ids); err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, nil)
}
