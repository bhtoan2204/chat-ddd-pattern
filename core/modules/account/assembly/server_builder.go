package assembly

import (
	"context"
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/account/application/command"
	"wechat-clone/core/modules/account/application/provider"
	"wechat-clone/core/modules/account/application/provider/google"
	"wechat-clone/core/modules/account/application/query"
	accountservice "wechat-clone/core/modules/account/application/service"
	accountrepo "wechat-clone/core/modules/account/infra/persistent/repository"
	accountserver "wechat-clone/core/modules/account/transport/server"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
	"wechat-clone/core/shared/transport/http"
)

func buildHTTPServer(ctx context.Context, appContext *appCtx.AppContext) (http.HTTPServer, error) {
	accountRepos := accountrepo.NewRepoImpl(appContext.GetDB(), appContext.GetCache())
	accountServices := accountservice.NewServices(appContext, accountRepos)
	authProviderRegistry := provider.NewProviderRegistry()
	authProviderRegistry.Register(google.NewGoogleProvider(ctx, appContext.GetConfig()))

	login := cqrs.NewDispatcher(command.NewLoginHandler(appContext, accountRepos, accountServices))
	register := cqrs.NewDispatcher(command.NewRegisterHandler(appContext, accountRepos, accountServices))
	logout := cqrs.NewDispatcher(command.NewLogoutHandler(appContext, accountRepos, accountServices))
	getProfile := cqrs.NewDispatcher(query.NewGetProfileHandler(appContext, accountRepos, accountServices))
	getAvatar := cqrs.NewDispatcher(query.NewGetAvatarHandler(appContext, accountRepos, accountServices))
	getPresignedUrl := cqrs.NewDispatcher(command.NewCreatePresignedUrlHandler(appContext, accountRepos, accountServices))
	updateProfile := cqrs.NewDispatcher(command.NewUpdateProfileHandler(appContext, accountRepos, accountServices))
	verifyEmail := cqrs.NewDispatcher(command.NewVerifyEmailHandler(appContext, accountRepos, accountServices))
	confirmVerifyEmail := cqrs.NewDispatcher(command.NewConfirmVerifyEmailHandler(appContext, accountRepos, accountServices))
	changePassword := cqrs.NewDispatcher(command.NewChangePasswordHandler(appContext, accountRepos, accountServices))
	searchUsers := cqrs.NewDispatcher(query.NewSearchUsers(appContext, accountRepos, accountServices))
	refresh := cqrs.NewDispatcher(command.NewRefresh(appContext, accountRepos, accountServices))
	loginGoogle := cqrs.NewDispatcher(command.NewLoginGoogle(appContext, accountRepos, accountServices, authProviderRegistry))
	callbackGoogle := cqrs.NewDispatcher(command.NewCallbackGoogle(appContext, accountRepos, accountServices, authProviderRegistry))
	server, err := accountserver.NewHTTPServer(
		login,
		register,
		logout,
		refresh,
		getProfile,
		updateProfile,
		verifyEmail,
		confirmVerifyEmail,
		changePassword,
		getAvatar,
		getPresignedUrl,
		searchUsers,
		loginGoogle,
		callbackGoogle,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return server, nil
}
