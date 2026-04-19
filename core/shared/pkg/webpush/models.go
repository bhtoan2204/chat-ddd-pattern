package webpush

import (
	"strings"

	lib "github.com/SherClockHolmes/webpush-go"
)

type Keys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
}

type Subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     Keys   `json:"keys"`
}

func (s Subscription) Validate() error {
	if strings.TrimSpace(s.Endpoint) == "" {
		return ErrSubscriptionEndpointRequired
	}
	if strings.TrimSpace(s.Keys.Auth) == "" || strings.TrimSpace(s.Keys.P256dh) == "" {
		return ErrSubscriptionKeysRequired
	}
	return nil
}

func (s Subscription) ToLibSubscription() *lib.Subscription {
	return &lib.Subscription{
		Endpoint: s.Endpoint,
		Keys: lib.Keys{
			Auth:   s.Keys.Auth,
			P256dh: s.Keys.P256dh,
		},
	}
}
