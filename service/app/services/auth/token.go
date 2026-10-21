package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DefaultTokenIssuer 是 JWT 默认签发方；可通过 config.yaml 的 jwt.issuer 覆盖。
const DefaultTokenIssuer = "xiaoque-admin"

// 令牌受众与类型：现阶段只有 admin，frontend 为后续前台登录预留。
const (
	// AudienceAdminAPI 是管理后台接口的 JWT 受众，写入 claims.aud。
	AudienceAdminAPI = "admin-api"
	// AudienceFrontendAPI 是官网前台接口的预留受众。
	AudienceFrontendAPI = "frontend-api"

	// TokenTypeAdmin 表示管理后台用户，sub 对应 xqyun_users.id。
	TokenTypeAdmin = "admin"
	// TokenTypeFrontend 表示前台用户，预留给后续前台用户表。
	TokenTypeFrontend = "frontend"
)

// TokenClaims 是 Access Token 中的业务声明，对应 JWT 的 payload 部分。
// 示例：
//
//	{
//	  "sub": "1",
//	  "username": "admin",
//	  "iss": "xiaoque-admin",
//	  "aud": "admin-api",
//	  "token_type": "admin",
//	  "iat": 1792375200,
//	  "exp": 1792376100
//	}
//
// aud 和 token_type 用来区分管理后台与前台用户，避免两套用户表的 ID 冲突。
type TokenClaims struct {
	// Subject 是用户 ID（JWT 标准字段 sub）。
	// 管理后台对应 xqyun_users.id；前台接入后需要结合 TokenType 判断属于哪张用户表。
	Subject string `json:"sub"`

	// Username 是管理后台登录名；前台 token 可留空或放前台账号 / 昵称。
	Username string `json:"username"`

	// Issuer 是签发方（JWT 标准字段 iss），默认 xiaoque-admin。
	// 校验方可以用它确认 token 是否由本服务签发。
	Issuer string `json:"iss"`

	// Audience 是 token 的接收方（JWT 标准字段 aud）。
	// 管理后台为 admin-api，前台预留 frontend-api；不同接口组应校验对应的 aud。
	Audience string `json:"aud"`

	// TokenType 是客户端类型：admin / frontend。
	// 中间件用它决定查询 xqyun_users 还是后续的前台用户表。
	TokenType string `json:"token_type"`

	// IssuedAt 是签发时间，Unix 秒（JWT 标准字段 iat）。
	IssuedAt int64 `json:"iat"`

	// ExpiresAt 是过期时间，Unix 秒（JWT 标准字段 exp）。
	// 校验时要求当前时间早于该值。
	ExpiresAt int64 `json:"exp"`
}

// TokenOptions 是签发 Access Token 的参数。
// 现阶段管理后台传入 admin 的受众 / 类型；后续前台登录可传入 frontend 对应值。
type TokenOptions struct {
	// Issuer 写入 iss；为空时使用 DefaultTokenIssuer。
	Issuer string
	// Audience 写入 aud，例如 admin-api。
	Audience string
	// TokenType 写入 token_type，例如 admin。
	TokenType string
	// TTL 是 Access Token 的有效期。
	TTL time.Duration
}

// SignAccessToken 使用 HS256 签发 Access Token。
// 结构为 header.payload.signature：
//   - header 固定 {"alg":"HS256","typ":"JWT"}；
//   - payload 是 TokenClaims 的 JSON；
//   - signature 是 HMAC-SHA256(secret, header.payload) 的 base64url 编码。
func SignAccessToken(secret []byte, options TokenOptions, userID uint64, username string, now time.Time) (string, error) {
	if len(secret) == 0 {
		return "", fmt.Errorf("认证服务：JWT 密钥未配置")
	}
	if options.TTL <= 0 {
		return "", fmt.Errorf("认证服务：令牌有效期无效")
	}
	if options.Audience == "" || options.TokenType == "" {
		return "", fmt.Errorf("认证服务：令牌受众或类型未配置")
	}

	issuer := options.Issuer
	if issuer == "" {
		issuer = DefaultTokenIssuer
	}

	payload, err := json.Marshal(TokenClaims{
		Subject:   strconv.FormatUint(userID, 10),
		Username:  username,
		Issuer:    issuer,
		Audience:  options.Audience,
		TokenType: options.TokenType,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(options.TTL).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("认证服务：编码访问令牌失败：%w", err)
	}

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	body := base64.RawURLEncoding.EncodeToString(payload)
	signingInput := header + "." + body

	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature, nil
}

// TokenVerifyOptions 是校验 Access Token 时要求匹配的签发方、受众与客户端类型。
type TokenVerifyOptions struct {
	Issuer    string
	Audience  string
	TokenType string
}

// VerifyAccessToken 校验 JWT 结构、签名、有效期、签发方、受众与客户端类型。
func VerifyAccessToken(secret []byte, token string, options TokenVerifyOptions) (*TokenClaims, error) {
	return verifyAccessToken(secret, token, options, time.Now())
}

func verifyAccessToken(
	secret []byte,
	token string,
	options TokenVerifyOptions,
	now time.Time,
) (*TokenClaims, error) {
	if len(secret) == 0 || options.Audience == "" || options.TokenType == "" {
		return nil, ErrInvalidAccessToken
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, ErrInvalidAccessToken
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrInvalidAccessToken
	}
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(parts[0] + "." + parts[1]))
	if !hmac.Equal(signature, mac.Sum(nil)) {
		return nil, ErrInvalidAccessToken
	}

	headerRaw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrInvalidAccessToken
	}
	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	if err := json.Unmarshal(headerRaw, &header); err != nil ||
		header.Algorithm != "HS256" ||
		header.Type != "JWT" {
		return nil, ErrInvalidAccessToken
	}

	payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidAccessToken
	}
	var claims TokenClaims
	if err := json.Unmarshal(payloadRaw, &claims); err != nil {
		return nil, ErrInvalidAccessToken
	}
	if claims.Subject == "" || claims.IssuedAt <= 0 || claims.ExpiresAt <= 0 {
		return nil, ErrInvalidAccessToken
	}
	if claims.ExpiresAt <= now.Unix() {
		return nil, ErrAccessTokenExpired
	}
	if options.Issuer != "" && claims.Issuer != options.Issuer {
		return nil, ErrInvalidAccessToken
	}
	if claims.Audience != options.Audience || claims.TokenType != options.TokenType {
		return nil, ErrInvalidAccessToken
	}

	return &claims, nil
}

// NewRefreshToken 生成刷新令牌原文与存储用的 SHA-256 哈希。
func NewRefreshToken() (plain, hash string, err error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("认证服务：生成刷新令牌失败：%w", err)
	}

	plain = base64.RawURLEncoding.EncodeToString(raw)
	return plain, HashRefreshToken(plain), nil
}

// HashRefreshToken 计算刷新令牌的 SHA-256 十六进制摘要。
// 数据库只保存摘要，明文只通过 HttpOnly Cookie 往返；刷新 / 退出时用同一算法比对。
func HashRefreshToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
