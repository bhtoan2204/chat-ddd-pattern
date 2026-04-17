package assembly

import (
	"context"
	appCtx "wechat-clone/core/context"
	notificationcommand "wechat-clone/core/modules/notification/application/command"
	notificationquery "wechat-clone/core/modules/notification/application/query"
	notificationrepo "wechat-clone/core/modules/notification/infra/persistent/repository"
	notificationserver "wechat-clone/core/modules/notification/transport/server"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
	"wechat-clone/core/shared/transport/http"
)

func buildHTTPServer(_ context.Context, appCtx *appCtx.AppContext) (http.HTTPServer, error) {
	notificationRepos := notificationrepo.NewRepoImpl(appCtx)
	notificationReadRepo := notificationrepo.NewNotificationReadRepository(appCtx.GetDB())
	savePushSubscription := cqrs.NewDispatcher(notificationcommand.NewSavePushSubscriptionHandler(notificationRepos))
	listNotification := cqrs.NewDispatcher(notificationquery.NewListNotificationHandler(notificationReadRepo))
	server, err := notificationserver.NewHTTPServer(savePushSubscription, listNotification)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return server, nil
}
