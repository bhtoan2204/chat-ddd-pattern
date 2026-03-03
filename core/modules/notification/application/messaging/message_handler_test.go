package messaging

import (
	"testing"
)

func TestParseAccountCreatedPayload_Object(t *testing.T) {
	raw := []byte(`{"AccountID":"acc-1","Email":"a@example.com","CreatedAt":"2026-03-03T13:05:32.218937909+07:00"}`)

	payload, err := parseAccountCreatedPayload(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if payload.AccountID != "acc-1" {
		t.Fatalf("expected account_id acc-1, got %s", payload.AccountID)
	}
	if payload.Email != "a@example.com" {
		t.Fatalf("expected email a@example.com, got %s", payload.Email)
	}
}

func TestParseAccountCreatedPayload_EncodedString(t *testing.T) {
	raw := []byte(`"{\"AccountID\":\"acc-2\",\"Email\":\"b@example.com\",\"CreatedAt\":\"2026-03-03T13:05:32.218937909+07:00\"}"`)

	payload, err := parseAccountCreatedPayload(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if payload.AccountID != "acc-2" {
		t.Fatalf("expected account_id acc-2, got %s", payload.AccountID)
	}
	if payload.Email != "b@example.com" {
		t.Fatalf("expected email b@example.com, got %s", payload.Email)
	}
}

func TestParseAccountCreatedPayload_Empty(t *testing.T) {
	_, err := parseAccountCreatedPayload([]byte(`""`))
	if err == nil {
		t.Fatalf("expected error when event_data is empty")
	}
}
