package types

import roomtypes "wechat-clone/core/modules/room/types"

type CreateRoomCommand struct {
	Name        string
	Description string
	RoomType    roomtypes.RoomType
}

type UpdateRoomCommand struct {
	Name        string
	Description string
	RoomType    roomtypes.RoomType
}

type JoinRoomCommand struct {
	RoomID string
}
