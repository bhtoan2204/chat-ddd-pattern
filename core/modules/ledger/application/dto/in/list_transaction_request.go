// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"

	"go-socket/core/shared/pkg/stackErr"
)

type ListTransactionRequest struct {
	Cursor   string `json:"cursor" form:"cursor"`
	Limit    int    `json:"limit" form:"limit"`
	Currency string `json:"currency" form:"currency"`
}

func (r *ListTransactionRequest) Normalize() {
	r.Cursor = strings.TrimSpace(r.Cursor)
	r.Currency = strings.ToUpper(strings.TrimSpace(r.Currency))
}

func (r *ListTransactionRequest) Validate() error {
	r.Normalize()
	if r.Limit < 0 {
		return stackErr.Error(errors.New("limit must be greater than or equal to zero"))
	}
	return nil
}
