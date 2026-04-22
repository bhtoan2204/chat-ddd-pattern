// CODE_GENERATOR - do not edit: response
package out

type ListIncomingFriendRequestsResponse struct {
	Items      []RelationshipAccountSummaryResponse `json:"items,omitempty"`
	NextCursor string                               `json:"next_cursor,omitempty"`
}

type RelationshipAccountSummaryResponse struct {
	AccountID       string `json:"account_id,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	Username        string `json:"username,omitempty"`
	AvatarObjectKey string `json:"avatar_object_key,omitempty"`
}
