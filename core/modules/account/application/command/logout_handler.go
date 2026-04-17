package command

import (
	"context"

	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/account/application/dto/in"
	"wechat-clone/core/modules/account/application/dto/out"
	"wechat-clone/core/modules/account/application/service"
	"wechat-clone/core/modules/account/application/support"
	repos "wechat-clone/core/modules/account/domain/repos"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
)

type logoutHandler struct {
	authService service.AuthenticationService
}

func NewLogoutHandler(appCtx *appCtx.AppContext, baseRepo repos.Repos, services service.Services) cqrs.Handler[*in.LogoutRequest, *out.LogoutResponse] {
	return &logoutHandler{
		authService: services.AuthenticationService(),
	}
}

func (u *logoutHandler) Handle(ctx context.Context, req *in.LogoutRequest) (*out.LogoutResponse, error) {
	accountID, err := support.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	if err := u.authService.Logout(ctx, service.LogoutCommand{
		AccountID:    accountID,
		RefreshToken: req.Token,
	}); err != nil {
		return nil, stackErr.Error(err)
	}

	return &out.LogoutResponse{
		Message: "Logout successful",
	}, nil
}
