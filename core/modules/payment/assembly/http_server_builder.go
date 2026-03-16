package assembly

import (
	"context"
	appCtx "go-socket/core/context"
	paymentserver "go-socket/core/modules/payment/transport/server"
	"go-socket/core/shared/transport/http"
)

func BuildHTTPServer(_ context.Context, appContext *appCtx.AppContext) (http.HTTPServer, error) {
	buses := BuildBuses(appContext)
	commandBus := buses.commandBus
	queryBus := buses.queryBus
	return paymentserver.NewHTTPServer(commandBus, queryBus)
}
