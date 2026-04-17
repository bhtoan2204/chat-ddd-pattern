package models

import (
	"time"
	"wechat-clone/core/modules/room/types"
)

type RoomModel struct {
	ID              string         `gorm:"primaryKey"`
	Name            string         `gorm:"not null"`
	Description     string         `gorm:"default:''"`
	RoomType        types.RoomType `gorm:"not null"`
	OwnerID         string         `gorm:"not null"`
	DirectKey       *string        `gorm:"index"`
	PinnedMessageID *string
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (RoomModel) TableName() string {
	return "rooms"
}
