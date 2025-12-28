package model

import (
	"ice-chat/pkg/snowflake"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserId          uint64         `gorm:"column:user_id;type:bigint unsigned;not null;primaryKey" json:"userId"`
	Username        string         `gorm:"column:username;type:varchar(50);not null;unique" json:"username"`
	Avatar          string         `gorm:"column:avatar;type:varchar(255);default:''" json:"avatar"`
	Password        string         `gorm:"column:password;type:varchar(100);not null" json:"-"`
	Email           string         `gorm:"column:email;type:varchar(100);not null" json:"email"`
	Status          int8           `gorm:"column:status;type:tinyint;default:1" json:"status"`
	LastOfflineTime *time.Time     `gorm:"column:last_offline_time;type:datetime;default:null" json:"lastOfflineTime"` // 新增字段
	CreatedAt       time.Time      `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index" json:"-"`
}

// TableName 指定表名
func (u User) TableName() string {
	return "users"
}

// BeforeCreate 钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserId == 0 {
		u.UserId = snowflake.NewID() // 调用雪花ID生成函数
	}
	return nil
}
