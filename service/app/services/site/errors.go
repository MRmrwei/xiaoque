// Package site 提供管理后台站点编排（页面、导航菜单、站点设置）的业务用例。
package site

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
	// ErrPageNotFound 表示页面不存在或已删除。
	ErrPageNotFound = &Error{Code: response.CodeNotFound, Message: "页面不存在"}
	// ErrMenuNotFound 表示菜单容器不存在。
	ErrMenuNotFound = &Error{Code: response.CodeNotFound, Message: "菜单不存在"}
	// ErrPathConflict 表示同一语言下路径重复。
	ErrPathConflict = &Error{
		Code:    response.CodeConflict,
		Message: "该语言下已存在相同路径的页面",
	}
)

// badRequest 构造参数错误。
func badRequest(message string) *Error {
	return &Error{Code: response.CodeBadRequest, Message: message}
}
