package models

import "encoding/json"

// FormSubmission 对应 xqyun_form_submissions 表，保存联系表单等提交的留言。
// name / email / phone / company / message 是从 Data 中提取的常用字段，
// 完整内容以 Data 为准。
type FormSubmission struct {
	ID        uint64          `gorm:"column:id;primaryKey;autoIncrement"`
	FormKey   string          `gorm:"column:form_key"`
	Locale    *string         `gorm:"column:locale"`
	Name      *string         `gorm:"column:name"`
	Email     *string         `gorm:"column:email"`
	Phone     *string         `gorm:"column:phone"`
	Company   *string         `gorm:"column:company"`
	Message   *string         `gorm:"column:message"`
	Data      json.RawMessage `gorm:"column:data;type:json"`
	Status    string          `gorm:"column:status"`
	IP        *string         `gorm:"column:ip"`
	UserAgent *string         `gorm:"column:user_agent"`
	CreatedAt int64           `gorm:"column:created_at"`
	UpdatedAt int64           `gorm:"column:updated_at"`
}

func (FormSubmission) TableName() string {
	return "xqyun_form_submissions"
}
