package models

// User 对应 xqyun_users 表，保存后台账号。
type User struct {
	ID           uint64  `gorm:"column:id;primaryKey;autoIncrement"`
	Username     string  `gorm:"column:username"`
	PasswordHash string  `gorm:"column:password_hash"`
	Name         string  `gorm:"column:name"`
	Email        *string `gorm:"column:email"`
	Mobile       *string `gorm:"column:mobile"`
	Avatar       *string `gorm:"column:avatar;type:mediumtext"`
	Status       int8    `gorm:"column:status"`
	LastLoginAt  int64   `gorm:"column:last_login_at"`
	LastLoginIP  *string `gorm:"column:last_login_ip"`
	CreatedAt    int64   `gorm:"column:created_at"`
	UpdatedAt    int64   `gorm:"column:updated_at"`
	DeletedAt    uint64  `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "xqyun_users"
}
