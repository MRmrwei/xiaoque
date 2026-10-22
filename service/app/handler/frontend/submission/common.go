package submission

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/response"
	submissionrepository "suki/xqyun/app/repository/submission"
	submissionservice "suki/xqyun/app/services/submission"
)

// newService 按请求组装表单留言用例服务。
func newService(app application.Application) *submissionservice.Service {
	return submissionservice.New(submissionrepository.New(app.DB()))
}

// fail 把用例错误转换为统一响应；未知错误返回 500 并记录日志。
func fail(ctx context.Context, c *app.RequestContext, err error) {
	var businessError *submissionservice.Error
	if errors.As(err, &businessError) {
		response.Fail(c, businessError.Code, businessError.Message)
		return
	}

	hlog.CtxErrorf(ctx, "frontend submission: %v", err)
	response.Fail(c, response.CodeInternal, "提交失败，请稍后重试")
}
