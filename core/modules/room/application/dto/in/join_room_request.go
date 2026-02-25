package in

import "errors"

type JoinRoomRequest struct {
	RoomID string `json:"room_id" form:"room_id"`
}

func (r *JoinRoomRequest) Validate() error {
	if r.RoomID == "" {
		return errors.New("room_id is required")
	}
	return nil
}
