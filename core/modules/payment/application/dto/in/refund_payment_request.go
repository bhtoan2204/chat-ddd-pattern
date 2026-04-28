// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"

	"wechat-clone/core/shared/pkg/stackErr"
)

type RefundPaymentRequest struct {
	Provider      string `json:"provider" form:"provider" binding:"required"`
	TransactionID string `json:"transaction_id" form:"transaction_id" binding:"required"`
	Reason        string `json:"reason" form:"reason"`
}

func (r *RefundPaymentRequest) Normalize() {
	r.Provider = strings.TrimSpace(r.Provider)
	r.TransactionID = strings.TrimSpace(r.TransactionID)
	r.Reason = strings.TrimSpace(r.Reason)
}

func (r *RefundPaymentRequest) Validate() error {
	r.Normalize()
	if r.Provider == "" {
		return stackErr.Error(errors.New("provider is required"))
	}
	if r.TransactionID == "" {
		return stackErr.Error(errors.New("transaction_id is required"))
	}
	return nil
}
