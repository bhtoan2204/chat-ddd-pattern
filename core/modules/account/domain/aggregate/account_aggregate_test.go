package aggregate

import (
	"testing"
	"time"
)

func TestAccountAggregateRecordAccountCreated(t *testing.T) {
	agg, err := NewAccountAggregate("account-1")
	if err != nil {
		t.Fatalf("NewAccountAggregate() error = %v", err)
	}

	createdAt := time.Now().UTC()
	if err := agg.RecordAccountCreated("user@example.com", "User", "active", createdAt); err != nil {
		t.Fatalf("RecordAccountCreated() error = %v", err)
	}

	if agg.AccountID != "account-1" {
		t.Fatalf("AccountID = %q, want %q", agg.AccountID, "account-1")
	}
	if agg.Email != "user@example.com" {
		t.Fatalf("Email = %q, want %q", agg.Email, "user@example.com")
	}
	if agg.DisplayName != "User" {
		t.Fatalf("DisplayName = %q, want %q", agg.DisplayName, "User")
	}
	if agg.Status != "active" {
		t.Fatalf("Status = %q, want %q", agg.Status, "active")
	}
	if !agg.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt = %v, want %v", agg.CreatedAt, createdAt)
	}
}
