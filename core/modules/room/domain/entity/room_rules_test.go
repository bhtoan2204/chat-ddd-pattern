package entity

import (
	"errors"
	"testing"
	"time"

	roomtypes "go-socket/core/modules/room/types"
)

func TestNewDirectConversationRoomRejectsSelfConversation(t *testing.T) {
	_, err := NewDirectConversationRoom("room-1", "user-1", "user-1", time.Now().UTC())
	if !errors.Is(err, ErrRoomDirectSelfNotAllowed) {
		t.Fatalf("expected self-conversation error, got %v", err)
	}
}

func TestBuildGroupMemberRolesIncludesOwnerOnce(t *testing.T) {
	memberSet, err := BuildGroupMemberRoles("owner", []string{"member-1", "owner", "member-1", "member-2"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(memberSet) != 3 {
		t.Fatalf("unexpected member count: %d", len(memberSet))
	}
	if memberSet["owner"] != roomtypes.RoomRoleOwner {
		t.Fatalf("expected owner role, got %s", memberSet["owner"])
	}
}

func TestRoomMemberCanManageAndRemoveFromGroup(t *testing.T) {
	room, err := NewRoom("room-1", "Group", "", "owner", roomtypes.RoomTypeGroup, "", time.Now().UTC())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	owner, err := NewRoomMember("member-1", room.ID, "owner", roomtypes.RoomRoleOwner, time.Now().UTC())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	member, err := NewRoomMember("member-2", room.ID, "member-1", roomtypes.RoomRoleMember, time.Now().UTC())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := owner.CanManageGroup(room); err != nil {
		t.Fatalf("expected owner to manage group, got %v", err)
	}
	if err := member.CanManageGroup(room); !errors.Is(err, ErrRoomInsufficientPermission) {
		t.Fatalf("expected insufficient permission error, got %v", err)
	}
	if err := owner.CanRemoveFrom(room, "member-1"); err != nil {
		t.Fatalf("expected owner to remove another member, got %v", err)
	}
	if err := owner.CanRemoveFrom(room, "owner"); !errors.Is(err, ErrRoomOwnerCannotLeave) {
		t.Fatalf("expected owner cannot leave error, got %v", err)
	}
}
