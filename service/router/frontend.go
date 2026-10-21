package router

import (
	"github.com/cloudwego/hertz/pkg/route"
	"suki/xqyun/app/core/application"
)

// FrontendRouters 注册官网前台公开接口。
// 前缀：/api/v1/frontend
func FrontendRouters(r *route.RouterGroup, app application.Application) {
	// TODO: 注册 render/page、forms 等官网公开接口。
}
