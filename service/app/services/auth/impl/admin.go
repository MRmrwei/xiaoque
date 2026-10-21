// Package impl 提供 auth.Authenticator 的具体实现。
package impl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sessionrepository "suki/xqyun/app/repository/session"
	userrepository "suki/xqyun/app/repository/user"
	"suki/xqyun/app/services/auth"
)

const (
	maxUserAgentLength = 500
	maxIPLength        = 45
)

// AdminLoginConfig 描述管理后台登录服务的令牌参数。
type AdminLoginConfig struct {
	Secret     []byte
	Issuer     string
	Audience   string
	TokenType  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// AdminLogin 校验 xqyun_users 中的后台账号，签发 JWT Access Token 与刷新令牌。
type AdminLogin struct {
	users      userrepository.UserRepository
	sessions   sessionrepository.Repository
	secret     []byte
	issuer     string
	audience   string
	tokenType  string
	accessTTL  time.Duration
	refreshTTL time.Duration
	now        func() time.Time
}

var _ auth.Authenticator = (*AdminLogin)(nil)

// NewAdminLogin 创建管理后台登录服务。
func NewAdminLogin(
	users userrepository.UserRepository,
	sessions sessionrepository.Repository,
	cfg AdminLoginConfig,
) *AdminLogin {
	if cfg.Issuer == "" {
		cfg.Issuer = auth.DefaultTokenIssuer
	}
	if cfg.Audience == "" {
		cfg.Audience = auth.AudienceAdminAPI
	}
	if cfg.TokenType == "" {
		cfg.TokenType = auth.TokenTypeAdmin
	}
	return &AdminLogin{
		users:      users,
		sessions:   sessions,
		secret:     append([]byte(nil), cfg.Secret...),
		issuer:     cfg.Issuer,
		audience:   cfg.Audience,
		tokenType:  cfg.TokenType,
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
		now:        time.Now,
	}
}

// Login 校验账号密码，签发 JWT Access Token，并保存刷新令牌的 SHA-256。
func (s *AdminLogin) Login(ctx context.Context, credentials auth.Credentials) (*auth.LoginResult, error) {
	username := strings.TrimSpace(credentials.Username)
	if username == "" || credentials.Password == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if s.users == nil || s.sessions == nil {
		return nil, fmt.Errorf("认证服务：数据访问层未初始化")
	}

	user, err := s.users.FindByUsername(ctx, username)
	if errors.Is(err, userrepository.ErrNotFound) {
		return nil, auth.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if !auth.VerifyPassword(credentials.Password, user.PasswordHash) {
		return nil, auth.ErrInvalidCredentials
	}
	if user.Status != 1 {
		return nil, auth.ErrAccountDisabled
	}

	now := s.now()
	accessToken, err := auth.SignAccessToken(s.secret, auth.TokenOptions{
		Issuer:    s.issuer,
		Audience:  s.audience,
		TokenType: s.tokenType,
		TTL:       s.accessTTL,
	}, user.ID, user.Username, now)
	if err != nil {
		return nil, err
	}
	refreshPlain, refreshHash, err := auth.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	userAgent := optionalString(credentials.UserAgent, maxUserAgentLength)
	ip := optionalString(credentials.IP, maxIPLength)

	nowUnix := now.Unix()
	sessionToken := sessionrepository.Token{
		UserID:    user.ID,
		TokenHash: refreshHash,
		TokenType: s.tokenType,
		UserAgent: userAgent,
		IP:        ip,
		ExpiresAt: now.Add(s.refreshTTL).Unix(),
		CreatedAt: nowUnix,
	}

	if err := s.sessions.Save(ctx, sessionToken); err != nil {
		return nil, err
	}
	// 登录信息只用于展示，写入失败时回滚刚创建的会话，避免留下不可用的令牌。
	if err := s.users.UpdateLoginState(ctx, user.ID, nowUnix, ip); err != nil {
		_ = s.sessions.Revoke(ctx, refreshHash, s.tokenType)
		return nil, err
	}

	user.LastLoginAt = nowUnix
	user.LastLoginIP = ip
	user.UpdatedAt = nowUnix

	return &auth.LoginResult{
		Pair: &auth.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshPlain,
			ExpiresIn:    int(s.accessTTL / time.Second),
		},
		User: user,
	}, nil
}

func truncate(value string, max int) string {
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}
