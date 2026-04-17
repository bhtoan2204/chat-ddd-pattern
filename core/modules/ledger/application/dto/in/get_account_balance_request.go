// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type GetAccountBalanceRequest struct {
	AccountID string `json:"account_id" form:"account_id" binding:"required"`
	Currency  string `json:"currency" form:"currency" binding:"required"`
}

func (r *GetAccountBalanceRequest) Normalize() {
	r.AccountID = strings.TrimSpace(r.AccountID)
	r.Currency = strings.TrimSpace(r.Currency)
}

func (r *GetAccountBalanceRequest) Validate() error {
	r.Normalize()
	if r.AccountID == "" {
		return stackErr.Error(errors.New("account_id is required"))
	}
	if r.Currency == "" {
		return stackErr.Error(errors.New("currency is required"))
	}
	return nil
}
