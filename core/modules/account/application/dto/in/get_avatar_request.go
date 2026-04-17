// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type GetAvatarRequest struct {
	AccountID string `json:"account_id" form:"account_id" binding:"required"`
}

func (r *GetAvatarRequest) Normalize() {
	r.AccountID = strings.TrimSpace(r.AccountID)
}

func (r *GetAvatarRequest) Validate() error {
	r.Normalize()
	if r.AccountID == "" {
		return stackErr.Error(errors.New("account_id is required"))
	}
	return nil
}
