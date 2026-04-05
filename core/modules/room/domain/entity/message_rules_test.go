package entity

import (
	"errors"
	"testing"
	"time"
)

func TestNewMessageValidatesMediaPayload(t *testing.T) {
	_, err := NewMessage("msg-1", "room-1", "user-1", MessageParams{
		MessageType: MessageTypeImage,
	}, time.Now().UTC())
	if !errors.Is(err, ErrMessageObjectKeyRequired) {
		t.Fatalf("expected object key error, got %v", err)
	}
}

func TestMessageEditAndDeleteRules(t *testing.T) {
	message, err := NewMessage("msg-1", "room-1", "user-1", MessageParams{
		Message:     "hello",
		MessageType: MessageTypeText,
	}, time.Now().UTC())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := message.Edit("user-2", "updated", time.Now().UTC()); !errors.Is(err, ErrMessageCannotEditOther) {
		t.Fatalf("expected cannot edit other error, got %v", err)
	}
	if err := message.Edit("user-1", "updated", time.Now().UTC()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if message.Message != "updated" || message.EditedAt == nil {
		t.Fatalf("expected edited message state, got %+v", message)
	}

	if err := message.DeleteForEveryone("user-2", time.Now().UTC()); !errors.Is(err, ErrMessageCannotDeleteEveryone) {
		t.Fatalf("expected cannot delete error, got %v", err)
	}
	if err := message.DeleteForEveryone("user-1", time.Now().UTC()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if message.Message != "" || message.DeletedForEveryoneAt == nil {
		t.Fatalf("expected deleted message state, got %+v", message)
	}
}
