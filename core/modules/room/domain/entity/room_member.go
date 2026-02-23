package entity

import (
	"go-socket/core/modules/room/types"
	"time"
)

type RoomMemberEntity struct {
	ID        string         `json:"id"`
	RoomID    string         `json:"room_id"`
	AccountID string         `json:"account_id"`
	Role      types.RoomRole `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
