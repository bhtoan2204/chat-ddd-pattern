package aggregate

import (
	"go-socket/core/modules/room/types"
	"time"
)

type EventRoomCreated struct {
	RoomID              string
	RoomType            types.RoomType
	MemberCount         int
	LastMessageID       string
	LastMessageAt       time.Time
	LastMessageContent  string
	LastMessageSenderID string
}

type EventRoomMemberAdded struct {
	RoomID         string
	MemberID       string
	MemberName     string
	MemberEmail    string
	MemberRole     types.RoomRole
	MemberJoinedAt time.Time
}

type EventRoomMemberRemoved struct {
	RoomID         string
	MemberID       string
	MemberName     string
	MemberEmail    string
	MemberRole     types.RoomRole
	MemberJoinedAt time.Time
}

type EventRoomMessageCreated struct {
	RoomID             string
	MessageID          string
	MessageContent     string
	MessageSenderID    string
	MessageSenderName  string
	MessageSenderEmail string
	MessageSentAt      time.Time
}
