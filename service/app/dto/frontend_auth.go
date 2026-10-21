package dto

// FrontendLoginRequest 是官网前台用户登录请求。
// 具体登录方式（邮箱密码、手机验证码、第三方登录）确定后再扩展字段。
type FrontendLoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// FrontendLoginResponse 是官网前台用户登录响应。
type FrontendLoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         FrontendUser `json:"user"`
}

// FrontendUser 是官网前台登录用户信息。
type FrontendUser struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}
