// CODE_GENERATOR - do not edit: response
package out

type ListFollowingResponse struct {
	Items      []RelationshipAccountSummaryResponse `json:"items,omitempty"`
	NextCursor string                               `json:"next_cursor,omitempty"`
	Total      int64                                `json:"total,omitempty"`
}
