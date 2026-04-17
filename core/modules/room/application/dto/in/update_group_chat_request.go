// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type UpdateGroupChatRequest struct {
	RoomID      string `json:"room_id" form:"room_id" binding:"required"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}

func (r *UpdateGroupChatRequest) Normalize() {
	r.RoomID = strings.TrimSpace(r.RoomID)
	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)
}

func (r *UpdateGroupChatRequest) Validate() error {
	r.Normalize()
	if r.RoomID == "" {
		return stackErr.Error(errors.New("room_id is required"))
	}
	return nil
}
