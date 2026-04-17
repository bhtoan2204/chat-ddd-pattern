// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type CreateDirectConversationRequest struct {
	PeerAccountID string `json:"peer_account_id" form:"peer_account_id" binding:"required"`
}

func (r *CreateDirectConversationRequest) Normalize() {
	r.PeerAccountID = strings.TrimSpace(r.PeerAccountID)
}

func (r *CreateDirectConversationRequest) Validate() error {
	r.Normalize()
	if r.PeerAccountID == "" {
		return stackErr.Error(errors.New("peer_account_id is required"))
	}
	return nil
}
