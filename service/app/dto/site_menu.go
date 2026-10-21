package dto

// AdminMenuItem 是后台菜单项，Label 同时保存中英文文案。
// 顶部菜单只有一级；页脚中顶级项是栏目，Children 是栏目内链接。
type AdminMenuItem struct {
	ID       string            `json:"id"`
	Label    map[string]string `json:"label"`
	Path     string            `json:"path"`
	LinkType string            `json:"link_type"`
	Target   string            `json:"target"`
	Visible  *bool             `json:"visible"`
	Sort     int               `json:"sort"`
	Children []AdminMenuItem   `json:"children,omitempty"`
}

// AdminMenuSummary 是菜单容器列表项。
type AdminMenuSummary struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   int64  `json:"updated_at"`
}

// AdminMenuDetail 是菜单详情与菜单项。
type AdminMenuDetail struct {
	AdminMenuSummary
	Items []AdminMenuItem `json:"items"`
}

// AdminMenuListResponse 是菜单列表响应。
type AdminMenuListResponse struct {
	Items []AdminMenuSummary `json:"items"`
}

// AdminMenuSaveRequest 是整体保存菜单请求，Items 的顺序即展示顺序。
type AdminMenuSaveRequest struct {
	Name  string          `json:"name"`
	Items []AdminMenuItem `json:"items"`
}
