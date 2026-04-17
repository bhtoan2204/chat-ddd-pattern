// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type GetChatPresenceRequest struct {
	AccountID string `json:"account_id" form:"account_id" binding:"required"`
}

func (r *GetChatPresenceRequest) Normalize() {
	r.AccountID = strings.TrimSpace(r.AccountID)
}

func (r *GetChatPresenceRequest) Validate() error {
	r.Normalize()
	if r.AccountID == "" {
		return stackErr.Error(errors.New("account_id is required"))
	}
	return nil
}
