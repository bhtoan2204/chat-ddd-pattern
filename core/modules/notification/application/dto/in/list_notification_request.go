package in

type ListNotificationRequest struct {
	Cursor string `json:"cursor" form:"cursor"`
	Limit  int    `json:"limit" form:"limit"`
}

func (r *ListNotificationRequest) Validate() error {
	return nil
}
