package assembly

import (
	appCtx "wechat-clone/core/context"
	paymentmessaging "wechat-clone/core/modules/payment/application/messaging"
	paymentserver "wechat-clone/core/modules/payment/transport/server"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"
	modruntime "wechat-clone/core/shared/runtime"
)

func BuildMessagingRuntime(cfg *config.Config, appContext *appCtx.AppContext) (modruntime.Module, error) {
	commandService, err := buildPaymentCommandService(appContext)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	messageHandler, err := paymentmessaging.NewMessageHandler(cfg, appContext, commandService)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return paymentserver.NewMessageServer(messageHandler)
}
