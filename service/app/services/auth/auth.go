// Package auth 定义登录业务接口与令牌结构。
package auth

import (
	"context"

	"suki/xqyun/app/models"
)

// TokenPair 是登录成功后签发的令牌。
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	// ExpiresIn 是 Access Token 的有效期，单位为秒。
	ExpiresIn int
}

// LoginResult 是登录成功后的令牌与账号信息。
type LoginResult struct {
	Pair *TokenPair
	User *models.User
}

// Credentials 是登录凭据与请求元信息。
type Credentials struct {
	Username  string
	Password  string
	UserAgent string
	IP        string
}

// Authenticator 定义登录能力。
type Authenticator interface {
	Login(ctx context.Context, credentials Credentials) (*LoginResult, error)
}

// RefreshCredentials 是刷新令牌的凭据与请求元信息。
// RefreshToken 是 HttpOnly Cookie 中的明文，服务端用 SHA-256 后作为 Redis 会话 key 的一部分。
type RefreshCredentials struct {
	RefreshToken string
	UserAgent    string
	IP           string
}

// TokenRefresher 定义刷新 Access Token 的能力。
// 实现需要校验刷新令牌是否有效、未过期、未吊销，且账号仍可用；
// 刷新成功后必须轮换刷新令牌（旧令牌立即失效），返回新的 Access Token 与刷新令牌。
type TokenRefresher interface {
	Refresh(ctx context.Context, credentials RefreshCredentials) (*LoginResult, error)
}

// LogoutCredentials 是退出登录的凭据。
// 现阶段只吊销当前这次会话的刷新令牌；如需「退出所有设备」，后续再加独立用例。
type LogoutCredentials struct {
	RefreshToken string
}

// Logouter 定义退出登录能力。
// 实现按刷新令牌哈希吊销记录；令牌缺失、已过期或已吊销都应视为成功（幂等），
// 以便 handler 无论结果如何都能清除 Cookie。
type Logouter interface {
	Logout(ctx context.Context, credentials LogoutCredentials) error
}

// SessionManager 是管理登录会话的完整能力，由同一个实现同时提供刷新与退出。
// 用例层只需要注入一次，不需要分别传入两个相同实例。
type SessionManager interface {
	TokenRefresher
	Logouter
}
