package query

import (
	"context"
	"go-socket/core/modules/notification/application/dto/in"
	"go-socket/core/modules/notification/application/dto/out"
	"go-socket/core/modules/notification/domain/repos"
	"go-socket/core/shared/pkg/logging"
	stackerr "go-socket/core/shared/pkg/stackErr"
	"go-socket/core/shared/utils"

	"go.uber.org/zap"
)

type listNotificationHandler struct {
	notificationRepo repos.NotificationRepository
}

func NewListNotificationHandler(notificationRepo repos.NotificationRepository) ListNotificationHandler {
	return &listNotificationHandler{
		notificationRepo: notificationRepo,
	}
}

func (h *listNotificationHandler) Handle(ctx context.Context, req *in.ListNotificationRequest) (*out.ListNotificationResponse, error) {
	log := logging.FromContext(ctx).Named("ListNotification")
	options := utils.QueryOptions{}
	if req.Limit > 0 {
		options.Limit = &req.Limit
	}
	notifications, err := h.notificationRepo.ListNotifications(ctx, options)
	if err != nil {
		log.Errorw("Failed to list notifications", zap.Error(err))
		return nil, stackerr.Error(err)
	}

	items := make([]out.NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		if notification == nil {
			continue
		}
		items = append(items, *notification)
	}

	return &out.ListNotificationResponse{
		Notifications: items,
		Limit:         req.Limit,
		Total:         len(items),
	}, nil
}
