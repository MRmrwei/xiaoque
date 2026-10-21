package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestSignAccessTokenUsesHS256JWT(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	now := time.Unix(1792375200, 0)
	options := TokenOptions{
		Issuer:    "xiaoque-admin",
		Audience:  AudienceAdminAPI,
		TokenType: TokenTypeAdmin,
		TTL:       15 * time.Minute,
	}

	token, err := SignAccessToken(secret, options, 42, "admin", now)
	if err != nil {
		t.Fatalf("SignAccessToken() 返回错误：%v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("JWT 应该有 3 段，实际 %d 段：%q", len(parts), token)
	}

	headerRaw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("解码 JWT header 失败：%v", err)
	}
	var header map[string]string
	if err := json.Unmarshal(headerRaw, &header); err != nil {
		t.Fatalf("解析 JWT header 失败：%v", err)
	}
	if header["alg"] != "HS256" || header["typ"] != "JWT" {
		t.Fatalf("JWT header 不正确：%#v", header)
	}

	payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("解码 JWT payload 失败：%v", err)
	}
	var claims TokenClaims
	if err := json.Unmarshal(payloadRaw, &claims); err != nil {
		t.Fatalf("解析 JWT payload 失败：%v", err)
	}
	if claims.Subject != "42" || claims.Username != "admin" || claims.Issuer != "xiaoque-admin" {
		t.Fatalf("JWT claims 不正确：%#v", claims)
	}
	if claims.Audience != AudienceAdminAPI || claims.TokenType != TokenTypeAdmin {
		t.Fatalf("JWT 受众 / 类型不正确：%#v", claims)
	}
	if claims.IssuedAt != now.Unix() || claims.ExpiresAt != now.Add(15*time.Minute).Unix() {
		t.Fatalf("JWT 时间字段不正确：%#v", claims)
	}

	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(parts[0] + "." + parts[1]))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		t.Fatal("JWT 签名与密钥不匹配")
	}
}

func TestSignAccessTokenRequiresSecretAndAudience(t *testing.T) {
	now := time.Now()
	options := TokenOptions{
		Audience:  AudienceAdminAPI,
		TokenType: TokenTypeAdmin,
		TTL:       time.Minute,
	}
	if _, err := SignAccessToken(nil, options, 1, "admin", now); err == nil {
		t.Fatal("未配置密钥时 SignAccessToken() 应返回错误")
	}

	options.Audience = ""
	options.TokenType = ""
	if _, err := SignAccessToken([]byte("secret"), options, 1, "admin", now); err == nil {
		t.Fatal("未配置受众 / 类型时 SignAccessToken() 应返回错误")
	}
}

func TestVerifyAccessToken(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	now := time.Unix(1792375200, 0)
	options := TokenOptions{
		Issuer:    "xiaoque-admin",
		Audience:  AudienceAdminAPI,
		TokenType: TokenTypeAdmin,
		TTL:       15 * time.Minute,
	}
	token, err := SignAccessToken(secret, options, 42, "admin", now)
	if err != nil {
		t.Fatalf("SignAccessToken() 返回错误：%v", err)
	}

	claims, err := verifyAccessToken(secret, token, TokenVerifyOptions{
		Issuer:    "xiaoque-admin",
		Audience:  AudienceAdminAPI,
		TokenType: TokenTypeAdmin,
	}, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("VerifyAccessToken() 返回错误：%v", err)
	}
	if claims.Subject != "42" || claims.Username != "admin" ||
		claims.Audience != AudienceAdminAPI || claims.TokenType != TokenTypeAdmin {
		t.Fatalf("VerifyAccessToken() claims 不正确：%#v", claims)
	}
}

func TestVerifyAccessTokenRejectsInvalidTokens(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	now := time.Unix(1792375200, 0)
	verifyOptions := TokenVerifyOptions{
		Issuer:    "xiaoque-admin",
		Audience:  AudienceAdminAPI,
		TokenType: TokenTypeAdmin,
	}
	sign := func(options TokenOptions) string {
		token, err := SignAccessToken(secret, options, 42, "admin", now)
		if err != nil {
			t.Fatalf("SignAccessToken() 返回错误：%v", err)
		}
		return token
	}

	valid := sign(TokenOptions{
		Issuer:    "xiaoque-admin",
		Audience:  AudienceAdminAPI,
		TokenType: TokenTypeAdmin,
		TTL:       15 * time.Minute,
	})
	parts := strings.Split(valid, ".")
	parts[1] = "A" + parts[1][1:]
	tampered := strings.Join(parts, ".")

	wrongAudience := sign(TokenOptions{
		Issuer:    "xiaoque-admin",
		Audience:  AudienceFrontendAPI,
		TokenType: TokenTypeAdmin,
		TTL:       15 * time.Minute,
	})
	wrongType := sign(TokenOptions{
		Issuer:    "xiaoque-admin",
		Audience:  AudienceAdminAPI,
		TokenType: TokenTypeFrontend,
		TTL:       15 * time.Minute,
	})
	wrongIssuer := sign(TokenOptions{
		Issuer:    "other-service",
		Audience:  AudienceAdminAPI,
		TokenType: TokenTypeAdmin,
		TTL:       15 * time.Minute,
	})

	cases := []struct {
		name  string
		token string
		at    time.Time
	}{
		{name: "tampered", token: tampered, at: now},
		{name: "wrong audience", token: wrongAudience, at: now},
		{name: "wrong type", token: wrongType, at: now},
		{name: "wrong issuer", token: wrongIssuer, at: now},
		{name: "expired", token: valid, at: now.Add(16 * time.Minute)},
		{name: "not a jwt", token: "not-a-jwt", at: now},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if _, err := verifyAccessToken(secret, testCase.token, verifyOptions, testCase.at); !errors.Is(err, ErrInvalidAccessToken) && !errors.Is(err, ErrAccessTokenExpired) {
				t.Fatalf("应返回令牌错误，实际 %v", err)
			}
		})
	}

	if _, err := verifyAccessToken(secret, valid, verifyOptions, now.Add(16*time.Minute)); !errors.Is(err, ErrAccessTokenExpired) {
		t.Fatalf("过期 token 应返回 ErrAccessTokenExpired，实际 %v", err)
	}
}
