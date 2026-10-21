// Package logger 初始化 Hertz hlog 的输出方式。
//
// 开发环境输出到控制台；生产环境写入日志目录，并使用 lumberjack 按大小轮转。
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/natefinch/lumberjack.v2"

	"suki/xqyun/app/core/config"
)

// Init 根据 APP_ENV 和日志配置初始化 hlog。
// 返回的关闭函数用于进程退出时关闭日志文件；开发环境返回空操作。
func Init(cfg *config.Config) (func() error, error) {
	if cfg == nil {
		return nil, fmt.Errorf("logger: config is nil")
	}

	hlog.SetLevel(parseLevel(cfg.Log.Level))

	if !isProduction(cfg.App.Env) {
		var output io.Writer = os.Stdout
		if cfg.Log.Color && os.Getenv("NO_COLOR") == "" && isTerminal(os.Stdout) {
			output = &colorWriter{writer: os.Stdout}
		}
		hlog.SetOutput(output)
		return func() error { return nil }, nil
	}

	dir := strings.TrimSpace(cfg.Log.Dir)
	if dir == "" {
		dir = "logs"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("logger: create log dir %s: %w", dir, err)
	}

	filename := strings.TrimSpace(cfg.Log.Filename)
	if filename == "" {
		filename = "app.log"
	}

	writer := &lumberjack.Logger{
		Filename:   filepath.Join(dir, filename),
		MaxSize:    cfg.Log.MaxSizeMB,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAgeDays,
		Compress:   cfg.Log.Compress,
		LocalTime:  true,
	}
	hlog.SetOutput(writer)

	return writer.Close, nil
}

// isTerminal 判断 stdout 是否连接到终端；重定向到文件或管道时不加颜色。
func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func isProduction(env string) bool {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "production", "prod", "release":
		return true
	default:
		return false
	}
}

func parseLevel(level string) hlog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "notice":
		return hlog.LevelNotice
	case "warn", "warning":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}
