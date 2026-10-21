package models

import "encoding/json"

// Page 对应 xqyun_pages 表，保存一种语言、一个路径的页面草稿与线上内容。
// 中英文页面通过 TranslationGroupID 配对，各自独立发布。
type Page struct {
	ID                 uint64          `gorm:"column:id;primaryKey;autoIncrement"`
	TranslationGroupID string          `gorm:"column:translation_group_id;type:char(26)"`
	Locale             string          `gorm:"column:locale"`
	Path               string          `gorm:"column:path"`
	Name               string          `gorm:"column:name"`
	Status             string          `gorm:"column:status"`
	DraftContent       json.RawMessage `gorm:"column:draft_content;type:json"`
	PublishedContent   json.RawMessage `gorm:"column:published_content;type:json"`
	Version            uint32          `gorm:"column:version"`
	CreatedBy          *uint64         `gorm:"column:created_by"`
	UpdatedBy          *uint64         `gorm:"column:updated_by"`
	CreatedAt          int64           `gorm:"column:created_at"`
	UpdatedAt          int64           `gorm:"column:updated_at"`
	DeletedAt          uint64          `gorm:"column:deleted_at"`
}

func (Page) TableName() string {
	return "xqyun_pages"
}
