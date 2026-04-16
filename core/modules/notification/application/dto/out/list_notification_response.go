// CODE_GENERATOR - do not edit: response
package out

type ListNotificationResponse struct {
	Notifications []NotificationResponse `json:"notifications,omitempty"`
	NextCursor    string                 `json:"next_cursor,omitempty"`
	HasMore       bool                   `json:"has_more,omitempty"`
	Total         int                    `json:"total,omitempty"`
	Limit         int                    `json:"limit,omitempty"`
}

type NotificationResponse struct {
	ID        string `json:"id,omitempty"`
	AccountID string `json:"account_id,omitempty"`
	Type      string `json:"type,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Body      string `json:"body,omitempty"`
	IsRead    bool   `json:"is_read,omitempty"`
	ReadAt    string `json:"read_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
