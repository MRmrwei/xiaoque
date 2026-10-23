// Package media 提供管理后台媒体素材（视频 / 图片）的 HTTP 控制器。
//
// 控制器是「消费者」：只依赖 app/services/media 里的接口，
// 具体用哪个实现由 app/utils 的装配函数决定；每个方法内按请求创建用例。
//
// 文件划分：
//   - upload.go 上传素材
//   - list.go   素材列表
package media

import (
	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
)

// Controller 处理管理后台素材相关接口。
// 路由前缀：/api/v1/admin/media
type Controller struct {
	*controller.BaseController
}

// NewController 创建素材控制器。
func NewController(app application.Application) *Controller {
	return &Controller{BaseController: controller.New(app)}
}
