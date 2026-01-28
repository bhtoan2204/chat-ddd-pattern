package entity

import (
	"go-socket/types"
	"time"
)

type Room struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	OwnerID     string         `json:"owner_id"`
	RoomType    types.RoomType `json:"room_type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
