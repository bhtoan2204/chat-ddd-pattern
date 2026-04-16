// CODE_GENERATOR - do not edit: response
package out

type ListTransactionResponse struct {
	Limit      int                   `json:"limit"`
	Size       int                   `json:"size"`
	Total      int64                 `json:"total"`
	HasMore    bool                  `json:"has_more"`
	NextCursor string                `json:"next_cursor"`
	Records    []TransactionResponse `json:"records"`
}
