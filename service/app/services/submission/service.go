// Package submission 定义表单留言的用例契约：接口、入参、状态与错误。
//
// 具体实现见 app/services/submission/impl；控制器只依赖这里的接口。
package submission

import (
	"context"

	"suki/xqyun/app/dto"
)

// 留言处理状态。
const (
	StatusNew        = "new"
	StatusProcessing = "processing"
	StatusDone       = "done"
	StatusSpam       = "spam"
)

// ListInput 是留言列表 / 导出的查询条件。
type ListInput struct {
	Status      string
	FormKey     string
	Locale      string
	Keyword     string
	CreatedFrom int64
	CreatedTo   int64
	Page        int
	PageSize    int
}

// SubmitInput 是前台提交表单的原始内容与请求信息。
type SubmitInput struct {
	FormKey   string
	Locale    string
	Name      string
	Email     string
	Company   string
	Phone     string
	Message   string
	IP        string
	UserAgent string
	Now       int64
}

// Submitter 是前台提交能力。
type Submitter interface {
	Submit(ctx context.Context, input SubmitInput) (dto.FrontendSubmissionResult, error)
}

// Lister 是后台列表与导出能力。
type Lister interface {
	List(ctx context.Context, input ListInput) (dto.AdminSubmissionListResponse, error)
	Export(ctx context.Context, input ListInput) ([]dto.AdminSubmission, error)
}

// Manager 是后台详情与状态管理能力。
type Manager interface {
	Detail(ctx context.Context, id uint64) (dto.AdminSubmission, error)
	UpdateStatus(ctx context.Context, id uint64, status string, now int64) (dto.AdminSubmission, error)
}

// Service 是表单留言的用例集合。
type Service interface {
	Submitter
	Lister
	Manager
}
