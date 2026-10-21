package models

import "encoding/json"

// Menu 对应 xqyun_menus 表，是顶部 / 底部导航的菜单容器。
// Code 取 header / footer / footer_bottom，唯一。
type Menu struct {
	ID          uint64  `gorm:"column:id;primaryKey;autoIncrement"`
	Code        string  `gorm:"column:code"`
	Name        string  `gorm:"column:name"`
	Description *string `gorm:"column:description"`
	CreatedAt   int64   `gorm:"column:created_at"`
	UpdatedAt   int64   `gorm:"column:updated_at"`
}

func (Menu) TableName() string {
	return "xqyun_menus"
}

// MenuItem 对应 xqyun_menu_items 表，同时承载顶部菜单项和页脚栏目 / 链接。
// 页脚中 ParentID 为空的顶级项是栏目，子项是栏目内的链接。
type MenuItem struct {
	ID        uint64          `gorm:"column:id;primaryKey;autoIncrement"`
	MenuID    uint64          `gorm:"column:menu_id"`
	ParentID  *uint64         `gorm:"column:parent_id"`
	Label     json.RawMessage `gorm:"column:label;type:json"`
	LinkType  string          `gorm:"column:link_type"`
	LinkValue *string         `gorm:"column:link_value"`
	Target    string          `gorm:"column:target"`
	Sort      int             `gorm:"column:sort"`
	Visible   int8            `gorm:"column:visible"`
	CreatedAt int64           `gorm:"column:created_at"`
	UpdatedAt int64           `gorm:"column:updated_at"`
}

func (MenuItem) TableName() string {
	return "xqyun_menu_items"
}
