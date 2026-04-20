// CODE_GENERATOR - do not edit: request

package in

import (
	"strings"
)

type ListIncomingFriendRequestsRequest struct {
	Cursor string `json:"cursor" form:"cursor"`
	Limit  int    `json:"limit" form:"limit"`
}

func (r *ListIncomingFriendRequestsRequest) Normalize() {
	r.Cursor = strings.TrimSpace(r.Cursor)
}

func (r *ListIncomingFriendRequestsRequest) Validate() error {
	r.Normalize()
	return nil
}
