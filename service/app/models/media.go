package models

// Media 对应 xqyun_media 表，保存上传的媒体素材。
// 文件本体存放在本地磁盘（或对象存储），表中只记录访问地址与元数据。
type Media struct {
	ID           uint64  `gorm:"column:id;primaryKey;autoIncrement"`
	Kind         string  `gorm:"column:kind"`
	OriginalName string  `gorm:"column:original_name"`
	Filename     string  `gorm:"column:filename"`
	URL          string  `gorm:"column:url"`
	Mime         string  `gorm:"column:mime"`
	Size         uint64  `gorm:"column:size"`
	CreatedBy    *uint64 `gorm:"column:created_by"`
	CreatedAt    int64   `gorm:"column:created_at"`
	DeletedAt    uint64  `gorm:"column:deleted_at"`
}

func (Media) TableName() string {
	return "xqyun_media"
}
