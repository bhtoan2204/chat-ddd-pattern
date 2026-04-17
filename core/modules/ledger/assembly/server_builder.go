package assembly

import (
	appCtx "wechat-clone/core/context"
	ledgermessaging "wechat-clone/core/modules/ledger/application/messaging"
	ledgerserver "wechat-clone/core/modules/ledger/transport/server"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"
	modruntime "wechat-clone/core/shared/runtime"
)

func buildMessagingRuntime(cfg *config.Config, appCtx *appCtx.AppContext) (modruntime.Module, error) {
	messageHandler, err := ledgermessaging.NewMessageHandler(cfg, appCtx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return ledgerserver.NewServer(messageHandler)
}
