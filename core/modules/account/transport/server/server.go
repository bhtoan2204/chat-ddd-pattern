package server

import (
	"context"
	accountin "go-socket/core/modules/account/application/dto/in"
	accountout "go-socket/core/modules/account/application/dto/out"
	accounthttp "go-socket/core/modules/account/transport/http"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/transport/http"

	"github.com/gin-gonic/gin"
)

type accountServer struct {
	login              cqrs.Dispatcher[*accountin.LoginRequest, *accountout.LoginResponse]
	register           cqrs.Dispatcher[*accountin.RegisterRequest, *accountout.RegisterResponse]
	logout             cqrs.Dispatcher[*accountin.LogoutRequest, *accountout.LogoutResponse]
	getProfile         cqrs.Dispatcher[*accountin.GetProfileRequest, *accountout.GetProfileResponse]
	getAvatar          cqrs.Dispatcher[*accountin.GetAvatarRequest, *accountout.GetAvatarResponse]
	updateProfile      cqrs.Dispatcher[*accountin.UpdateProfileRequest, *accountout.UpdateProfileResponse]
	verifyEmail        cqrs.Dispatcher[*accountin.VerifyEmailRequest, *accountout.VerifyEmailResponse]
	confirmVerifyEmail cqrs.Dispatcher[*accountin.ConfirmVerifyEmailRequest, *accountout.ConfirmVerifyEmailResponse]
	changePassword     cqrs.Dispatcher[*accountin.ChangePasswordRequest, *accountout.ChangePasswordResponse]
}

func NewServer(
	login cqrs.Dispatcher[*accountin.LoginRequest, *accountout.LoginResponse],
	register cqrs.Dispatcher[*accountin.RegisterRequest, *accountout.RegisterResponse],
	logout cqrs.Dispatcher[*accountin.LogoutRequest, *accountout.LogoutResponse],
	getProfile cqrs.Dispatcher[*accountin.GetProfileRequest, *accountout.GetProfileResponse],
	getAvatar cqrs.Dispatcher[*accountin.GetAvatarRequest, *accountout.GetAvatarResponse],
	updateProfile cqrs.Dispatcher[*accountin.UpdateProfileRequest, *accountout.UpdateProfileResponse],
	verifyEmail cqrs.Dispatcher[*accountin.VerifyEmailRequest, *accountout.VerifyEmailResponse],
	confirmVerifyEmail cqrs.Dispatcher[*accountin.ConfirmVerifyEmailRequest, *accountout.ConfirmVerifyEmailResponse],
	changePassword cqrs.Dispatcher[*accountin.ChangePasswordRequest, *accountout.ChangePasswordResponse],
) (http.HTTPServer, error) {
	return &accountServer{
		login:              login,
		register:           register,
		logout:             logout,
		getProfile:         getProfile,
		getAvatar:          getAvatar,
		updateProfile:      updateProfile,
		verifyEmail:        verifyEmail,
		confirmVerifyEmail: confirmVerifyEmail,
		changePassword:     changePassword,
	}, nil
}

func (s *accountServer) RegisterPublicRoutes(routes *gin.RouterGroup) {
	accounthttp.RegisterPublicRoutes(routes, s.login, s.register, s.confirmVerifyEmail)
}

func (s *accountServer) RegisterPrivateRoutes(routes *gin.RouterGroup) {
	accounthttp.RegisterPrivateRoutes(routes, s.logout, s.getProfile, s.getAvatar, s.updateProfile, s.verifyEmail, s.changePassword)
}

func (s *accountServer) Stop(_ context.Context) error {
	return nil
}
