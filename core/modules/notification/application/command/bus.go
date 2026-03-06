package command

import (
	"go-socket/core/modules/notification/application/dto/in"
	"go-socket/core/modules/notification/application/dto/out"
	"go-socket/core/shared/pkg/cqrs"
)

type SavePushSubscriptionHandler = cqrs.Handler[*in.SavePushSubscriptionRequest, *out.SavePushSubscriptionResponse]

type Bus struct {
	SavePushSubscription cqrs.Dispatcher[*in.SavePushSubscriptionRequest, *out.SavePushSubscriptionResponse]
}

func NewBus(savePushSubscriptionHandler SavePushSubscriptionHandler) Bus {
	return Bus{
		SavePushSubscription: cqrs.NewDispatcher(savePushSubscriptionHandler),
	}
}
