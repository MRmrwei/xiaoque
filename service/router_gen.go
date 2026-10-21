package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"suki/xqyun/app/core/application"
)

func register(r *server.Hertz, app application.Application) {
	customizedRegister(r, app)
}
