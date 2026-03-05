package models

import (
	"go-socket/core/modules/notification/types"
	"time"
)

type NotificationModel struct {
	ID        string                 `gorm:"primaryKey;type:char(36)"`
	AccountID string                 `gorm:"not null;index"`
	Type      types.NotificationType `gorm:"not null;index"`
	Subject   string                 `gorm:"not null"`
	Body      string                 `gorm:"not null"`
	IsRead    bool                   `gorm:"default:false"`
	ReadAt    *time.Time             `gorm:"index"`
	CreatedAt time.Time              `gorm:"autoCreateTime;index"`
}

func (NotificationModel) TableName() string {
	return "notifications"
}
