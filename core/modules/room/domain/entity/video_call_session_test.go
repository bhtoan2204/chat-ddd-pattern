package entity

import (
	"testing"
	"time"
)

func TestVideoCallSessionLeaveEndsWhenLastParticipantLeaves(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.April, 23, 8, 0, 0, 0, time.UTC)
	session, err := NewVideoCallSession("session-1", "room-1", "account-1", now)
	if err != nil {
		t.Fatalf("NewVideoCallSession() error = %v", err)
	}

	ended, err := session.Leave("account-1", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("Leave() error = %v", err)
	}
	if !ended {
		t.Fatalf("Leave() ended = false, want true")
	}
	if session.Status != VideoCallStatusEnded {
		t.Fatalf("Status = %s, want %s", session.Status, VideoCallStatusEnded)
	}
	if session.EndedByAccountID != "account-1" {
		t.Fatalf("EndedByAccountID = %s, want account-1", session.EndedByAccountID)
	}
}

func TestVideoCallSessionJoinIsIdempotent(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.April, 23, 8, 0, 0, 0, time.UTC)
	session, err := NewVideoCallSession("session-1", "room-1", "account-1", now)
	if err != nil {
		t.Fatalf("NewVideoCallSession() error = %v", err)
	}

	if err := session.Join("account-2", now.Add(time.Minute)); err != nil {
		t.Fatalf("Join() first error = %v", err)
	}
	if err := session.Join("account-2", now.Add(2*time.Minute)); err != nil {
		t.Fatalf("Join() second error = %v", err)
	}
	if len(session.ParticipantAccountIDs) != 2 {
		t.Fatalf("participant count = %d, want 2", len(session.ParticipantAccountIDs))
	}
}
