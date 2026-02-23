// CODE_GENERATOR: request

package in

type ListRoomsRequest struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

func (r *ListRoomsRequest) Validate() error {
	return nil
}
