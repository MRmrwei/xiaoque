// Package middleware 提供 Hertz 中间件。
package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"suki/xqyun/app/core/application"
	"suki/xqyun/app/core/response"
	"suki/xqyun/app/services/auth"
)

const adminClaimsKey = "xqyun.admin.claims"

// AdminAuth 校验管理后台接口的 Bearer Access Token。
type AdminAuth struct {
	secret    []byte
	issuer    string
	audience  string
	tokenType string
}

// NewAdminAuth 根据 Application 的 admin 客户端配置创建认证中间件。
func NewAdminAuth(app application.Application) *AdminAuth {
	jwtConfig := app.JWT()
	client, ok := jwtConfig.Client("admin")
	if !ok {
		return &AdminAuth{}
	}
	return newAdminAuth(client.Secret, jwtConfig.Issuer, client.Audience, client.TokenType)
}

func newAdminAuth(secret []byte, issuer, audience, tokenType string) *AdminAuth {
	if issuer == "" {
		issuer = auth.DefaultTokenIssuer
	}
	return &AdminAuth{
		secret:    append([]byte(nil), secret...),
		issuer:    issuer,
		audience:  audience,
		tokenType: tokenType,
	}
}

// Handle 校验 Authorization: Bearer <token>，成功后将 claims 写入请求上下文。
func (m *AdminAuth) Handle(ctx context.Context, c *app.RequestContext) {
	token, ok := bearerToken(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized, "请先登录")
		c.Abort()
		return
	}

	claims, err := auth.VerifyAccessToken(m.secret, token, auth.TokenVerifyOptions{
		Issuer:    m.issuer,
		Audience:  m.audience,
		TokenType: m.tokenType,
	})
	if err != nil {
		response.Fail(c, response.CodeUnauthorized, "登录状态已失效，请重新登录")
		c.Abort()
		return
	}

	c.Set(adminClaimsKey, claims)
	c.Next(ctx)
}

// AdminClaimsFromContext 读取当前请求的管理员 claims。
func AdminClaimsFromContext(c *app.RequestContext) (*auth.TokenClaims, bool) {
	value, exists := c.Get(adminClaimsKey)
	if !exists {
		return nil, false
	}
	claims, ok := value.(*auth.TokenClaims)
	return claims, ok
}

func bearerToken(c *app.RequestContext) (string, bool) {
	header := strings.TrimSpace(string(c.GetHeader("Authorization")))
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}
	return token, true
}
