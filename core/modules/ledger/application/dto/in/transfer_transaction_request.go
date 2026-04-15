// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"go-socket/core/shared/pkg/stackErr"
	"strings"
)

type TransferTransactionRequest struct {
	ToAccountID string            `json:"to_account_id" form:"to_account_id" binding:"required"`
	Currency    string            `json:"currency" form:"currency" binding:"required"`
	Amount      int64             `json:"amount" form:"amount" binding:"required"`
	Metadata    map[string]string `json:"metadata" form:"metadata"`
}

func (r *TransferTransactionRequest) Normalize() {
	r.ToAccountID = strings.TrimSpace(r.ToAccountID)
	r.Currency = strings.TrimSpace(r.Currency)
	for key, value := range r.Metadata {
		r.Metadata[key] = strings.TrimSpace(value)
	}
}

func (r *TransferTransactionRequest) Validate() error {
	r.Normalize()
	if r.ToAccountID == "" {
		return stackErr.Error(errors.New("to_account_id is required"))
	}
	if r.Currency == "" {
		return stackErr.Error(errors.New("currency is required"))
	}
	if r.Amount == 0 {
		return stackErr.Error(errors.New("amount is required"))
	}
	return nil
}
