package submission

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/core/response"
	submissionservice "suki/xqyun/app/services/submission"
	"suki/xqyun/app/core/application"
	submissionrepository "suki/xqyun/app/repository/submission"
	submissionimpl "suki/xqyun/app/services/submission/impl"
)

// newService 按请求组装表单留言用例。
func newService(app application.Application) *submissionimpl.Service {
	return submissionimpl.New(submissionrepository.New(app.DB()))
}

// fail 把用例错误转换为统一响应；未知错误返回 500 并记录日志。
func fail(ctx context.Context, c *app.RequestContext, err error) {
	var businessError *submissionservice.Error
	if errors.As(err, &businessError) {
		response.Fail(c, businessError.Code, businessError.Message)
		return
	}

	hlog.CtxErrorf(ctx, "admin submission: %v", err)
	response.Fail(c, response.CodeInternal, "操作失败，请稍后重试")
}

// pathID 解析路径参数中的数字 ID，失败时直接写入参数错误响应。
func pathID(c *app.RequestContext) (uint64, bool) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, response.CodeBadRequest, "留言 ID 无效")
		return 0, false
	}
	return id, true
}

// listInput 读取列表 / 导出共用的筛选与分页参数。
func listInput(c *app.RequestContext) (submissionservice.ListInput, bool) {
	createdFrom, ok := queryInt64(c, "created_from")
	if !ok {
		return submissionservice.ListInput{}, false
	}
	createdTo, ok := queryInt64(c, "created_to")
	if !ok {
		return submissionservice.ListInput{}, false
	}
	page, ok := queryInt(c, "page")
	if !ok {
		return submissionservice.ListInput{}, false
	}
	pageSize, ok := queryInt(c, "page_size")
	if !ok {
		return submissionservice.ListInput{}, false
	}

	return submissionservice.ListInput{
		Status:      c.Query("status"),
		FormKey:     c.Query("form_key"),
		Locale:      c.Query("locale"),
		Keyword:     c.Query("keyword"),
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		Page:        page,
		PageSize:    pageSize,
	}, true
}

func queryInt64(c *app.RequestContext, key string) (int64, bool) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0, true
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		response.Fail(c, response.CodeBadRequest, "时间筛选参数不正确")
		return 0, false
	}
	return value, true
}

func queryInt(c *app.RequestContext, key string) (int, bool) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		response.Fail(c, response.CodeBadRequest, "分页参数不正确")
		return 0, false
	}
	return value, true
}
