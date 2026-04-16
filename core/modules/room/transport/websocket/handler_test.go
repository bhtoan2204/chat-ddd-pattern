package socket

import (
	"context"
	"testing"

	"go-socket/core/shared/pkg/pubsub"

	"go.uber.org/mock/gomock"
)

func TestHandleRealtimeMessageIgnoresUnknownTopic(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	hub := NewMockIHub(ctrl)
	if err := handleRealtimeMessage(context.Background(), hub, pubsub.Message{
		Topic: "unknown",
		Data:  "ignored",
	}); err != nil {
		t.Fatalf("handleRealtimeMessage returned error: %v", err)
	}

}
