package logger

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"suki/xqyun/app/core/config"
)

func TestIsProduction(t *testing.T) {
	production := []string{"production", "prod", "release", " Production "}
	for _, env := range production {
		if !isProduction(env) {
			t.Fatalf("%q 应识别为生产环境", env)
		}
	}

	for _, env := range []string{"", "dev", "development", "test", "staging"} {
		if isProduction(env) {
			t.Fatalf("%q 不应识别为生产环境", env)
		}
	}
}

func TestParseLevel(t *testing.T) {
	cases := map[string]hlog.Level{
		"trace":   hlog.LevelTrace,
		"debug":   hlog.LevelDebug,
		"info":    hlog.LevelInfo,
		"notice":  hlog.LevelNotice,
		"warn":    hlog.LevelWarn,
		"warning": hlog.LevelWarn,
		"error":   hlog.LevelError,
		"fatal":   hlog.LevelFatal,
		"unknown": hlog.LevelInfo,
		"":        hlog.LevelInfo,
	}

	for value, want := range cases {
		if got := parseLevel(value); got != want {
			t.Fatalf("parseLevel(%q) = %v，期望 %v", value, got, want)
		}
	}
}

func TestInitWritesFileInProduction(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		App: config.AppConfig{Env: "production"},
		Log: config.LogConfig{
			Level:      "info",
			Dir:        dir,
			Filename:   "app.log",
			MaxSizeMB:  1,
			MaxBackups: 2,
			MaxAgeDays: 1,
			Compress:   false,
		},
	}

	closeLogger, err := Init(cfg)
	if err != nil {
		t.Fatalf("Init() 返回错误：%v", err)
	}
	defer func() {
		_ = closeLogger()
		hlog.SetOutput(os.Stdout)
	}()

	hlog.CtxErrorf(context.Background(), "logger test message")

	data, err := os.ReadFile(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatalf("读取日志文件失败：%v", err)
	}
	if !strings.Contains(string(data), "logger test message") {
		t.Fatalf("日志文件未写入预期内容：%q", string(data))
	}
}

func TestInitRejectsNilConfig(t *testing.T) {
	if _, err := Init(nil); err == nil {
		t.Fatal("Init(nil) 应返回错误")
	}
}

func TestColorizeLevelPrefix(t *testing.T) {
	line := "2026/10/21 13:48:00.000000 login.go:73: [Error] admin login failed\n"
	got := colorize(line)
	if !strings.Contains(got, colorError+"[Error]"+colorReset) {
		t.Fatalf("Error 等级未正确着色：%q", got)
	}

	plain := "message without level"
	if colorize(plain) != plain {
		t.Fatalf("无等级标签的日志不应被修改：%q", colorize(plain))
	}
}

func TestColorWriterReturnsOriginalLength(t *testing.T) {
	var buffer strings.Builder
	writer := &colorWriter{writer: &buffer}
	input := []byte("[Warn] warning message\n")

	written, err := writer.Write(input)
	if err != nil {
		t.Fatalf("Write() 返回错误：%v", err)
	}
	if written != len(input) {
		t.Fatalf("Write() 返回 %d，期望 %d", written, len(input))
	}
	if !strings.Contains(buffer.String(), colorWarn+"[Warn]"+colorReset) {
		t.Fatalf("输出未包含 Warn 颜色：%q", buffer.String())
	}
}

func TestIsTerminalOnRegularFile(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "log-*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败：%v", err)
	}
	defer file.Close()

	if isTerminal(file) {
		t.Fatal("普通文件不应被识别为终端")
	}
}
