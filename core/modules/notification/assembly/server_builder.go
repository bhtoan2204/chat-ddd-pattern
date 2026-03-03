package assembly

import (
	"go-socket/config"
	appCtx "go-socket/core/context"
	notificationserver "go-socket/core/modules/notification/transport/server"
	stackerr "go-socket/core/shared/pkg/stackErr"
)

func BuildServer(cfg *config.Config, appCtx *appCtx.AppContext) (notificationserver.Server, error) {
	messageHandler, err := BuildMessageHandler(cfg, appCtx)
	if err != nil {
		return nil, stackerr.Error(err)
	}

	return notificationserver.NewServer(messageHandler)
}
