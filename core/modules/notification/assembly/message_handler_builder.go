package assembly

import (
	appCtx "wechat-clone/core/context"
	notificationmessaging "wechat-clone/core/modules/notification/application/messaging"
	notificationrepo "wechat-clone/core/modules/notification/infra/persistent/repository"
	"wechat-clone/core/shared/config"
)

func buildMessagingHandler(cfg *config.Config, appCtx *appCtx.AppContext) (notificationmessaging.MessageHandler, error) {
	repos := notificationrepo.NewRepoImpl(appCtx)
	return notificationmessaging.NewMessageHandler(cfg, appCtx.GetSMTP(), repos.NotificationRepository())
}
