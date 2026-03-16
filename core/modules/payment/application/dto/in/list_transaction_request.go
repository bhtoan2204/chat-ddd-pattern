package in

type ListTransactionRequest struct {
	Page  int64 `json:"page" form:"page"`
	Limit int64 `json:"limit" form:"limit"`
}

func (r *ListTransactionRequest) Validate() error {
	if r.Page < 0 {
		r.Page = 0
	}
	if r.Limit < 0 {
		r.Limit = 10
	}

	return nil
}
