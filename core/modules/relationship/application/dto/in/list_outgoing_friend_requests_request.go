// CODE_GENERATOR - do not edit: request

package in

import (
	"strings"
)

type ListOutgoingFriendRequestsRequest struct {
	Cursor string `json:"cursor" form:"cursor"`
	Limit  int    `json:"limit" form:"limit"`
}

func (r *ListOutgoingFriendRequestsRequest) Normalize() {
	r.Cursor = strings.TrimSpace(r.Cursor)
}

func (r *ListOutgoingFriendRequestsRequest) Validate() error {
	r.Normalize()
	return nil
}
