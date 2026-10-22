package dto

// AdminMediaItem 是后台素材库中的一条素材。
type AdminMediaItem struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	URL          string `json:"url"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	Mime         string `json:"mime"`
	Size         int64  `json:"size"`
	CreatedAt    int64  `json:"created_at"`
}

// AdminMediaUpdateRequest 是修改素材信息请求。
type AdminMediaUpdateRequest struct {
	OriginalName string `json:"original_name"`
}

// AdminMediaListResponse 是素材列表分页响应。
type AdminMediaListResponse struct {
	Items    []AdminMediaItem `json:"items"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}
