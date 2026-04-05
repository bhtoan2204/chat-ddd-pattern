package http

import (
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/transport/http/handler"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/transport/httpx"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(
	routes *gin.RouterGroup,
	login cqrs.Dispatcher[*in.LoginRequest, *out.LoginResponse],
	register cqrs.Dispatcher[*in.RegisterRequest, *out.RegisterResponse],
	confirmVerifyEmail cqrs.Dispatcher[*in.ConfirmVerifyEmailRequest, *out.ConfirmVerifyEmailResponse],
) {
	routes.POST("/auth/login", httpx.Wrap(handler.NewLoginHandler(login)))
	routes.POST("/auth/register", httpx.Wrap(handler.NewRegisterHandler(register)))
	routes.POST("/auth/verify-email/confirm", httpx.Wrap(handler.NewConfirmVerifyEmailHandler(confirmVerifyEmail)))
}

func RegisterPrivateRoutes(
	routes *gin.RouterGroup,
	logout cqrs.Dispatcher[*in.LogoutRequest, *out.LogoutResponse],
	getProfile cqrs.Dispatcher[*in.GetProfileRequest, *out.GetProfileResponse],
	getAvatar cqrs.Dispatcher[*in.GetAvatarRequest, *out.GetAvatarResponse],
	updateProfile cqrs.Dispatcher[*in.UpdateProfileRequest, *out.UpdateProfileResponse],
	verifyEmail cqrs.Dispatcher[*in.VerifyEmailRequest, *out.VerifyEmailResponse],
	changePassword cqrs.Dispatcher[*in.ChangePasswordRequest, *out.ChangePasswordResponse],
) {
	routes.POST("/auth/logout", httpx.Wrap(handler.NewLogoutHandler(logout)))
	routes.GET("/auth/profile", httpx.Wrap(handler.NewGetProfileHandler(getProfile)))
	routes.GET("/auth/avatar/:account_id", httpx.Wrap(handler.NewGetAvatarHandler(getAvatar)))
	routes.PUT("/auth/profile", httpx.Wrap(handler.NewUpdateProfileHandler(updateProfile)))
	routes.POST("/auth/verify-email", httpx.Wrap(handler.NewVerifyEmailHandler(verifyEmail)))
	routes.PUT("/auth/change-password", httpx.Wrap(handler.NewChangePasswordHandler(changePassword)))
}
