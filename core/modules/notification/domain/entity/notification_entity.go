package entity

import (
	"time"
	"wechat-clone/core/modules/notification/types"
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
