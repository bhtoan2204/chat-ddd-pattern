package assembly

import (
	appCtx "go-socket/core/context"
	notificationserver "go-socket/core/modules/notification/transport/server"
	"go-socket/core/shared/config"
	"go-socket/core/shared/pkg/stackErr"
	modruntime "go-socket/core/shared/runtime"
)

func buildMessagingRuntime(cfg *config.Config, appCtx *appCtx.AppContext) (modruntime.Module, error) {
	messageHandler, err := buildMessagingHandler(cfg, appCtx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return notificationserver.NewServer(messageHandler)
}
