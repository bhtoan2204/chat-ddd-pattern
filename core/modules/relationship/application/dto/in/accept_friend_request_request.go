// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type AcceptFriendRequestRequest struct {
	RequesterUserID string `json:"requester_user_id" form:"requester_user_id" binding:"required"`
}

func (r *AcceptFriendRequestRequest) Normalize() {
	r.RequesterUserID = strings.TrimSpace(r.RequesterUserID)
}

func (r *AcceptFriendRequestRequest) Validate() error {
	r.Normalize()
	if r.RequesterUserID == "" {
		return stackErr.Error(errors.New("requester_user_id is required"))
	}
	return nil
}
