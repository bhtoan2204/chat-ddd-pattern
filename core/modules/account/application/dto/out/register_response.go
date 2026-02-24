// CODE_GENERATOR: response
package out

type RegisterResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}
