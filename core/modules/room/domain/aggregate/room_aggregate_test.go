package aggregate

import (
	"testing"
	"time"

	roomtypes "go-socket/core/modules/room/types"
)

func TestRoomAggregateRecordMemberAddedBindsRoomID(t *testing.T) {
	agg, err := NewRoomAggregate("room-1")
	if err != nil {
		t.Fatalf("NewRoomAggregate() error = %v", err)
	}

	if err := agg.RecordMemberAdded("member-1", roomtypes.RoomRoleMember, time.Now().UTC()); err != nil {
		t.Fatalf("RecordMemberAdded() error = %v", err)
	}

	if agg.RoomID != "room-1" {
		t.Fatalf("RoomID = %q, want %q", agg.RoomID, "room-1")
	}
	if agg.MemberCount != 1 {
		t.Fatalf("MemberCount = %d, want 1", agg.MemberCount)
	}
}

func TestRoomAggregateRecordMessageCreatedBindsRoomID(t *testing.T) {
	agg, err := NewRoomAggregate("room-1")
	if err != nil {
		t.Fatalf("NewRoomAggregate() error = %v", err)
	}

	sentAt := time.Now().UTC()
	if err := agg.RecordMessageCreated("msg-1", "sender-1", "Alice", "alice@example.com", "hello", sentAt); err != nil {
		t.Fatalf("RecordMessageCreated() error = %v", err)
	}

	if agg.RoomID != "room-1" {
		t.Fatalf("RoomID = %q, want %q", agg.RoomID, "room-1")
	}
	if agg.LastMessageID != "msg-1" {
		t.Fatalf("LastMessageID = %q, want %q", agg.LastMessageID, "msg-1")
	}
	if !agg.LastMessageAt.Equal(sentAt) {
		t.Fatalf("LastMessageAt = %v, want %v", agg.LastMessageAt, sentAt)
	}
}
