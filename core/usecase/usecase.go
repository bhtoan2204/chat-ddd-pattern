package usecase

import (
	accountusecase "go-socket/core/domain/account/usecase"
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

func NewUsecase(
	auth accountusecase.AuthUsecase,
	room roomusecase.RoomUsecase,
) Usecase {
	return &usecase{
		auth: auth,
		room: room,
	}
}

func (u *usecase) AuthUsecase() accountusecase.AuthUsecase {
	return u.auth
}

func (u *usecase) RoomUsecase() roomusecase.RoomUsecase {
	return u.room
}
