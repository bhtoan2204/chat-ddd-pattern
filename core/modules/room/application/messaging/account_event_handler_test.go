package messaging

import (
	"context"
	"testing"
	"time"

	"go-socket/core/modules/room/domain/entity"
	sharedevents "go-socket/core/shared/contracts/events"
)

type roomAccountProjectionRepoStub struct {
	projected *entity.AccountEntity
}

func (s *roomAccountProjectionRepoStub) ProjectAccount(_ context.Context, account *entity.AccountEntity) error {
	s.projected = account
	return nil
}

func (s *roomAccountProjectionRepoStub) ListByAccountIDs(context.Context, []string) ([]*entity.AccountEntity, error) {
	return nil, nil
}

func TestDecodeAccountCreatedPayloadUsesSharedContract(t *testing.T) {
	raw := []byte(`{"AccountID":"acc-1","Email":"a@example.com","DisplayName":"Alice","CreatedAt":"2026-03-03T13:05:32.218937909+07:00"}`)

	payloadAny, err := decodeEventPayload(context.Background(), sharedevents.EventAccountCreated, raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	payload, ok := payloadAny.(*sharedevents.AccountCreatedEvent)
	if !ok {
		t.Fatalf("expected AccountCreatedEvent, got %T", payloadAny)
	}

	if payload.AccountID != "acc-1" {
		t.Fatalf("expected account_id acc-1, got %s", payload.AccountID)
	}
	if payload.DisplayName != "Alice" {
		t.Fatalf("expected display name Alice, got %s", payload.DisplayName)
	}
}

func TestHandleAccountEventProfileUpdatedProjectsUsernameAndAvatar(t *testing.T) {
	repo := &roomAccountProjectionRepoStub{}
	handler := &messageHandler{accountRepo: repo}

	raw := []byte(`{
		"id": 1,
		"aggregate_id": "acc-3",
		"aggregate_type": "account",
		"version": 2,
		"event_name": "EventAccountProfileUpdated",
		"event_data": {
			"AccountID":"acc-3",
			"DisplayName":"Alice Updated",
			"Username":"alice",
			"AvatarObjectKey":"avatars/alice.png",
			"UpdatedAt":"2026-03-03T13:05:32.218937909+07:00"
		},
		"created_at": "2026-03-03T13:05:32.218937909+07:00"
	}`)

	if err := handler.handleAccountEvent(context.Background(), raw); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.projected == nil {
		t.Fatalf("expected projected account to be saved")
	}
	if repo.projected.AccountID != "acc-3" {
		t.Fatalf("expected account_id acc-3, got %s", repo.projected.AccountID)
	}
	if repo.projected.DisplayName != "Alice Updated" {
		t.Fatalf("expected display name Alice Updated, got %s", repo.projected.DisplayName)
	}
	if repo.projected.Username != "alice" {
		t.Fatalf("expected username alice, got %s", repo.projected.Username)
	}
	if repo.projected.AvatarObjectKey != "avatars/alice.png" {
		t.Fatalf("expected avatar avatars/alice.png, got %s", repo.projected.AvatarObjectKey)
	}
}

func TestHandleAccountEventCreatedFallsBackToEmailWhenDisplayNameMissing(t *testing.T) {
	repo := &roomAccountProjectionRepoStub{}
	handler := &messageHandler{accountRepo: repo}

	raw := []byte(`{
		"id": 22,
		"aggregate_id": "acc-legacy",
		"aggregate_type": "AccountAggregate",
		"version": 1,
		"event_name": "EventAccountCreated",
		"event_data": "{\"AccountID\":\"acc-legacy\",\"Email\":\"legacy@example.com\",\"CreatedAt\":\"2026-04-06T02:16:22.067488606+07:00\"}",
		"created_at": "2026-04-05T19:16:22.000000Z"
	}`)

	if err := handler.handleAccountEvent(context.Background(), raw); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.projected == nil {
		t.Fatalf("expected projected account to be saved")
	}
	if repo.projected.DisplayName != "legacy@example.com" {
		t.Fatalf("expected display name fallback to email, got %q", repo.projected.DisplayName)
	}
	if repo.projected.UpdatedAt.IsZero() {
		t.Fatalf("expected updated_at to be populated")
	}
}

func TestResolveAccountCreatedDisplayNameFallsBackInPriorityOrder(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name    string
		payload *sharedevents.AccountCreatedEvent
		want    string
	}{
		{
			name: "uses display name when present",
			payload: &sharedevents.AccountCreatedEvent{
				AccountID:   "acc-1",
				Email:       "user@example.com",
				DisplayName: "Alice",
				CreatedAt:   now,
			},
			want: "Alice",
		},
		{
			name: "falls back to email for legacy payload",
			payload: &sharedevents.AccountCreatedEvent{
				AccountID: "acc-2",
				Email:     "legacy@example.com",
				CreatedAt: now,
			},
			want: "legacy@example.com",
		},
		{
			name: "falls back to account id when email missing",
			payload: &sharedevents.AccountCreatedEvent{
				AccountID: "acc-3",
				CreatedAt: now,
			},
			want: "acc-3",
		},
	}

	for _, tt := range tests {
		if got := resolveAccountCreatedDisplayName(tt.payload); got != tt.want {
			t.Fatalf("%s: expected %q, got %q", tt.name, tt.want, got)
		}
	}
}
