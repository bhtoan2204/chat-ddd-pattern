package assembly

import (
	appCtx "go-socket/core/context"
	accountusecase "go-socket/core/modules/account/application/usecase"
	accountrepo "go-socket/core/modules/account/infra/persistent/repository"
)

func BuildAuthUsecase(appCtx *appCtx.AppContext) accountusecase.AuthUsecase {
	accountRepos := accountrepo.NewRepoImpl(appCtx)
	return accountusecase.NewAuthUsecase(appCtx, accountRepos)
}
