package repository

import (
	appCtx "go-socket/core/context"
	"go-socket/core/modules/room/domain/repos"
)

type repoImpl struct {
	roomRepo       repos.RoomRepository
	messageRepo    repos.MessageRepository
	roomMemberRepo repos.RoomMemberRepository
}

func NewRepoImpl(appCtx *appCtx.AppContext) repos.Repos {
	db := appCtx.GetDB()
	roomRepo := NewRoomRepoImpl(db, appCtx.GetCache())
	messageRepo := NewMessageRepoImpl(db)
	roomMemberRepo := NewRoomMemberImpl(db)
	return &repoImpl{
		roomRepo:       roomRepo,
		messageRepo:    messageRepo,
		roomMemberRepo: roomMemberRepo,
	}
}

func (r *repoImpl) RoomRepository() repos.RoomRepository {
	return r.roomRepo
}

func (r *repoImpl) MessageRepository() repos.MessageRepository {
	return r.messageRepo
}

func (r *repoImpl) RoomMemberRepository() repos.RoomMemberRepository {
	return r.roomMemberRepo
}
