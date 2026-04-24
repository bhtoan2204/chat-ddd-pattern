// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type CreateStripeTopUpRequest struct {
	Amount   int64             `json:"amount" form:"amount" binding:"required"`
	Currency string            `json:"currency" form:"currency" binding:"required"`
	Metadata map[string]string `json:"metadata" form:"metadata"`
}

func (r *CreateStripeTopUpRequest) Normalize() {
	r.Currency = strings.TrimSpace(r.Currency)
	for key, value := range r.Metadata {
		r.Metadata[key] = strings.TrimSpace(value)
	}
}

func (r *CreateStripeTopUpRequest) Validate() error {
	r.Normalize()
	if r.Amount == 0 {
		return stackErr.Error(errors.New("amount is required"))
	}
	if r.Currency == "" {
		return stackErr.Error(errors.New("currency is required"))
	}
	return nil
}
