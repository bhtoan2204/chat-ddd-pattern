// CODE_GENERATOR - do not edit: response
package out

type ListTransactionResponse struct {
	Records    []TransactionResponse `json:"records,omitempty"`
	NextCursor string                `json:"next_cursor,omitempty"`
	HasMore    bool                  `json:"has_more,omitempty"`
	Total      int64                 `json:"total,omitempty"`
	Size       int                   `json:"size,omitempty"`
	Limit      int                   `json:"limit,omitempty"`
}
