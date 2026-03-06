package assembly

import (
	appCtx "go-socket/core/context"
	notificationcommand "go-socket/core/modules/notification/application/command"
	notificationrepo "go-socket/core/modules/notification/infra/persistent/repository"
	notificationserver "go-socket/core/modules/notification/transport/server"
	stackerr "go-socket/core/shared/pkg/stackErr"
)

func BuildHTTPServer(appCtx *appCtx.AppContext) (notificationserver.HTTPServer, error) {
	notificationRepos := notificationrepo.NewRepoImpl(appCtx)
	savePushSubscriptionHandler := notificationcommand.NewSavePushSubscriptionHandler(notificationRepos)
	commandBus := notificationcommand.NewBus(savePushSubscriptionHandler)

	server, err := notificationserver.NewHTTPServer(commandBus)
	if err != nil {
		return nil, stackerr.Error(err)
	}

	return server, nil
}
