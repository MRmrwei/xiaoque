package auth

import "errors"

var (
	// ErrInvalidCredentials 表示账号不存在或密码错误。
	// 对外不区分两者，避免泄漏账号是否存在。
	ErrInvalidCredentials = errors.New("登录失败：账号或密码错误")
	// ErrAccountDisabled 表示账号已被禁用。
	ErrAccountDisabled = errors.New("登录失败：账号已被禁用")
	// ErrInvalidAccessToken 表示 Access Token 格式、签名或 claims 不合法。
	ErrInvalidAccessToken = errors.New("访问令牌无效")
	// ErrAccessTokenExpired 表示 Access Token 已过期。
	ErrAccessTokenExpired = errors.New("访问令牌已过期")
	// ErrInvalidRefreshToken 表示刷新令牌缺失、无效、已过期或已被吊销。
	// 对外统一提示重新登录，不区分具体原因。
	ErrInvalidRefreshToken = errors.New("登录状态已失效，请重新登录")
)
