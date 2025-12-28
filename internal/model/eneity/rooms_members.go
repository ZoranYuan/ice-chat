package model

import "time"

type RoomsMember struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement"`
	UserID    uint64     `gorm:"not null;uniqueIndex:idx_group_user;index" json:"userId"`
	RoomID    uint64     `gorm:"not null;uniqueIndex:idx_group_user;index" json:"roomId"`
	JoinedAt  time.Time  `gorm:"not null;autoCreateTime" json:"joinedAt"`
	LeftAt    *time.Time `gorm:"index" json:"leftAt,omitempty"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (RoomsMember) TableName() string {
	return "room_members"
}
