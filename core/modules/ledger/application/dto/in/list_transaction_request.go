// CODE_GENERATOR - do not edit: request

package in

import (
	"strings"
)

type ListTransactionRequest struct {
	Cursor   string `json:"cursor" form:"cursor"`
	Limit    int    `json:"limit" form:"limit"`
	Currency string `json:"currency" form:"currency"`
}

func (r *ListTransactionRequest) Normalize() {
	r.Cursor = strings.TrimSpace(r.Cursor)
	r.Currency = strings.TrimSpace(r.Currency)
}

func (r *ListTransactionRequest) Validate() error {
	r.Normalize()
	return nil
}
