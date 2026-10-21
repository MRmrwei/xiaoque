package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"

	"suki/xqyun/app/core/response"
	"suki/xqyun/app/services/auth"
)

type testResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Username string `json:"username"`
		Subject  string `json:"subject"`
	} `json:"data"`
}

func newTestAdminAuth(secret []byte) *AdminAuth {
	return newAdminAuth(secret, auth.DefaultTokenIssuer, auth.AudienceAdminAPI, auth.TokenTypeAdmin)
}

func signTestToken(t *testing.T, secret []byte, audience, tokenType string, now time.Time, ttl time.Duration) string {
	t.Helper()
	token, err := auth.SignAccessToken(secret, auth.TokenOptions{
		Issuer:    auth.DefaultTokenIssuer,
		Audience:  audience,
		TokenType: tokenType,
		TTL:       ttl,
	}, 7, "admin", now)
	if err != nil {
		t.Fatalf("SignAccessToken() 返回错误：%v", err)
	}
	return token
}

func performProtectedRequest(t *testing.T, middleware *AdminAuth, authorization string) (int, testResponse) {
	t.Helper()

	h := server.New()
	h.Use(middleware.Handle)
	h.GET("/protected", func(ctx context.Context, c *app.RequestContext) {
		claims, ok := AdminClaimsFromContext(c)
		if !ok {
			response.Fail(c, response.CodeInternal, "缺少 claims")
			return
		}
		response.Success(c, map[string]string{
			"username": claims.Username,
			"subject":  claims.Subject,
		})
	})

	var headers []ut.Header
	if authorization != "" {
		headers = append(headers, ut.Header{Key: "Authorization", Value: authorization})
	}
	recorder := ut.PerformRequest(h.Engine, http.MethodGet, "/protected", nil, headers...)

	var body testResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("解析响应失败：%v，body=%s", err, recorder.Body.String())
	}
	return recorder.Code, body
}

func TestAdminAuthAllowsAdminToken(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	token := signTestToken(t, secret, auth.AudienceAdminAPI, auth.TokenTypeAdmin, time.Now(), time.Minute)

	status, body := performProtectedRequest(t, newTestAdminAuth(secret), "Bearer "+token)
	if status != http.StatusOK {
		t.Fatalf("HTTP 状态 = %d，期望 200", status)
	}
	if body.Code != response.CodeOK {
		t.Fatalf("业务码 = %d，期望 0：%s", body.Code, body.Message)
	}
	if body.Data.Username != "admin" || body.Data.Subject != "7" {
		t.Fatalf("claims 未正确写入上下文：%#v", body.Data)
	}
}

func TestAdminAuthAcceptsLowercaseBearer(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	token := signTestToken(t, secret, auth.AudienceAdminAPI, auth.TokenTypeAdmin, time.Now(), time.Minute)

	_, body := performProtectedRequest(t, newTestAdminAuth(secret), "bearer "+token)
	if body.Code != response.CodeOK {
		t.Fatalf("小写 bearer 应被接受，实际 code=%d message=%s", body.Code, body.Message)
	}
}

func TestAdminAuthRejectsInvalidRequests(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	validToken := signTestToken(t, secret, auth.AudienceAdminAPI, auth.TokenTypeAdmin, time.Now(), time.Minute)
	frontendToken := signTestToken(t, secret, auth.AudienceFrontendAPI, auth.TokenTypeFrontend, time.Now(), time.Minute)
	expiredToken := signTestToken(t, secret, auth.AudienceAdminAPI, auth.TokenTypeAdmin, time.Now().Add(-2*time.Hour), time.Hour)

	parts := strings.Split(validToken, ".")
	parts[1] = "A" + parts[1][1:]
	tamperedToken := strings.Join(parts, ".")

	cases := []struct {
		name          string
		authorization string
	}{
		{name: "missing header", authorization: ""},
		{name: "invalid scheme", authorization: "Token " + validToken},
		{name: "tampered token", authorization: "Bearer " + tamperedToken},
		{name: "frontend token", authorization: "Bearer " + frontendToken},
		{name: "expired token", authorization: "Bearer " + expiredToken},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			status, body := performProtectedRequest(t, newTestAdminAuth(secret), testCase.authorization)
			if status != http.StatusOK {
				t.Fatalf("失败响应 HTTP 状态 = %d，期望 200", status)
			}
			if body.Code != response.CodeUnauthorized {
				t.Fatalf("业务码 = %d，期望 %d", body.Code, response.CodeUnauthorized)
			}
			if body.Data.Username != "" {
				t.Fatalf("未认证请求不应执行受保护 handler：%#v", body.Data)
			}
		})
	}
}
