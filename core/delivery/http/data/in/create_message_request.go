// CODE_GENERATOR: request

package in

import "errors"

type CreateMessageRequest struct {
	RoomID  string `json:"room_id" form:"room_id"`
	Message string `json:"message" form:"message"`
}

func (r *CreateMessageRequest) Validate() error {
	if r.RoomID == "" {
		return errors.New("room_id is required")
	}
	if r.Message == "" {
		return errors.New("message is required")
	}
	return nil
}
