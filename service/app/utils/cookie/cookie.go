// Package cookie 提供管理后台 Refresh Token Cookie 的读写工具。
package cookie

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"

	"suki/xqyun/app/core/application"
)

// SetRefresh 写入 Refresh Token 的 HttpOnly Cookie。
func SetRefresh(
	c *app.RequestContext,
	value string,
	maxAge int,
	cfg application.CookieConfig,
) {
	if cfg.Name == "" || value == "" {
		return
	}
	sameSite := parseSameSite(cfg.SameSite)
	secure := cfg.Secure || sameSite == protocol.CookieSameSiteNoneMode
	c.SetCookie(cfg.Name, value, maxAge, cfg.Path, cfg.Domain, sameSite, secure, true)
}

// ClearRefresh 删除 Refresh Token Cookie，用于退出登录。
func ClearRefresh(c *app.RequestContext, cfg application.CookieConfig) {
	if cfg.Name == "" {
		return
	}
	sameSite := parseSameSite(cfg.SameSite)
	secure := cfg.Secure || sameSite == protocol.CookieSameSiteNoneMode
	c.SetCookie(cfg.Name, "", -1, cfg.Path, cfg.Domain, sameSite, secure, true)
}

// RefreshValue 读取 Refresh Token Cookie。
func RefreshValue(c *app.RequestContext, cfg application.CookieConfig) (string, bool) {
	if cfg.Name == "" {
		return "", false
	}
	value := strings.TrimSpace(string(c.Cookie(cfg.Name)))
	return value, value != ""
}

func parseSameSite(value string) protocol.CookieSameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return protocol.CookieSameSiteStrictMode
	case "none":
		return protocol.CookieSameSiteNoneMode
	default:
		return protocol.CookieSameSiteLaxMode
	}
}
