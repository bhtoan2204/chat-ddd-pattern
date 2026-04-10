package assembly

import (
	appCtx "go-socket/core/context"
	"go-socket/core/modules/room/application/messaging"
	"go-socket/core/shared/config"
)

func BuildMessageHandler(cfg *config.Config, appCtx *appCtx.AppContext) (messaging.MessageHandler, error) {
	return messaging.NewMessageHandler(cfg, appCtx)
}
