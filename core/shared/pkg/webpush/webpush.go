package webpush

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"

	lib "github.com/SherClockHolmes/webpush-go"
)

var sendNotification = func(
	ctx context.Context,
	payload []byte,
	subscription *lib.Subscription,
	options *lib.Options,
) (*http.Response, error) {
	return lib.SendNotificationWithContext(ctx, payload, subscription, options)
}

var (
	ErrInvalidConfig                = errors.New("webpush configuration is invalid")
	ErrSubscriptionEndpointRequired = errors.New("webpush subscription endpoint is required")
	ErrSubscriptionKeysRequired     = errors.New("webpush subscription keys are required")
	ErrPushDeliveryRejected         = errors.New("webpush delivery rejected")
)

//go:generate mockgen -package=webpush -destination=webpush_mock.go -source=webpush.go
type WebPush interface {
	Send(ctx context.Context, payload []byte, subscription Subscription) error
	SendMany(ctx context.Context, payload []byte, subscriptions []Subscription) error
}

type webPush struct {
	vapidPublicKey  string
	vapidPrivateKey string
	ttl             int
}

func NewWebPush(cfg *config.Config) (WebPush, error) {
	if cfg == nil {
		return nil, stackErr.Error(ErrInvalidConfig)
	}
	if strings.TrimSpace(cfg.WebPushConfig.VAPIDPublicKey) == "" || strings.TrimSpace(cfg.WebPushConfig.VAPIDPrivateKey) == "" {
		return nil, stackErr.Error(ErrInvalidConfig)
	}

	ttl := cfg.WebPushConfig.TTL
	if ttl <= 0 {
		ttl = 30
	}

	return &webPush{
		vapidPublicKey:  cfg.WebPushConfig.VAPIDPublicKey,
		vapidPrivateKey: cfg.WebPushConfig.VAPIDPrivateKey,
		ttl:             ttl,
	}, nil
}

func (s *webPush) Send(
	ctx context.Context,
	payload []byte,
	subscription Subscription,
) error {
	if err := subscription.Validate(); err != nil {
		return stackErr.Error(err)
	}

	resp, err := sendNotification(ctx, payload, subscription.ToLibSubscription(), &lib.Options{
		VAPIDPublicKey:  s.vapidPublicKey,
		VAPIDPrivateKey: s.vapidPrivateKey,
		TTL:             s.ttl,
	})
	if err != nil {
		return stackErr.Error(err)
	}

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp != nil && resp.StatusCode >= http.StatusBadRequest {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return stackErr.Error(fmt.Errorf("%w: status=%d body_read_failed=%v", ErrPushDeliveryRejected, resp.StatusCode, readErr))
		}
		return stackErr.Error(fmt.Errorf("%w: status=%d body=%s", ErrPushDeliveryRejected, resp.StatusCode, strings.TrimSpace(string(body))))
	}

	return nil
}

func (s *webPush) SendMany(
	ctx context.Context,
	payload []byte,
	subscriptions []Subscription,
) error {
	for idx, subscription := range subscriptions {
		if err := s.Send(ctx, payload, subscription); err != nil {
			return stackErr.Error(fmt.Errorf("send webpush subscription #%d failed: %w", idx, err))
		}
	}

	return nil
}
