package socket

import (
	"context"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestHubRegisterDoesNotBlockAndStoresClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := NewMockIClient(ctrl)
	client.EXPECT().GetID().Return("client-1").AnyTimes()
	client.EXPECT().GetUserID().Return("user-1").AnyTimes()

	hub := &Hub{
		clients:       make(map[string]IClient),
		rooms:         make(map[string]IRoom),
		clientRooms:   make(map[string]map[string]struct{}),
		subscriptions: make(map[string]*roomSubscription),
	}

	done := make(chan struct{})
	go func() {
		hub.Register(context.Background(), client)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Register blocked")
	}

	got, ok := hub.clients["client-1"]
	if !ok {
		t.Fatalf("expected client %s to be registered", "client-1")
	}
	if got != client {
		t.Fatalf("expected registered client to match input client")
	}

	rooms, ok := hub.clientRooms["client-1"]
	if !ok {
		t.Fatalf("expected clientRooms entry for client %s", "client-1")
	}
	if len(rooms) != 0 {
		t.Fatalf("expected no joined rooms after register, got %d", len(rooms))
	}
}
