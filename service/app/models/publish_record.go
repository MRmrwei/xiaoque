package models

// PublishRecord 对应 xqyun_publish_records 表，记录页面 / 内容的发布与下架。
type PublishRecord struct {
	ID         uint64  `gorm:"column:id;primaryKey;autoIncrement"`
	TargetType string  `gorm:"column:target_type"`
	TargetID   uint64  `gorm:"column:target_id"`
	Locale     string  `gorm:"column:locale"`
	Status     string  `gorm:"column:status"`
	Message    *string `gorm:"column:message"`
	OperatorID *uint64 `gorm:"column:operator_id"`
	CreatedAt  int64   `gorm:"column:created_at"`
}

func (PublishRecord) TableName() string {
	return "xqyun_publish_records"
}
