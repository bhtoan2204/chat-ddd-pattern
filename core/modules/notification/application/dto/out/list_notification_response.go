package out

import "time"

type ListNotificationResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	NextCursor    string                 `json:"next_cursor"`
	HasMore       bool                   `json:"has_more"`
	Total         int                    `json:"total"`
	Limit         int                    `json:"limit"`
}

type NotificationResponse struct {
	ID        string     `json:"id"`
	AccountID string     `json:"account_id"`
	Type      string     `json:"type"`
	Subject   string     `json:"subject"`
	Body      string     `json:"body"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}
