package in

import (
	"errors"
	"go-socket/core/shared/pkg/webpush"
)

type SavePushSubscriptionRequest struct {
	Endpoint string       `json:"endpoint" form:"endpoint" binding:"required"`
	Keys     webpush.Keys `json:"keys" form:"keys" binding:"required"`
}

func (r *SavePushSubscriptionRequest) Validate() error {
	if r.Endpoint == "" {
		return errors.New("endpoint is required")
	}
	if r.Keys.Auth == "" {
		return errors.New("keys.auth is required")
	}
	if r.Keys.P256dh == "" {
		return errors.New("keys.p256dh is required")
	}
	return nil
}
