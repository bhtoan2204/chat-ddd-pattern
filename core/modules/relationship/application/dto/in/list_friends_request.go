// CODE_GENERATOR - do not edit: request

package in

import (
	"strings"
)

type ListFriendsRequest struct {
	UserID string `json:"user_id" form:"user_id"`
	Cursor string `json:"cursor" form:"cursor"`
	Limit  int    `json:"limit" form:"limit"`
}

func (r *ListFriendsRequest) Normalize() {
	r.UserID = strings.TrimSpace(r.UserID)
	r.Cursor = strings.TrimSpace(r.Cursor)
}

func (r *ListFriendsRequest) Validate() error {
	r.Normalize()
	return nil
}
