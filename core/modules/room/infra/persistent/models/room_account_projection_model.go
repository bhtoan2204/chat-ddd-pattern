package models

import "time"

type RoomAccountProjection struct {
	AccountID       string    `gorm:"primaryKey"`
	DisplayName     string    `gorm:"default:''"`
	Username        string    `gorm:"default:''"`
	AvatarObjectKey string    `gorm:"default:''"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (r *RoomAccountProjection) TableName() string {
	return "room_account_projections"
}
