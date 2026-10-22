// Package impl 提供 user 包业务接口的具体实现。
package impl

import (
	"context"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	sessionrepository "suki/xqyun/app/repository/session"
	userrepository "suki/xqyun/app/repository/user"
	"suki/xqyun/app/services/auth"
	"suki/xqyun/app/services/user"
)

// minPasswordRunes 是新密码的最小字符数（中文按字符计算）。
const minPasswordRunes = 8

// PasswordChanger 是 user.PasswordChanger 的实现。
type PasswordChanger struct {
	users     userrepository.UserRepository
	sessions  sessionrepository.Repository
	tokenType string
	now       func() time.Time
}

var _ user.PasswordChanger = (*PasswordChanger)(nil)

// NewPasswordChanger 创建修改密码用例。
// tokenType 与管理后台 JWT 的 token_type 一致（admin），用于定位 Redis 会话。
func NewPasswordChanger(
	users userrepository.UserRepository,
	sessions sessionrepository.Repository,
	tokenType string,
) *PasswordChanger {
	return &PasswordChanger{
		users:     users,
		sessions:  sessions,
		tokenType: tokenType,
		now:       time.Now,
	}
}

// ChangePassword 校验当前密码后重新哈希并保存新密码。
func (s *PasswordChanger) ChangePassword(
	ctx context.Context,
	input user.ChangePasswordInput,
) error {
	if input.UserID == 0 {
		return user.ErrUserNotFound
	}

	account, err := s.users.FindByID(ctx, input.UserID)
	if errors.Is(err, userrepository.ErrNotFound) {
		return user.ErrUserNotFound
	}
	if err != nil {
		return err
	}

	if !auth.VerifyPassword(input.CurrentPassword, account.PasswordHash) {
		return user.ErrCurrentPasswordIncorrect
	}
	if utf8.RuneCountInString(input.NewPassword) < minPasswordRunes {
		return user.ErrPasswordTooShort
	}

	hashed, err := auth.HashPassword(input.NewPassword)
	if err != nil {
		if errors.Is(err, auth.ErrPasswordTooLong) {
			return user.ErrPasswordTooLong
		}
		return err
	}

	if err := s.users.UpdatePassword(ctx, input.UserID, hashed, s.now().Unix()); err != nil {
		if errors.Is(err, userrepository.ErrNotFound) {
			return user.ErrUserNotFound
		}
		return err
	}

	s.revokeSessions(ctx, input)
	return nil
}

// revokeSessions 吊销当前令牌与该用户的全部刷新令牌，强制所有设备重新登录。
// 密码已经修改成功，这里失败只记日志：用户可以先用新密码重新登录，
// 未吊销的会话也会在 refresh_ttl（7 天）后自然失效。
func (s *PasswordChanger) revokeSessions(ctx context.Context, input user.ChangePasswordInput) {
	if s.sessions == nil {
		return
	}

	if token := input.RefreshToken; token != "" {
		hash := auth.HashRefreshToken(token)
		if err := s.sessions.Revoke(ctx, hash, s.tokenType); err != nil {
			hlog.CtxWarnf(ctx, "改密后吊销当前刷新令牌失败：%v", err)
		}
	}
	if err := s.sessions.RevokeByUser(ctx, input.UserID, s.tokenType); err != nil {
		hlog.CtxWarnf(ctx, "改密后吊销用户全部刷新令牌失败：%v", err)
	}
}
