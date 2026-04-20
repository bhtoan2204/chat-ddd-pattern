// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type FollowUserRequest struct {
	TargetUserID string `json:"target_user_id" form:"target_user_id" binding:"required"`
}

func (r *FollowUserRequest) Normalize() {
	r.TargetUserID = strings.TrimSpace(r.TargetUserID)
}

func (r *FollowUserRequest) Validate() error {
	r.Normalize()
	if r.TargetUserID == "" {
		return stackErr.Error(errors.New("target_user_id is required"))
	}
	return nil
}
