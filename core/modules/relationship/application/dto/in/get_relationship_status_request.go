// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type GetRelationshipStatusRequest struct {
	TargetUserID string `json:"target_user_id" form:"target_user_id" binding:"required"`
}

func (r *GetRelationshipStatusRequest) Normalize() {
	r.TargetUserID = strings.TrimSpace(r.TargetUserID)
}

func (r *GetRelationshipStatusRequest) Validate() error {
	r.Normalize()
	if r.TargetUserID == "" {
		return stackErr.Error(errors.New("target_user_id is required"))
	}
	return nil
}
