package dto

// PageBlock 是页面区块，只保存类型与属性，不保存渲染后的 HTML。
type PageBlock struct {
	ID    string         `json:"id"`
	Type  string         `json:"type"`
	Props map[string]any `json:"props"`
}

// SEO 是页面 SEO 信息。
type SEO struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
}

// AdminPageSummary 是后台页面列表项。
type AdminPageSummary struct {
	ID           string                 `json:"id"`
	GroupID      string                 `json:"group_id"`
	Locale       string                 `json:"locale"`
	Path         string                 `json:"path"`
	Name         string                 `json:"name"`
	Status       string                 `json:"status"`
	Version      int                    `json:"version"`
	UpdatedAt    int64                  `json:"updated_at"`
	UpdatedBy    string                 `json:"updated_by"`
	Translations []AdminPageTranslation `json:"translations"`
}

// AdminPageTranslation 是同一翻译组下的其他语言版本。
type AdminPageTranslation struct {
	PageID string `json:"page_id"`
	Locale string `json:"locale"`
	Path   string `json:"path"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// AdminPageDraft 是后台页面详情（当前草稿）。
type AdminPageDraft struct {
	ID        string      `json:"id"`
	GroupID   string      `json:"group_id"`
	Locale    string      `json:"locale"`
	Path      string      `json:"path"`
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	Version   int         `json:"version"`
	SEO       SEO         `json:"seo"`
	Blocks    []PageBlock `json:"blocks"`
	UpdatedAt int64       `json:"updated_at"`
}

// AdminPageListResponse 是后台页面列表的分页响应。
type AdminPageListResponse struct {
	Items    []AdminPageSummary `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// AdminPageCreateRequest 是新建页面请求。
// CopyFromID 非空时复制来源页面的草稿结构（区块 ID 重新生成）。
type AdminPageCreateRequest struct {
	Locale     string      `json:"locale"`
	Path       string      `json:"path"`
	Name       string      `json:"name"`
	SEO        *SEO        `json:"seo"`
	Blocks     []PageBlock `json:"blocks"`
	CopyFromID string      `json:"copy_from_id"`
}

// AdminPageSaveDraftRequest 是保存草稿请求。
// SEO 为 nil 时保留原值；Blocks 为 nil 时保留原区块，传空数组表示清空。
type AdminPageSaveDraftRequest struct {
	Path   string      `json:"path"`
	Name   string      `json:"name"`
	SEO    *SEO        `json:"seo"`
	Blocks []PageBlock `json:"blocks"`
}

// AdminPageTranslateRequest 是创建翻译配对草稿请求。
// Blocks 为空时服务端复制源页面结构并清空可翻译文本；传入时以传入内容为准。
type AdminPageTranslateRequest struct {
	FromLocale string      `json:"from_locale"`
	ToLocale   string      `json:"to_locale"`
	Path       string      `json:"path"`
	Name       string      `json:"name"`
	SEO        *SEO        `json:"seo"`
	Blocks     []PageBlock `json:"blocks"`
}

// AdminPagePublishRequest 是发布页面请求，Note 写入发布记录。
type AdminPagePublishRequest struct {
	Note string `json:"note"`
}

// AdminPagePublishResponse 是发布结果。
type AdminPagePublishResponse struct {
	Version     int   `json:"version"`
	PublishedAt int64 `json:"published_at"`
}
