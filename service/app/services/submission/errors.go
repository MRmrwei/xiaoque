// Package submission 提供后台表单留言的业务用例。
package submission

import "suki/xqyun/app/core/response"

// Error 是可直接映射为统一响应结构的业务错误。
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	// ErrNotFound 表示留言不存在。
	ErrNotFound = &Error{Code: response.CodeNotFound, Message: "留言不存在"}
)

// badRequest 构造参数错误。
func badRequest(message string) *Error {
	return &Error{Code: response.CodeBadRequest, Message: message}
}
