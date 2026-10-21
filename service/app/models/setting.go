package models

import "encoding/json"

// Setting 对应 xqyun_settings 表，按 key + locale 保存全站设置。
// Locale 为空字符串表示全局配置，Value 为 JSON。
type Setting struct {
	ID        uint64          `gorm:"column:id;primaryKey;autoIncrement"`
	Key       string          `gorm:"column:key"`
	Locale    string          `gorm:"column:locale"`
	Value     json.RawMessage `gorm:"column:value;type:json"`
	UpdatedBy *uint64         `gorm:"column:updated_by"`
	UpdatedAt int64           `gorm:"column:updated_at"`
}

func (Setting) TableName() string {
	return "xqyun_settings"
}
