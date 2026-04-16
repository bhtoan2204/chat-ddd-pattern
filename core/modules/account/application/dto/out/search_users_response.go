// CODE_GENERATOR - do not edit: response
package out

type SearchUsersResponse struct {
	Total  int64            `json:"total,omitempty"`
	Limit  int              `json:"limit,omitempty"`
	Offset int              `json:"offset,omitempty"`
	Items  []SearchUserItem `json:"items,omitempty"`
}

type SearchUserItem struct {
	ID              string `json:"id,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	Username        string `json:"username,omitempty"`
	AvatarObjectKey string `json:"avatar_object_key,omitempty"`
	Status          string `json:"status,omitempty"`
	EmailVerified   bool   `json:"email_verified,omitempty"`
}
