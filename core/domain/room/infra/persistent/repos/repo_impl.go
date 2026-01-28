package repository

import (
	appCtx "go-socket/core/context"
	"go-socket/core/domain/room/repos"
)

type repoImpl struct {
	roomRepo repos.RoomRepository
}

func NewRepoImpl(appCtx *appCtx.AppContext) repos.Repos {
	roomRepo := NewRoomRepoImpl(appCtx.GetDB(), appCtx.GetCache())
	return &repoImpl{roomRepo: roomRepo}
}

func (r *repoImpl) RoomRepository() repos.RoomRepository {
	return r.roomRepo
}
