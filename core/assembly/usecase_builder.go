package assembly

import (
	appCtx "go-socket/core/context"
	accountrepo "go-socket/core/domain/account/infra/persistent/repos"
	accountusecase "go-socket/core/domain/account/usecase"
	roomrepo "go-socket/core/domain/room/infra/persistent/repos"
	roomusecase "go-socket/core/domain/room/usecase"
	coreusecase "go-socket/core/usecase"
)

func BuildUsecase(appCtx *appCtx.AppContext) coreusecase.Usecase {
	accountRepos := accountrepo.NewRepoImpl(appCtx)
	authUC := accountusecase.NewAuthUsecase(appCtx, accountRepos)
	roomRepos := roomrepo.NewRepoImpl(appCtx)
	roomUC := roomusecase.NewRoomUsecase(appCtx, roomRepos)
	return coreusecase.NewUsecase(authUC, roomUC)
}
