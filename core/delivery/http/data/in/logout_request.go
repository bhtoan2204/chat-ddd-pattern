// CODE_GENERATOR: request

package in

type LogoutRequest struct {
	Token string `json:"token" form:"token"`
}

func (r *LogoutRequest) Validate() error {
	return nil
}
