package submission

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/response"
	"suki/xqyun/app/dto"
	submissionservice "suki/xqyun/app/services/submission"
)

// Submit 提交表单留言。
// POST /api/v1/frontend/forms/:key/submissions
func (s *Controller) Submit(ctx context.Context, c *app.RequestContext) {
	var request dto.FrontendSubmissionRequest
	if err := c.Bind(&request); err != nil {
		response.Fail(c, response.CodeBadRequest, "提交内容格式不正确")
		return
	}

	result, err := newService(s.App()).Submit(ctx, submissionservice.SubmitInput{
		FormKey:   c.Param("key"),
		Locale:    request.Locale,
		Name:      request.Name,
		Email:     request.Email,
		Company:   request.Company,
		Phone:     request.Phone,
		Message:   request.Message,
		IP:        c.ClientIP(),
		UserAgent: string(c.UserAgent()),
		Now:       time.Now().Unix(),
	})
	if err != nil {
		fail(ctx, c, err)
		return
	}
	response.Success(c, result)
}
