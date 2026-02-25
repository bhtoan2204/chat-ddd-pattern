package command

import (
	"context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
)

type logoutHandler struct{}

func NewLogoutHandler() LogoutHandler {
	return &logoutHandler{}
}

func (u *logoutHandler) Handle(ctx context.Context, req *in.LogoutRequest) (*out.LogoutResponse, error) {
	_ = ctx
	_ = req
	return &out.LogoutResponse{
		Message: "Logout successful",
	}, nil
}
