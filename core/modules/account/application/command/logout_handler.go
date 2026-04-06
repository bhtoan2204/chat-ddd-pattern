package command

import (
	"context"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/application/service"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
)

type logoutHandler struct{}

func NewLogoutHandler(appCtx *appCtx.AppContext, baseRepo repos.Repos, services service.Services) cqrs.Handler[*in.LogoutRequest, *out.LogoutResponse] {
	return &logoutHandler{}
}

func (u *logoutHandler) Handle(ctx context.Context, req *in.LogoutRequest) (*out.LogoutResponse, error) {
	return &out.LogoutResponse{
		Message: "Logout successful",
	}, nil
}
