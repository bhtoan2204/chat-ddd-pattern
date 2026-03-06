package query

import (
	"go-socket/core/modules/notification/application/dto/in"
	"go-socket/core/modules/notification/application/dto/out"
	"go-socket/core/shared/pkg/cqrs"
)

type ListNotificationHandler = cqrs.Handler[*in.ListNotificationRequest, *out.ListNotificationResponse]

type Bus struct {
	ListNotification cqrs.Dispatcher[*in.ListNotificationRequest, *out.ListNotificationResponse]
}

func NewBus(listNotificationHandler ListNotificationHandler) Bus {
	return Bus{
		ListNotification: cqrs.NewDispatcher(listNotificationHandler),
	}
}
