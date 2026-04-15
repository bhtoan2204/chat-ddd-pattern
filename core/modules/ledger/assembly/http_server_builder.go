package assembly

import (
	"context"

	appCtx "go-socket/core/context"
	ledgercommand "go-socket/core/modules/ledger/application/command"
	ledgerquery "go-socket/core/modules/ledger/application/query"
	repos "go-socket/core/modules/ledger/infra/persistent/repository"
	ledgerserver "go-socket/core/modules/ledger/transport/server"
	"go-socket/core/shared/pkg/cqrs"
	infrahttp "go-socket/core/shared/transport/http"
)

func buildHTTPServer(_ context.Context, appContext *appCtx.AppContext) (infrahttp.HTTPServer, error) {
	ledgerService := BuildService(appContext)
	ledgerRepos := repos.NewRepoImpl(appContext)

	ledgerQueryService := BuildQueryService(appContext)
	createTransaction := cqrs.NewDispatcher(ledgercommand.NewCreateTransactionHandler(ledgerService))
	getAccountBalance := cqrs.NewDispatcher(ledgerquery.NewGetAccountBalanceHandler(ledgerQueryService))
	getTransaction := cqrs.NewDispatcher(ledgerquery.NewGetTransactionHandler(ledgerQueryService))
	transferTransaction := cqrs.NewDispatcher(ledgercommand.NewTransferTransaction(appContext, ledgerRepos, ledgerService))

	return ledgerserver.NewHTTPServer(createTransaction, getAccountBalance, getTransaction, transferTransaction)
}
