package entity

import (
	"go-socket/core/modules/notification/types"
	"time"
)

type NotificationEntity struct {
	ID        string
	AccountID string
	Type      types.NotificationType
	Subject   string
	Body      string
	IsRead    bool
	ReadAt    *time.Time
	CreatedAt time.Time
}
