package usecase

import (
	appCtx "go-socket/core/context"
	accountrepo "go-socket/core/domain/account/infra/persistent/repos"
	accountusecase "go-socket/core/domain/account/usecase"
	roomrepo "go-socket/core/domain/room/infra/persistent/repos"
	roomusecase "go-socket/core/domain/room/usecase"
)

type AuthUsecase = accountusecase.AuthUsecase
type RoomUsecase = roomusecase.RoomUsecase

type Usecase interface {
	AuthUsecase() accountusecase.AuthUsecase
	RoomUsecase() roomusecase.RoomUsecase
}

type usecase struct {
	auth accountusecase.AuthUsecase
	room roomusecase.RoomUsecase
}

func NewUsecase(appCtx *appCtx.AppContext) Usecase {
	accountRepos := accountrepo.NewRepoImpl(appCtx)
	authUC := accountusecase.NewAuthUsecase(appCtx, accountRepos)
	roomRepos := roomrepo.NewRepoImpl(appCtx)
	roomUC := roomusecase.NewRoomUsecase(appCtx, roomRepos)
	return &usecase{
		auth: authUC,
		room: roomUC,
	}
}

func (u *usecase) AuthUsecase() accountusecase.AuthUsecase {
	return u.auth
}

func (u *usecase) RoomUsecase() roomusecase.RoomUsecase {
	return u.room
}
