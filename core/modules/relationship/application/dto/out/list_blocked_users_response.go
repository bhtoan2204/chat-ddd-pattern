// CODE_GENERATOR - do not edit: response
package out

type ListBlockedUsersResponse struct {
	Items      []RelationshipAccountSummaryResponse `json:"items,omitempty"`
	NextCursor string                               `json:"next_cursor,omitempty"`
	Total      int64                                `json:"total,omitempty"`
}
