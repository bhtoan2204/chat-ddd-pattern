package models

import (
	"time"
	"wechat-clone/core/modules/room/types"
)

type RoomMemberModel struct {
	ID              string         `gorm:"primaryKey"`
	RoomID          string         `gorm:"not null;index"`
	AccountID       string         `gorm:"not null;index"`
	Role            types.RoomRole `gorm:"default:member"`
	LastDeliveredAt *time.Time
	LastReadAt      *time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (RoomMemberModel) TableName() string {
	return "room_members"
}
