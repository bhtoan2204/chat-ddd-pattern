// CODE_GENERATOR: application-handler
package command

import (
	"context"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/application/provider"
	"go-socket/core/modules/account/application/service"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/stackErr"
)

type loginGoogleHandler struct {
	authProviderRegistry *provider.AuthProviderRegistry
}

func NewLoginGoogle(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
	services service.Services,
	authProviderRegistry *provider.AuthProviderRegistry,
) cqrs.Handler[*in.LoginGoogleRequest, *out.LoginGoogleResponse] {
	return &loginGoogleHandler{
		authProviderRegistry: authProviderRegistry,
	}
}

func (u *loginGoogleHandler) Handle(ctx context.Context, req *in.LoginGoogleRequest) (*out.LoginGoogleResponse, error) {
	googleProvider, err := u.authProviderRegistry.Get("google")
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return &out.LoginGoogleResponse{
		RedirectURL: googleProvider.Login(),
	}, nil
}
