package dto

// AdminProfileUpdateRequest 是修改当前管理员资料请求。
// 邮箱、手机号、头像传空字符串表示清空。
type AdminProfileUpdateRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
	Avatar string `json:"avatar"`
}

// AdminProfile 是当前管理员资料。
type AdminProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar"`
}

// AdminPasswordChangeRequest 是修改当前管理员密码请求。
type AdminPasswordChangeRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
