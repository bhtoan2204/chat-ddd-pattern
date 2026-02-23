package assembly

import (
	appCtx "go-socket/core/context"
	roomusecase "go-socket/core/modules/room/application/usecase"
	roomrepo "go-socket/core/modules/room/infra/persistent/repository"
)

type Usecases struct {
	Room    roomusecase.RoomUsecase
	Message roomusecase.MessageUsecase
}

func BuildUsecases(appCtx *appCtx.AppContext) Usecases {
	roomRepos := roomrepo.NewRepoImpl(appCtx)
	return Usecases{
		Room:    roomusecase.NewRoomUsecase(appCtx, roomRepos),
		Message: roomusecase.NewMessageUsecase(appCtx, roomRepos),
	}
}
