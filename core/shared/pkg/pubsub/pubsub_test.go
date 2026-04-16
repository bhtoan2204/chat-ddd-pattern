package pubsub

import (
	"context"
	"testing"
	"time"
)

func TestNew_DefaultBufferSize(t *testing.T) {
	bus := New(Config{})
	if bus == nil {
		t.Fatal("expected bus to be created")
	}

	sub, err := bus.Subscribe("topic-1")
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}

	if got, want := cap(sub.ch), 16; got != want {
		t.Fatalf("unexpected default buffer size: got=%d want=%d", got, want)
	}
}

func TestShoot_Msg(t *testing.T) {
	bus := New(Config{
		BufferSize:  1,
		PublishMode: PublishBlocking,
	})

	sub, err := bus.Subscribe("topic-1")
	if err != nil {
		t.Fatalf("subscribe sub1 failed: %v", err)
	}

	payload := 123
	if err := bus.Publish(context.Background(), "topic-1", payload); err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	select {
	case msg, ok := <-sub.C():
		if !ok {
			t.Fatalf("channel closed unexpectedly")
		}
		if msg.Topic != "topic-1" {
			t.Fatalf("%s unexpected topic", msg.Topic)
		}
		t.Log(msg)
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for message")
	}
}

func TestPublish_ToMultipleSubscribers(t *testing.T) {
	bus := New(Config{
		BufferSize:  1,
		PublishMode: PublishBlocking,
	})

	sub1, err := bus.Subscribe("topic-1")
	if err != nil {
		t.Fatalf("subscribe sub1 failed: %v", err)
	}
	sub2, err := bus.Subscribe("topic-1")
	if err != nil {
		t.Fatalf("subscribe sub2 failed: %v", err)
	}

	payload := 123
	if err := bus.Publish(context.Background(), "topic-1", payload); err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	assertReceive := func(name string, sub *Subscription) {
		t.Helper()

		select {
		case msg, ok := <-sub.C():
			if !ok {
				t.Fatalf("%s channel closed unexpectedly", name)
			}
			if msg.Topic != "topic-1" {
				t.Fatalf("%s unexpected topic: got=%s", name, msg.Topic)
			}
			if got, want := msg.Data, any(payload); got != want {
				t.Fatalf("%s unexpected payload: got=%v want=%v", name, got, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("%s timeout waiting for message", name)
		}
	}

	assertReceive("sub1", sub1)
	assertReceive("sub2", sub2)
}
