package dto

// FrontendSubmissionRequest 是前台提交表单请求。
// 字段与区块表单一致；email 选填，phone 必填。
type FrontendSubmissionRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
	Locale  string `json:"locale"`
}

// FrontendSubmissionResult 是前台表单提交结果。
type FrontendSubmissionResult struct {
	ID string `json:"id"`
}
