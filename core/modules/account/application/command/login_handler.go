package command

import (
	"context"
	"errors"

	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/account/application/dto/in"
	"wechat-clone/core/modules/account/application/dto/out"
	"wechat-clone/core/modules/account/application/service"
	repos "wechat-clone/core/modules/account/domain/repos"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

type loginHandler struct {
	authService service.AuthenticationService
}

func NewLoginHandler(_ *appCtx.AppContext, _ repos.Repos, services service.Services) cqrs.Handler[*in.LoginRequest, *out.LoginResponse] {
	return &loginHandler{
		authService: services.AuthenticationService(),
	}
}

func (u *loginHandler) Handle(ctx context.Context, req *in.LoginRequest) (*out.LoginResponse, error) {
	log := logging.FromContext(ctx).Named("Login")

	result, err := u.authService.Authenticate(ctx, service.AuthenticateAccountCommand{
		Email:    req.Email,
		Password: req.Password,
		Device: service.DeviceCommand{
			DeviceUID:  req.DeviceUid,
			DeviceName: req.DeviceName,
			DeviceType: req.DeviceType,
			OSName:     req.OsName,
			OSVersion:  req.OsVersion,
			AppVersion: req.AppVersion,
			UserAgent:  req.UserAgent,
			IPAddress:  req.IpAddress,
		},
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAuthenticationAccountNotFound):
			log.Errorw("Account not found", zap.String("email", req.Email))
			return nil, stackErr.Error(ErrAccountNotFound)
		case errors.Is(err, service.ErrAuthenticationInvalidPassword):
			log.Errorw("Invalid credentials", zap.String("email", req.Email))
			return nil, stackErr.Error(ErrInvalidCredentials)
		default:
			log.Errorw("Login failed", zap.Error(err), zap.String("email", req.Email))
			return nil, stackErr.Error(err)
		}
	}

	return &out.LoginResponse{
		AccessToken:      result.AccessToken,
		AccessExpiresAt:  result.AccessExpiresAt.UnixMilli(),
		RefreshToken:     result.RefreshToken,
		RefreshExpiresAt: result.RefreshExpiresAt.UnixMilli(),
	}, nil
}
