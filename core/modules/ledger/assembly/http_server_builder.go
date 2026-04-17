package assembly

import (
	"context"

	appCtx "wechat-clone/core/context"
	ledgercommand "wechat-clone/core/modules/ledger/application/command"
	ledgerquery "wechat-clone/core/modules/ledger/application/query"
	ledgerserver "wechat-clone/core/modules/ledger/transport/server"
	"wechat-clone/core/shared/pkg/cqrs"
	infrahttp "wechat-clone/core/shared/transport/http"
)

func buildHTTPServer(_ context.Context, appContext *appCtx.AppContext) (infrahttp.HTTPServer, error) {
	ledgerService := BuildService(appContext)
	ledgerQueryService := BuildQueryService(appContext)
	getAccountBalance := cqrs.NewDispatcher(ledgerquery.NewGetAccountBalanceHandler(ledgerQueryService))
	getTransaction := cqrs.NewDispatcher(ledgerquery.NewGetTransactionHandler(ledgerQueryService))
	listTransaction := cqrs.NewDispatcher(ledgerquery.NewListTransactionHandler(ledgerQueryService))
	transferTransaction := cqrs.NewDispatcher(ledgercommand.NewTransferTransaction(appContext, ledgerService))

	return ledgerserver.NewHTTPServer(getAccountBalance, getTransaction, transferTransaction, listTransaction)
}
