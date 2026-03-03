package assembly

import (
	"go-socket/config"
	appCtx "go-socket/core/context"
	notificationmessaging "go-socket/core/modules/notification/application/messaging"
)

func BuildMessageHandler(cfg *config.Config, appCtx *appCtx.AppContext) (notificationmessaging.MessageHandler, error) {
	return notificationmessaging.NewMessageHandler(cfg, appCtx.GetSMTP())
}
