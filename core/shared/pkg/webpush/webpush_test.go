package webpush

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"wechat-clone/core/shared/config"

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

	service, err := NewWebPush(cfg)
	if err != nil {
		t.Fatalf("NewWebPush() error = %v", err)
	}
	impl, ok := service.(*webPush)
	if !ok {
		t.Fatalf("expected *webPush, got %T", service)
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

func TestWebPushSendSuccess(t *testing.T) {
	service := &webPush{
		vapidPublicKey:  "pub",
		vapidPrivateKey: "pri",
		ttl:             60,
	}

	subscription := Subscription{
		Endpoint: "https://example.com/push",
		Keys: Keys{
			Auth:   "auth",
			P256dh: "p256dh",
		},
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
	if gotSub == nil {
		t.Fatalf("expected mapped subscription")
	}
	if gotSub.Endpoint != subscription.Endpoint {
		t.Fatalf("expected endpoint %s, got %s", subscription.Endpoint, gotSub.Endpoint)
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

func TestWebPushSendError(t *testing.T) {
	service := &webPush{
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

	err := service.Send(context.Background(), []byte(`{}`), Subscription{
		Endpoint: "https://example.com/push",
		Keys: Keys{
			Auth:   "auth",
			P256dh: "p256dh",
		},
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestWebPushSendRejectsInvalidSubscription(t *testing.T) {
	service := &webPush{
		vapidPublicKey:  "pub",
		vapidPrivateKey: "pri",
		ttl:             30,
	}

	err := service.Send(context.Background(), []byte(`{}`), Subscription{})
	if !errors.Is(err, ErrSubscriptionEndpointRequired) {
		t.Fatalf("expected ErrSubscriptionEndpointRequired, got %v", err)
	}
}

func TestWebPushSendRejectsBadResponse(t *testing.T) {
	service := &webPush{
		vapidPublicKey:  "pub",
		vapidPrivateKey: "pri",
		ttl:             30,
	}

	original := sendNotification
	sendNotification = func(_ context.Context, _ []byte, _ *lib.Subscription, _ *lib.Options) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusGone,
			Body:       io.NopCloser(strings.NewReader("expired")),
		}, nil
	}
	t.Cleanup(func() {
		sendNotification = original
	})

	err := service.Send(context.Background(), []byte(`{}`), Subscription{
		Endpoint: "https://example.com/push",
		Keys: Keys{
			Auth:   "auth",
			P256dh: "p256dh",
		},
	})
	if !errors.Is(err, ErrPushDeliveryRejected) {
		t.Fatalf("expected ErrPushDeliveryRejected, got %v", err)
	}
}

func TestWebPushSendMany(t *testing.T) {
	service := &webPush{
		vapidPublicKey:  "pub",
		vapidPrivateKey: "pri",
		ttl:             30,
	}

	callCount := 0
	original := sendNotification
	sendNotification = func(_ context.Context, _ []byte, _ *lib.Subscription, _ *lib.Options) (*http.Response, error) {
		callCount++
		return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader(""))}, nil
	}
	t.Cleanup(func() {
		sendNotification = original
	})

	err := service.SendMany(context.Background(), []byte(`{}`), []Subscription{
		{Endpoint: "https://example.com/1", Keys: Keys{Auth: "a", P256dh: "b"}},
		{Endpoint: "https://example.com/2", Keys: Keys{Auth: "c", P256dh: "d"}},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 calls, got %d", callCount)
	}
}

func TestNewWebPushServiceReal(t *testing.T) {
	cfg := &config.Config{
		WebPushConfig: config.WebPushConfig{
			VAPIDPublicKey:  "BB4d5HKEUbNTcRsIgIIyHeFNJPSPxrLlI6Ywsh5FANjz5RGK6-Z1MgD3wIeb8s8wpJtZ10VW5EHLoRp2ogGqTOM",
			VAPIDPrivateKey: "rIh0fg2Af_n4_FihdYaHwrb7-GxUF90j3gcmRvRN3K4",
			TTL:             45,
		},
	}
	webpush, err := NewWebPush(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	type webPushPayload struct {
		Title string                 `json:"title"`
		Body  string                 `json:"body"`
		Data  map[string]interface{} `json:"data,omitempty"`
	}
	payload := webPushPayload{
		Title: "PayChat",
		Body:  "Ban co thong bao moi.",
		Data: map[string]interface{}{
			"notification_id": "test-id",
			"account_id":      "test-account-id",
			"url":             "localhost:5173",
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}
	if err := webpush.Send(context.Background(), data, Subscription{
		Endpoint: "https://fcm.googleapis.com/fcm/send/cMeCgEP6Tl4:APA91bG5LNFibwmRJ54ncy_DZHJf0hOK5jz45g0AW9bX8lJ6w3CFlV71pd_t4R5-z5N-a_HznWDSdj02clE95hYVTdjgcgNrEp5tZ7X_PbXMxnJHdAnrExn6GrLF3N262edCcLMDxSaj",
		Keys: Keys{
			Auth:   "Jph1s8Hk3/pY2RN8YYJETA==",
			P256dh: "BK+2MT8OW2nyqncuW3/vv/rS3l0FtFOQ+DSv7K+ZJZEeHvfSYI96kUy74aW8c6MQVrjTinro13SvsnA4BYpAzAE=",
		},
	}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

}
