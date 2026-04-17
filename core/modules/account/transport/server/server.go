// CODE_GENERATOR: registry
package server

import (
	"context"

	"wechat-clone/core/modules/account/application/dto/in"
	"wechat-clone/core/modules/account/application/dto/out"
	accounthttp "wechat-clone/core/modules/account/transport/http"
	"wechat-clone/core/shared/pkg/cqrs"
	infrahttp "wechat-clone/core/shared/transport/http"

	"github.com/gin-gonic/gin"
)

type accountHTTPServer struct {
	login              cqrs.Dispatcher[*in.LoginRequest, *out.LoginResponse]
	register           cqrs.Dispatcher[*in.RegisterRequest, *out.RegisterResponse]
	logout             cqrs.Dispatcher[*in.LogoutRequest, *out.LogoutResponse]
	refresh            cqrs.Dispatcher[*in.RefreshRequest, *out.RefreshResponse]
	getProfile         cqrs.Dispatcher[*in.GetProfileRequest, *out.GetProfileResponse]
	updateProfile      cqrs.Dispatcher[*in.UpdateProfileRequest, *out.UpdateProfileResponse]
	verifyEmail        cqrs.Dispatcher[*in.VerifyEmailRequest, *out.VerifyEmailResponse]
	confirmVerifyEmail cqrs.Dispatcher[*in.ConfirmVerifyEmailRequest, *out.ConfirmVerifyEmailResponse]
	changePassword     cqrs.Dispatcher[*in.ChangePasswordRequest, *out.ChangePasswordResponse]
	getAvatar          cqrs.Dispatcher[*in.GetAvatarRequest, *out.GetAvatarResponse]
	createPresignedUrl cqrs.Dispatcher[*in.CreatePresignedUrlRequest, *out.CreatePresignedUrlResponse]
	searchUsers        cqrs.Dispatcher[*in.SearchUsersRequest, *out.SearchUsersResponse]
	loginGoogle        cqrs.Dispatcher[*in.LoginGoogleRequest, *out.LoginGoogleResponse]
	callbackGoogle     cqrs.Dispatcher[*in.CallbackGoogleRequest, *out.CallbackGoogleResponse]
}

func NewHTTPServer(
	login cqrs.Dispatcher[*in.LoginRequest, *out.LoginResponse],
	register cqrs.Dispatcher[*in.RegisterRequest, *out.RegisterResponse],
	logout cqrs.Dispatcher[*in.LogoutRequest, *out.LogoutResponse],
	refresh cqrs.Dispatcher[*in.RefreshRequest, *out.RefreshResponse],
	getProfile cqrs.Dispatcher[*in.GetProfileRequest, *out.GetProfileResponse],
	updateProfile cqrs.Dispatcher[*in.UpdateProfileRequest, *out.UpdateProfileResponse],
	verifyEmail cqrs.Dispatcher[*in.VerifyEmailRequest, *out.VerifyEmailResponse],
	confirmVerifyEmail cqrs.Dispatcher[*in.ConfirmVerifyEmailRequest, *out.ConfirmVerifyEmailResponse],
	changePassword cqrs.Dispatcher[*in.ChangePasswordRequest, *out.ChangePasswordResponse],
	getAvatar cqrs.Dispatcher[*in.GetAvatarRequest, *out.GetAvatarResponse],
	createPresignedUrl cqrs.Dispatcher[*in.CreatePresignedUrlRequest, *out.CreatePresignedUrlResponse],
	searchUsers cqrs.Dispatcher[*in.SearchUsersRequest, *out.SearchUsersResponse],
	loginGoogle cqrs.Dispatcher[*in.LoginGoogleRequest, *out.LoginGoogleResponse],
	callbackGoogle cqrs.Dispatcher[*in.CallbackGoogleRequest, *out.CallbackGoogleResponse],
) (infrahttp.HTTPServer, error) {
	return &accountHTTPServer{
		login:              login,
		register:           register,
		logout:             logout,
		refresh:            refresh,
		getProfile:         getProfile,
		updateProfile:      updateProfile,
		verifyEmail:        verifyEmail,
		confirmVerifyEmail: confirmVerifyEmail,
		changePassword:     changePassword,
		getAvatar:          getAvatar,
		createPresignedUrl: createPresignedUrl,
		searchUsers:        searchUsers,
		loginGoogle:        loginGoogle,
		callbackGoogle:     callbackGoogle,
	}, nil
}

func (s *accountHTTPServer) RegisterPublicRoutes(routes *gin.RouterGroup) {
	accounthttp.RegisterPublicRoutes(routes, s.login, s.register, s.refresh, s.confirmVerifyEmail, s.loginGoogle, s.callbackGoogle)
}

func (s *accountHTTPServer) RegisterPrivateRoutes(routes *gin.RouterGroup) {
	accounthttp.RegisterPrivateRoutes(routes, s.logout, s.getProfile, s.updateProfile, s.verifyEmail, s.changePassword, s.getAvatar, s.createPresignedUrl, s.searchUsers)
}

func (s *accountHTTPServer) RegisterSocketRoutes(routes *gin.RouterGroup) {
}

func (s *accountHTTPServer) Stop(ctx context.Context) error {
	return nil
}
