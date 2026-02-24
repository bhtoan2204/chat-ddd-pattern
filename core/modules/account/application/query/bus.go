package query

import (
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/shared/pkg/cqrs"
)

type GetProfileHandler = cqrs.Handler[*in.GetProfileRequest, *out.GetProfileResponse]

type Bus struct {
	GetProfile cqrs.Dispatcher[*in.GetProfileRequest, *out.GetProfileResponse]
}

func NewBus(getProfileHandler GetProfileHandler) Bus {
	return Bus{
		GetProfile: cqrs.NewDispatcher(getProfileHandler),
	}
}
