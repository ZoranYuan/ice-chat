package model

import (
	"ice-chat/pkg/snowflake"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID          uint64    `gorm:"column:id;type:bigint unsigned;not null;primaryKey" json:"id"`
	UserID      uint64    `gorm:"not null;index" json:"user_id"`
	Username    string    `gorm:"type:varchar(64);not null" json:"username"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	MessageType int8      `gorm:"type:tinyint;default:0;not null" json:"message_type"`
	Status      int8      `gorm:"type:tinyint;default:0;not null" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (c Message) TableName() string {
	return "messages"
}

// BeforeCreate 钩子
func (c *Message) BeforeCreate(tx *gorm.DB) error {
	if c.ID == 0 {
		c.ID = snowflake.NewID() // 调用雪花ID生成函数
	}
	return nil
}
