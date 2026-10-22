package dto

// PublishedPage 是官网前台读取的已发布页面快照，只包含线上内容。
// Alternates 只包含已发布语言，未发布语言不出现，避免 hreflang 指向 404。
type PublishedPage struct {
	ID          string            `json:"id"`
	GroupID     string            `json:"group_id"`
	Locale      string            `json:"locale"`
	Path        string            `json:"path"`
	Name        string            `json:"name"`
	SEO         SEO               `json:"seo"`
	Blocks      []PageBlock       `json:"blocks"`
	Version     int               `json:"version"`
	PublishedAt int64             `json:"published_at"`
	PublishedBy string            `json:"published_by"`
	Alternates  map[string]string `json:"alternates"`
}

// NavItem 是前台导航项。Label 已按请求语言取好，前台直接展示。
// 顶部菜单只有一级；页脚顶级项是栏目，Children 是栏目内链接。
type NavItem struct {
	ID       string    `json:"id"`
	Label    string    `json:"label"`
	Path     string    `json:"path"`
	LinkType string    `json:"link_type"`
	Target   string    `json:"target"`
	Visible  bool      `json:"visible"`
	Sort     int       `json:"sort"`
	Children []NavItem `json:"children,omitempty"`
}

// SiteNav 是顶部菜单与页脚栏目，对应 xqyun_menus 的三个菜单容器。
type SiteNav struct {
	Header       []NavItem `json:"header"`
	Footer       []NavItem `json:"footer"`
	FooterBottom []NavItem `json:"footer_bottom"`
}

// RenderPageData 是前台页面渲染的聚合结果，服务端一次返回避免逐区块请求。
// Site 是脱敏后的站点配置，按 key 直接给出已解析的 JSON 值（不含统计代码等敏感内容）。
// Collections 按区块类型（news / products / cases）预取已发布内容，内容库上线前为空数组。
type RenderPageData struct {
	Page        PublishedPage    `json:"page"`
	Nav         SiteNav          `json:"nav"`
	Site        map[string]any   `json:"site"`
	Collections map[string][]any `json:"collections"`
}
