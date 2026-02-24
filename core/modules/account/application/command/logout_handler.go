package command

import (
	"context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
)

type logoutUseCase struct{}

func NewLogoutUseCase() LogoutHandler {
	return &logoutUseCase{}
}

func (u *logoutUseCase) Handle(ctx context.Context, req *in.LogoutRequest) (*out.LogoutResponse, error) {
	_ = ctx
	_ = req
	return &out.LogoutResponse{
		Message: "Logout successful",
	}, nil
}
