package assembly

import (
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/ledger/application/service"
	ledgerrepo "wechat-clone/core/modules/ledger/infra/persistent/repository"
)

func BuildService(appContext *appCtx.AppContext) service.LedgerService {
	return service.NewLedgerService(ledgerrepo.NewRepoImpl(appContext))
}

func BuildQueryService(appContext *appCtx.AppContext) service.LedgerQueryService {
	return service.NewLedgerQueryService(ledgerrepo.NewRepoImpl(appContext))
}
