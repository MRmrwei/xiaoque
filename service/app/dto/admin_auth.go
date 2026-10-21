package dto

// AdminLoginRequest 是管理后台登录请求。
type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AdminLoginResponse 是管理后台登录响应。
type AdminLoginResponse struct {
	AccessToken string `json:"access_token"`
	// RefreshToken 不返回给客户端，由 handler 写入 HttpOnly Cookie。
	RefreshToken string    `json:"-"`
	ExpiresIn    int       `json:"expires_in"`
	User         AdminUser `json:"user"`
}

// AdminUser 是管理后台登录用户信息。
type AdminUser struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Mobile      string   `json:"mobile"`
	Avatar      string   `json:"avatar"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
