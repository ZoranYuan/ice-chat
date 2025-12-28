package model

import "time"

type GroupMember struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement"`
	UserID    uint64     `gorm:"not null;uniqueIndex:idx_group_user;index" json:"userId"`
	GroupID   uint64     `gorm:"not null;uniqueIndex:idx_group_user;index" json:"groupId"`
	JoinedAt  time.Time  `gorm:"not null;autoCreateTime" json:"joinedAt"`
	LeftAt    *time.Time `gorm:"index" json:"leftAt,omitempty"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (GroupMember) TableName() string {
	return "group_members"
}
