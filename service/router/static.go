package router

import (
	"os"
	"strings"

	hertzapp "github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/core/application"
)

// RegisterUploadStatic 注册本地媒体文件的静态访问。
// 只在 upload.provider 为空或 local 时启用，并支持 Range 以便视频拖动进度。
func RegisterUploadStatic(r *server.Hertz, cfg application.UploadConfig) {
	isLocal := cfg.Provider == "" || strings.EqualFold(cfg.Provider, "local")
	if !isLocal || cfg.Dir == "" {
		return
	}

	if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {
		hlog.Fatalf("create upload dir %s: %v", cfg.Dir, err)
	}

	prefix := strings.TrimRight(cfg.URLPrefix, "/")
	if prefix == "" {
		prefix = "/uploads"
	}

	// Hertz 的 FS 默认把完整请求路径拼到 Root 后面，需要先去掉 URL 前缀，
	// 否则会去找 uploads/uploads/2026/...，导致图片 / 视频全部 404。
	pathRewrite := func(ctx *hertzapp.RequestContext) []byte {
		requestPath := strings.TrimPrefix(string(ctx.Path()), prefix)
		if requestPath == "" {
			requestPath = "/"
		}
		if !strings.HasPrefix(requestPath, "/") {
			requestPath = "/" + requestPath
		}
		return []byte(requestPath)
	}

	r.StaticFS(prefix, &hertzapp.FS{
		Root:            cfg.Dir,
		AcceptByteRange: true,
		PathRewrite:     pathRewrite,
	})
}
