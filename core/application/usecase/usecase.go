package usecase

import (
	accountusecase "go-socket/core/domain/account/usecase"
	roomusecase "go-socket/core/domain/room/usecase"
)

type AuthUsecase = accountusecase.AuthUsecase
type RoomUsecase = roomusecase.RoomUsecase
type MessageUsecase = roomusecase.MessageUsecase

type Usecase interface {
	AuthUsecase() accountusecase.AuthUsecase
	RoomUsecase() roomusecase.RoomUsecase
	MessageUsecase() roomusecase.MessageUsecase
}

type usecase struct {
	// Account module
	auth accountusecase.AuthUsecase

	// Room module
	room    roomusecase.RoomUsecase
	message roomusecase.MessageUsecase
}

func NewUsecase(
	auth accountusecase.AuthUsecase,
	room roomusecase.RoomUsecase,
	message roomusecase.MessageUsecase,
) Usecase {
	return &usecase{
		auth:    auth,
		room:    room,
		message: message,
	}
}

func (u *usecase) AuthUsecase() accountusecase.AuthUsecase {
	return u.auth
}

func (u *usecase) RoomUsecase() roomusecase.RoomUsecase {
	return u.room
}

func (u *usecase) MessageUsecase() roomusecase.MessageUsecase {
	return u.message
}
