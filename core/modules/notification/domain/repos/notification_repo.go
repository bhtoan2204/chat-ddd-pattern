package repos

import (
	"context"
	"go-socket/core/modules/notification/application/dto/out"
	"go-socket/core/modules/notification/domain/entity"
	"go-socket/core/shared/utils"
)

type NotificationRepository interface {
	// Commands
	CreateNotification(ctx context.Context, notification *entity.NotificationEntity) error

	// Queries
	ListNotifications(ctx context.Context, options utils.QueryOptions) ([]*out.NotificationResponse, error)
}
