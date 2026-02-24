package command

import (
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/shared/pkg/cqrs"
)

type LoginHandler = cqrs.Handler[*in.LoginRequest, *out.LoginResponse]
type RegisterHandler = cqrs.Handler[*in.RegisterRequest, *out.RegisterResponse]
type LogoutHandler = cqrs.Handler[*in.LogoutRequest, *out.LogoutResponse]

type Bus struct {
	Login    cqrs.Dispatcher[*in.LoginRequest, *out.LoginResponse]
	Register cqrs.Dispatcher[*in.RegisterRequest, *out.RegisterResponse]
	Logout   cqrs.Dispatcher[*in.LogoutRequest, *out.LogoutResponse]
}

func NewBus(
	loginHandler LoginHandler,
	registerHandler RegisterHandler,
	logoutHandler LogoutHandler,
) Bus {
	return Bus{
		Login:    cqrs.NewDispatcher(loginHandler),
		Register: cqrs.NewDispatcher(registerHandler),
		Logout:   cqrs.NewDispatcher(logoutHandler),
	}
}
