package model

import (
	"time"

	"gorm.io/gorm"
)

// 用户模型
type User struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name;type:varchar(100);not null"`
	Emali     string         `gorm:"column:email;type:varchar(100);uniqueIndex;not null"`
	Password  string         `gorm:"column:password;type:varchar(255);not null"`
	Phone     string         `gorm:"column:phone;type:varchar(20);uniqueIndex"`
	Avater    string         `gorm:"column:avatar;type:varchar(500)"`
	Bio       string         `gorm:"column:bio;type:text"`
	Gender    string         `gorm:"column:gender;type:tinyint"`
	LastLogin int64          `gorm:"column:last_login"`
	Status    UserStatus     `gorm:"column:status;not null;default:1;index"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	BannedAt  *time.Time     `gorm:"column:banned_at"`
	BanReason string         `gorm:"column:ban_reason;type:varchar(255)"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

// 表名
func (User) TableName() string {
	return "users"
}
