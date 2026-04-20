package assembly

import (
	appCtx "wechat-clone/core/context"
	notificationmessaging "wechat-clone/core/modules/notification/application/messaging"
	notificationservice "wechat-clone/core/modules/notification/application/service"
	notificationrepo "wechat-clone/core/modules/notification/infra/persistent/repository"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"
)

func buildMessagingHandler(cfg *config.Config, appCtx *appCtx.AppContext) (notificationmessaging.MessageHandler, error) {
	repos, err := notificationrepo.NewRepoImpl(appCtx)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	services := notificationservice.NewServices(appCtx, repos)
	return notificationmessaging.NewMessageHandler(cfg, repos, services)
}
