package logger

import (
	"io"
	"strings"
)

// ANSI 颜色码；仅在开发环境终端输出时使用。
const (
	colorReset  = "\x1b[0m"
	colorTrace  = "\x1b[90m"
	colorDebug  = "\x1b[36m"
	colorInfo   = "\x1b[32m"
	colorNotice = "\x1b[34m"
	colorWarn   = "\x1b[33m"
	colorError  = "\x1b[1;31m"
	colorFatal  = "\x1b[1;35m"
)

var levelColors = []struct {
	label string
	color string
}{
	{label: "[Trace]", color: colorTrace},
	{label: "[Debug]", color: colorDebug},
	{label: "[Info]", color: colorInfo},
	{label: "[Notice]", color: colorNotice},
	{label: "[Warn]", color: colorWarn},
	{label: "[Error]", color: colorError},
	{label: "[Fatal]", color: colorFatal},
}

// colorWriter 给 hlog 输出的等级标签加 ANSI 颜色。
// 只用于开发环境终端；生产文件日志和重定向输出不会经过它。
type colorWriter struct {
	writer io.Writer
}

func (w *colorWriter) Write(p []byte) (int, error) {
	line := colorize(string(p))
	if _, err := io.WriteString(w.writer, line); err != nil {
		return 0, err
	}
	// 转换后长度会变化，但按 io.Writer 约定返回原始输入长度。
	return len(p), nil
}

func colorize(line string) string {
	for _, item := range levelColors {
		index := strings.Index(line, item.label)
		if index < 0 {
			continue
		}
		return line[:index] + item.color + item.label + colorReset + line[index+len(item.label):]
	}
	return line
}
