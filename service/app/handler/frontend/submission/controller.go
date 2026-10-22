package submission

import (
	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/controller"
)

// Controller 处理官网前台的公开表单提交接口。
type Controller struct {
	*controller.BaseController
}

// NewController 创建前台表单控制器。
func NewController(app application.Application) *Controller {
	return &Controller{BaseController: controller.New(app)}
}
