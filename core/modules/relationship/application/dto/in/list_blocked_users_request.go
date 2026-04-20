// CODE_GENERATOR - do not edit: request

package in

import (
	"strings"
)

type ListBlockedUsersRequest struct {
	Cursor string `json:"cursor" form:"cursor"`
	Limit  int    `json:"limit" form:"limit"`
}

func (r *ListBlockedUsersRequest) Normalize() {
	r.Cursor = strings.TrimSpace(r.Cursor)
}

func (r *ListBlockedUsersRequest) Validate() error {
	r.Normalize()
	return nil
}
