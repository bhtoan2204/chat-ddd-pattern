package command

import "errors"

var (
	ErrRoomFull            = errors.New("room is full")
	ErrRoomNotFound        = errors.New("room not found")
	ErrRoomAlreadyJoined   = errors.New("room already joined")
	ErrRoomAccountNotFound = errors.New("account not found")
)
