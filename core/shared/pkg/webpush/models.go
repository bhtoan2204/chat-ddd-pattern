package webpush

import lib "github.com/SherClockHolmes/webpush-go"

type Keys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
}

type Subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     Keys   `json:"keys"`
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
