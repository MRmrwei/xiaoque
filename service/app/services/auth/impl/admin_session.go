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

// AdminSession 负责管理后台刷新令牌的轮换与退出登录。
// 只处理刷新令牌，不校验 Access Token（过期后仍然需要能刷新）。
type AdminSession struct {
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

var (
	_ auth.TokenRefresher = (*AdminSession)(nil)
	_ auth.Logouter       = (*AdminSession)(nil)
)

// NewAdminSession 创建管理后台刷新 / 退出登录服务。
func NewAdminSession(
	users userrepository.UserRepository,
	sessions sessionrepository.Repository,
	cfg AdminLoginConfig,
) *AdminSession {
	if cfg.Issuer == "" {
		cfg.Issuer = auth.DefaultTokenIssuer
	}
	if cfg.Audience == "" {
		cfg.Audience = auth.AudienceAdminAPI
	}
	if cfg.TokenType == "" {
		cfg.TokenType = auth.TokenTypeAdmin
	}
	return &AdminSession{
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

// Refresh 校验刷新令牌、确认账号可用，然后轮换令牌并签发新的 Access Token。
func (s *AdminSession) Refresh(
	ctx context.Context,
	credentials auth.RefreshCredentials,
) (*auth.LoginResult, error) {
	plain := strings.TrimSpace(credentials.RefreshToken)
	if plain == "" {
		return nil, auth.ErrInvalidRefreshToken
	}
	if s.users == nil || s.sessions == nil {
		return nil, fmt.Errorf("认证服务：数据访问层未初始化")
	}

	tokenHash := auth.HashRefreshToken(plain)
	stored, err := s.sessions.FindByHash(ctx, tokenHash, s.tokenType)
	if errors.Is(err, sessionrepository.ErrNotFound) {
		return nil, auth.ErrInvalidRefreshToken
	}
	if err != nil {
		return nil, err
	}

	now := s.now()
	// Redis 中已吊销的令牌会被删除，这里只做一次过期兜底。
	if stored.ExpiresAt <= now.Unix() {
		return nil, auth.ErrInvalidRefreshToken
	}

	user, err := s.users.FindByID(ctx, stored.UserID)
	if errors.Is(err, userrepository.ErrNotFound) {
		return nil, auth.ErrInvalidRefreshToken
	}
	if err != nil {
		return nil, err
	}
	if user.Status != 1 {
		return nil, auth.ErrAccountDisabled
	}

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

	nowUnix := now.Unix()
	next := sessionrepository.Token{
		UserID:    user.ID,
		TokenHash: refreshHash,
		TokenType: s.tokenType,
		UserAgent: optionalString(credentials.UserAgent, maxUserAgentLength),
		IP:        optionalString(credentials.IP, maxIPLength),
		ExpiresAt: now.Add(s.refreshTTL).Unix(),
		CreatedAt: nowUnix,
	}

	if err := s.sessions.Rotate(ctx, tokenHash, s.tokenType, next); err != nil {
		if errors.Is(err, sessionrepository.ErrRefreshTokenRotated) {
			return nil, auth.ErrInvalidRefreshToken
		}
		return nil, err
	}

	return &auth.LoginResult{
		Pair: &auth.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshPlain,
			ExpiresIn:    int(s.accessTTL / time.Second),
		},
		User: user,
	}, nil
}

// Logout 吊销当前刷新令牌。
// 未携带令牌或令牌已失效时同样返回成功，保证重复退出不会报错。
func (s *AdminSession) Logout(ctx context.Context, credentials auth.LogoutCredentials) error {
	plain := strings.TrimSpace(credentials.RefreshToken)
	if plain == "" {
		return nil
	}
	if s.sessions == nil {
		return fmt.Errorf("认证服务：数据访问层未初始化")
	}

	if err := s.sessions.Revoke(
		ctx,
		auth.HashRefreshToken(plain),
		s.tokenType,
	); err != nil {
		return err
	}
	return nil
}

// optionalString 去除首尾空白并截断，空值返回 nil。
func optionalString(value string, max int) *string {
	value = truncate(strings.TrimSpace(value), max)
	if value == "" {
		return nil
	}
	return &value
}
