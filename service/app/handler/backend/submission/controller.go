package submission

import (
	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
)

// Controller 处理「线索 → 表单留言」相关接口。
type Controller struct {
	*controller.BaseController
}

// NewController 创建表单留言控制器。
func NewController(app application.Application) *Controller {
	return &Controller{BaseController: controller.New(app)}
}
