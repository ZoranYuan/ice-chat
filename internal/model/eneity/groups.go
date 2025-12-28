package model

import (
	"ice-chat/pkg/snowflake"
	"time"

	"gorm.io/gorm"
)

type Groups struct {
	GroupId    uint64         `gorm:"column:group_id;type:bigint unsigned;not null;primaryKey" json:"groupId"` // 群组ID
	GroupName  string         `json:"groupName" gorm:"column:group_name;not null"`                             // 群组名称
	Avatar     string         `json:"avatar" gorm:"avatar"`
	Desc       string         `json:"desc" gorm:"avatar"`
	CreateUser uint64         `json:"createUser" gorm:"createUser;not null"`
	Members    int            `json:"members" gorm:"column:mambers;not null;default:0"`           // 最大成员数
	CreatedAt  time.Time      `json:"createdAt" gorm:"column:created_at;not null;autoCreateTime"` // 创建时间
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"column:updated_at;not null;autoUpdateTime"` // 更新时间
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`                           // 软删除（-：不序列化到JSON）
}

// TableName 指定数据库表名
func (g *Groups) TableName() string {
	return "groups"
}

// BeforeCreate 钩子
func (g *Groups) BeforeCreate(tx *gorm.DB) error {
	if g.GroupId == 0 {
		g.GroupId = snowflake.NewID() // 调用雪花ID生成函数
	}
	return nil
}
