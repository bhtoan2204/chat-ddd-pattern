package webpush

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"go-socket/core/shared/config"

	lib "github.com/SherClockHolmes/webpush-go"
)

type spyReadCloser struct {
	closed bool
}

func (s *spyReadCloser) Read(_ []byte) (int, error) {
	return 0, io.EOF
}

func (s *spyReadCloser) Close() error {
	s.closed = true
	return nil
}

func TestNewWebPushService(t *testing.T) {
	cfg := &config.Config{
		WebPushConfig: config.WebPushConfig{
			VAPIDPublicKey:  "pub-key",
			VAPIDPrivateKey: "pri-key",
			TTL:             120,
		},
	}

	service := NewWebPushService(cfg)
	impl, ok := service.(*webPushService)
	if !ok {
		t.Fatalf("expected *webPushService, got %T", service)
	}

	if impl.vapidPublicKey != "pub-key" {
		t.Fatalf("expected vapidPublicKey pub-key, got %s", impl.vapidPublicKey)
	}
	if impl.vapidPrivateKey != "pri-key" {
		t.Fatalf("expected vapidPrivateKey pri-key, got %s", impl.vapidPrivateKey)
	}
	if impl.ttl != 120 {
		t.Fatalf("expected ttl 120, got %d", impl.ttl)
	}
}

func TestWebPushServiceSendSuccess(t *testing.T) {
	service := &webPushService{
		vapidPublicKey:  "pub",
		vapidPrivateKey: "pri",
		ttl:             60,
	}

	subscription := &lib.Subscription{
		Endpoint: "https://example.com/push",
	}
	payload := []byte(`{"title":"hello"}`)

	var gotPayload []byte
	var gotSub *lib.Subscription
	var gotOptions *lib.Options
	body := &spyReadCloser{}

	original := sendNotification
	sendNotification = func(_ context.Context, p []byte, s *lib.Subscription, o *lib.Options) (*http.Response, error) {
		gotPayload = p
		gotSub = s
		gotOptions = o
		return &http.Response{StatusCode: http.StatusCreated, Body: body}, nil
	}
	t.Cleanup(func() {
		sendNotification = original
	})

	if err := service.Send(context.Background(), payload, subscription); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(gotPayload) != string(payload) {
		t.Fatalf("expected payload %s, got %s", string(payload), string(gotPayload))
	}
	if gotSub != subscription {
		t.Fatalf("expected same subscription pointer")
	}
	if gotOptions == nil {
		t.Fatalf("expected options to be passed")
	}
	if gotOptions.VAPIDPublicKey != "pub" {
		t.Fatalf("expected VAPIDPublicKey pub, got %s", gotOptions.VAPIDPublicKey)
	}
	if gotOptions.VAPIDPrivateKey != "pri" {
		t.Fatalf("expected VAPIDPrivateKey pri, got %s", gotOptions.VAPIDPrivateKey)
	}
	if gotOptions.TTL != 60 {
		t.Fatalf("expected TTL 60, got %d", gotOptions.TTL)
	}
	if !body.closed {
		t.Fatalf("expected response body to be closed")
	}
}

func TestWebPushServiceSendError(t *testing.T) {
	service := &webPushService{
		vapidPublicKey:  "pub",
		vapidPrivateKey: "pri",
		ttl:             30,
	}

	expectedErr := errors.New("send failed")
	original := sendNotification
	sendNotification = func(_ context.Context, _ []byte, _ *lib.Subscription, _ *lib.Options) (*http.Response, error) {
		return nil, expectedErr
	}
	t.Cleanup(func() {
		sendNotification = original
	})

	err := service.Send(context.Background(), []byte(`{}`), &lib.Subscription{Endpoint: "https://example.com/push"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
