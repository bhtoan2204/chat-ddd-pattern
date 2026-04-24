// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type CreateQuoteRequest struct {
	FromCurrency string `json:"from_currency" form:"from_currency" binding:"required"`
	ToCurrency   string `json:"to_currency" form:"to_currency" binding:"required"`
	ToAmount     string `json:"to_amount" form:"to_amount" binding:"required"`
	Purpose      string `json:"purpose" form:"purpose" binding:"required"`
}

func (r *CreateQuoteRequest) Normalize() {
	r.FromCurrency = strings.TrimSpace(r.FromCurrency)
	r.ToCurrency = strings.TrimSpace(r.ToCurrency)
	r.ToAmount = strings.TrimSpace(r.ToAmount)
	r.Purpose = strings.TrimSpace(r.Purpose)
}

func (r *CreateQuoteRequest) Validate() error {
	r.Normalize()
	if r.FromCurrency == "" {
		return stackErr.Error(errors.New("from_currency is required"))
	}
	if r.ToCurrency == "" {
		return stackErr.Error(errors.New("to_currency is required"))
	}
	if r.ToAmount == "" {
		return stackErr.Error(errors.New("to_amount is required"))
	}
	if r.Purpose == "" {
		return stackErr.Error(errors.New("purpose is required"))
	}
	return nil
}
