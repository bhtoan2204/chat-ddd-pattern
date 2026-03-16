package assembly

import (
	appCtx "go-socket/core/context"
	"go-socket/core/modules/payment/application/command"
	"go-socket/core/modules/payment/application/query"
	paymentrepo "go-socket/core/modules/payment/infra/persistent/repository"
)

type Buses struct {
	commandBus command.Bus
	queryBus   query.Bus
}

func BuildBuses(appCtx *appCtx.AppContext) Buses {
	paymentRepos := paymentrepo.NewRepoImpl(appCtx)

	depositHandler := command.NewDepositHandler(paymentRepos)
	rebuildProjectionHandler := command.NewRebuildProjectionHandler(paymentRepos)
	transferHandler := command.NewTransferHandler(paymentRepos)
	withdrawalHandler := command.NewWithdrawalHandler(paymentRepos)
	commandBus := command.NewBus(depositHandler, rebuildProjectionHandler, transferHandler, withdrawalHandler)

	newListTransactionHandler := query.NewListTransactionHandler(paymentRepos)
	queryBus := query.NewBus(newListTransactionHandler)
	return Buses{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}
