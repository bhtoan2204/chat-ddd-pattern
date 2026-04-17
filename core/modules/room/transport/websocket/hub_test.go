package socket

import (
	"context"
	"testing"
	"time"
)

type stubClient struct {
	id     string
	userID string
}

func (c *stubClient) GetID() string {
	return c.id
}

func (c *stubClient) GetUserID() string {
	return c.userID
}

func (c *stubClient) Send(context.Context, []byte) {}

func (c *stubClient) ReadPump(context.Context, IHub) {}

func (c *stubClient) WritePump(context.Context) {}

func (c *stubClient) Close(context.Context) {}

func TestHubRegisterDoesNotDeadlockWhenAutoJoiningUserChannel(t *testing.T) {
	hub := &Hub{
		clients:       make(map[string]IClient),
		rooms:         make(map[string]IRoom),
		clientRooms:   make(map[string]map[string]struct{}),
		subscriptions: make(map[string]*roomSubscription),
	}

	client := &stubClient{
		id:     "client-1",
		userID: "user-1",
	}

	done := make(chan struct{})
	go func() {
		hub.Register(context.Background(), client)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Register blocked while auto-joining user channel")
	}

	if _, ok := hub.clients[client.id]; !ok {
		t.Fatalf("expected client %s to be registered", client.id)
	}

	userChannel := userChannelName(client.userID)
	if _, ok := hub.rooms[userChannel]; !ok {
		t.Fatalf("expected auto-joined user channel %s", userChannel)
	}
	if _, ok := hub.clientRooms[client.id][userChannel]; !ok {
		t.Fatalf("expected client %s to be tracked in user channel %s", client.id, userChannel)
	}
}
