// Package adminauth 组装管理后台登录用例。
// 输入校验、认证调用、业务错误映射与响应 DTO 组装都在这一层完成，handler 只负责 HTTP 参数与响应。
package adminauth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"suki/xqyun/app/core/response"
	"suki/xqyun/app/dto"
	"suki/xqyun/app/services/auth"
)

// Error 是可直接映射为统一响应结构的业务错误。
// HTTP 状态固定 200，客户端按 Code 判断具体错误。
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	// ErrInvalidInput 表示账号或密码为空。
	ErrInvalidInput = &Error{
		Code:    response.CodeBadRequest,
		Message: "请输入用户名和密码",
	}
	// ErrInvalidCredentials 表示账号不存在或密码错误，对外不区分两者。
	ErrInvalidCredentials = &Error{
		Code:    response.CodeUnauthorized,
		Message: "用户名或密码错误",
	}
	// ErrAccountDisabled 表示账号已禁用。
	ErrAccountDisabled = &Error{
		Code:    response.CodeForbidden,
		Message: "账号已被禁用",
	}
	// ErrSessionExpired 表示刷新令牌缺失、无效、过期或已吊销，需要重新登录。
	ErrSessionExpired = &Error{
		Code:    response.CodeUnauthorized,
		Message: "登录状态已失效，请重新登录",
	}
)

// LoginInput 是登录用例的输入。
type LoginInput struct {
	Username  string
	Password  string
	UserAgent string
	IP        string
}

// RefreshInput 是刷新令牌用例的输入。
type RefreshInput struct {
	RefreshToken string
	UserAgent    string
	IP           string
}

// LogoutInput 是退出登录用例的输入。
type LogoutInput struct {
	RefreshToken string
}

// Service 是管理后台认证用例服务，覆盖登录、刷新令牌与退出登录。
type Service struct {
	authenticator auth.Authenticator
	session       auth.SessionManager
}

// NewLoginService 创建管理后台登录用例服务。
func NewLoginService(authenticator auth.Authenticator) *Service {
	return &Service{authenticator: authenticator}
}

// NewSessionService 创建刷新令牌 / 退出登录用例服务。
// 传入 *authimpl.AdminSession 这类同时实现刷新与退出的会话服务。
func NewSessionService(session auth.SessionManager) *Service {
	return &Service{session: session}
}

// Login 校验输入、调用认证服务，并返回可直接输出的登录响应。
func (s *Service) Login(ctx context.Context, input LoginInput) (*dto.AdminLoginResponse, error) {
	username := strings.TrimSpace(input.Username)
	if username == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}
	if s.authenticator == nil {
		return nil, fmt.Errorf("登录服务未初始化")
	}

	result, err := s.authenticator.Login(ctx, auth.Credentials{
		Username:  username,
		Password:  input.Password,
		UserAgent: input.UserAgent,
		IP:        input.IP,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			return nil, ErrInvalidCredentials
		case errors.Is(err, auth.ErrAccountDisabled):
			return nil, ErrAccountDisabled
		default:
			return nil, err
		}
	}
	if result == nil || result.Pair == nil || result.User == nil {
		return nil, fmt.Errorf("登录服务返回了空结果")
	}
	return adminResponse(result), nil
}

// Refresh 校验刷新令牌并轮换，返回新的 Access Token 与用户信息。
func (s *Service) Refresh(ctx context.Context, input RefreshInput) (*dto.AdminLoginResponse, error) {
	refreshToken := strings.TrimSpace(input.RefreshToken)
	if refreshToken == "" {
		return nil, ErrSessionExpired
	}
	if s.session == nil {
		return nil, fmt.Errorf("刷新令牌服务未初始化")
	}

	result, err := s.session.Refresh(ctx, auth.RefreshCredentials{
		RefreshToken: refreshToken,
		UserAgent:    input.UserAgent,
		IP:           input.IP,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidRefreshToken):
			return nil, ErrSessionExpired
		case errors.Is(err, auth.ErrAccountDisabled):
			return nil, ErrAccountDisabled
		default:
			return nil, err
		}
	}
	if result == nil || result.Pair == nil || result.User == nil {
		return nil, fmt.Errorf("刷新令牌服务返回了空结果")
	}
	return adminResponse(result), nil
}

// Logout 吊销当前刷新令牌；未携带令牌时直接视为成功。
func (s *Service) Logout(ctx context.Context, input LogoutInput) error {
	refreshToken := strings.TrimSpace(input.RefreshToken)
	if refreshToken == "" {
		return nil
	}
	if s.session == nil {
		return fmt.Errorf("退出登录服务未初始化")
	}
	return s.session.Logout(ctx, auth.LogoutCredentials{RefreshToken: refreshToken})
}

// adminResponse 把认证结果组装为登录 / 刷新共用的响应 DTO。
func adminResponse(result *auth.LoginResult) *dto.AdminLoginResponse {
	user := result.User
	return &dto.AdminLoginResponse{
		AccessToken:  result.Pair.AccessToken,
		RefreshToken: result.Pair.RefreshToken,
		ExpiresIn:    result.Pair.ExpiresIn,
		User: dto.AdminUser{
			ID:          strconv.FormatUint(user.ID, 10),
			Username:    user.Username,
			Name:        user.Name,
			Email:       stringValue(user.Email),
			Mobile:      stringValue(user.Mobile),
			Avatar:      stringValue(user.Avatar),
			Roles:       []string{},
			Permissions: []string{},
		},
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
