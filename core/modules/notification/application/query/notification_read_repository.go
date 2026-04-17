package query

import (
	"context"

	"wechat-clone/core/modules/notification/application/dto/out"
	"wechat-clone/core/shared/utils"
)

//go:generate mockgen -package=query -destination=notification_read_repository_mock.go -source=notification_read_repository.go
type NotificationReadRepository interface {
	ListNotifications(ctx context.Context, options utils.QueryOptions) ([]*out.NotificationResponse, error)
}
