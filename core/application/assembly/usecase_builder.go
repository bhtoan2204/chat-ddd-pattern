package assembly

import (
	coreusecase "go-socket/core/application/usecase"
	appCtx "go-socket/core/context"
	accountrepo "go-socket/core/domain/account/infra/persistent/repository"
	accountusecase "go-socket/core/domain/account/usecase"
	roomrepo "go-socket/core/domain/room/infra/persistent/repository"
	roomusecase "go-socket/core/domain/room/usecase"
)

func BuildUsecase(appCtx *appCtx.AppContext) coreusecase.Usecase {
	accountRepos := accountrepo.NewRepoImpl(appCtx)
	authUC := accountusecase.NewAuthUsecase(appCtx, accountRepos)
	roomRepos := roomrepo.NewRepoImpl(appCtx)
	roomUC := roomusecase.NewRoomUsecase(appCtx, roomRepos)
	messageUC := roomusecase.NewMessageUsecase(appCtx, roomRepos)
	return coreusecase.NewUsecase(authUC, roomUC, messageUC)
}
