package messaging

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	out "wechat-clone/core/modules/room/application/dto/out"
	"wechat-clone/core/modules/room/application/service"
	"wechat-clone/core/modules/room/domain/aggregate"
	"wechat-clone/core/modules/room/domain/entity"
	"wechat-clone/core/modules/room/domain/repos"
	roomtypes "wechat-clone/core/modules/room/types"
	sharedevents "wechat-clone/core/shared/contracts/events"

	"go.uber.org/mock/gomock"
)

func TestDecodeLedgerAccountTransferPayloadDerivesParticipantsFromEntrySign(t *testing.T) {
	raw, err := json.Marshal(sharedevents.LedgerTransaction{
		TransactionID: "txn-1",
		Currency:      "vnd",
		CreatedAt:     time.Date(2026, 4, 19, 8, 30, 0, 0, time.UTC),
		Entries: []*sharedevents.LedgerEntry{
			{AccountID: "receiver-1", Currency: "VND", Amount: 2500},
			{AccountID: "sender-1", Currency: "VND", Amount: -2500},
		},
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	payload, err := decodeLedgerAccountTransferPayload(context.Background(), raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if payload.SenderAccountID != "sender-1" {
		t.Fatalf("expected sender sender-1, got %s", payload.SenderAccountID)
	}
	if payload.ReceiverAccountID != "receiver-1" {
		t.Fatalf("expected receiver receiver-1, got %s", payload.ReceiverAccountID)
	}
	if payload.AmountMinor != 2500 {
		t.Fatalf("expected amount_minor 2500, got %d", payload.AmountMinor)
	}
	if payload.Currency != "VND" {
		t.Fatalf("expected currency VND, got %s", payload.Currency)
	}
}

func TestDecodeLedgerAccountTransferPayloadRejectsUnbalancedEntries(t *testing.T) {
	raw, err := json.Marshal(sharedevents.LedgerTransaction{
		TransactionID: "txn-2",
		Currency:      "USD",
		Entries: []*sharedevents.LedgerEntry{
			{AccountID: "sender-1", Currency: "USD", Amount: -1000},
			{AccountID: "receiver-1", Currency: "USD", Amount: 999},
		},
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	if _, err := decodeLedgerAccountTransferPayload(context.Background(), raw); err == nil {
		t.Fatal("expected unbalanced ledger payload to fail")
	}
}

func TestHandleLedgerAccountTransferredToAccountSkipsDuplicateTransferMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseRepo := repos.NewMockRepos(ctrl)
	messageRepo := repos.NewMockMessageRepository(ctrl)

	handler := &messageHandler{
		baseRepo: baseRepo,
	}

	raw, err := json.Marshal(sharedevents.LedgerTransaction{
		TransactionID: "txn-dup",
		Currency:      "VND",
		Entries: []*sharedevents.LedgerEntry{
			{AccountID: "sender-1", Currency: "VND", Amount: -2500},
			{AccountID: "receiver-1", Currency: "VND", Amount: 2500},
		},
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	baseRepo.EXPECT().MessageRepository().Return(messageRepo).Times(1)
	messageRepo.EXPECT().
		GetMessageByID(gomock.Any(), transferMessageID("txn-dup")).
		Return(&entity.MessageEntity{ID: transferMessageID("txn-dup")}, nil).
		Times(1)

	if err := handler.handleLedgerAccountTransferredToAccount(context.Background(), raw); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestHandleLedgerAccountTransferredToAccountCreatesDeterministicTransferMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseRepo := repos.NewMockRepos(ctrl)
	txRepos := repos.NewMockRepos(ctrl)
	roomAggRepo := repos.NewMockRoomAggregateRepository(ctrl)
	txRoomAggRepo := repos.NewMockRoomAggregateRepository(ctrl)
	messageRepo := repos.NewMockMessageRepository(ctrl)
	realtime := service.NewMockService(ctrl)

	roomAgg := mustBuildDirectRoomAggregate(t, "room-1", "sender-1", "receiver-1")

	handler := &messageHandler{
		baseRepo: baseRepo,
		svc:      realtime,
	}

	createdAt := time.Date(2026, 4, 19, 8, 30, 0, 0, time.UTC)
	raw, err := json.Marshal(sharedevents.LedgerTransaction{
		TransactionID: "txn-123",
		Currency:      "VND",
		CreatedAt:     createdAt,
		Entries: []*sharedevents.LedgerEntry{
			{AccountID: "receiver-1", Currency: "VND", Amount: 2500},
			{AccountID: "sender-1", Currency: "VND", Amount: -2500},
		},
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	baseRepo.EXPECT().MessageRepository().Return(messageRepo).Times(1)
	messageRepo.EXPECT().
		GetMessageByID(gomock.Any(), transferMessageID("txn-123")).
		Return(nil, nil).
		Times(1)
	baseRepo.EXPECT().RoomAggregateRepository().Return(roomAggRepo).Times(1)
	roomAggRepo.EXPECT().
		LoadByDirectKey(gomock.Any(), entity.CanonicalDirectKey("sender-1", "receiver-1")).
		Return(roomAgg, nil).
		Times(1)
	baseRepo.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(repos.Repos) error) error {
			return fn(txRepos)
		}).
		Times(1)
	txRepos.EXPECT().RoomAggregateRepository().Return(txRoomAggRepo).Times(1)
	txRoomAggRepo.EXPECT().
		Save(gomock.Any(), roomAgg).
		Return(nil).
		Times(1)
	realtime.EXPECT().
		EmitMessage(gomock.Any(), gomock.AssignableToTypeOf(roomtypes.MessagePayload{})).
		DoAndReturn(func(_ context.Context, payload roomtypes.MessagePayload) error {
			out, ok := payload.Payload.(*out.ChatMessageResponse)
			if !ok {
				t.Fatalf("expected chat message response payload, got %T", payload.Payload)
			}
			if out.ID != transferMessageID("txn-123") {
				t.Fatalf("expected deterministic message id, got %s", out.ID)
			}
			if out.Message != "VND 2500" {
				t.Fatalf("expected exact transfer message body, got %q", out.Message)
			}
			if out.SenderID != "sender-1" {
				t.Fatalf("expected sender sender-1, got %s", out.SenderID)
			}
			if out.CreatedAt != createdAt.Format(time.RFC3339) {
				t.Fatalf("expected created_at %s, got %s", createdAt.Format(time.RFC3339), out.CreatedAt)
			}
			return nil
		}).
		Times(1)

	if err := handler.handleLedgerAccountTransferredToAccount(context.Background(), raw); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func mustBuildDirectRoomAggregate(t *testing.T, roomID, senderID, receiverID string) *aggregate.RoomStateAggregate {
	t.Helper()

	now := time.Date(2026, 4, 19, 8, 0, 0, 0, time.UTC)
	room, err := entity.NewDirectConversationRoom(roomID, senderID, receiverID, now)
	if err != nil {
		t.Fatalf("build room: %v", err)
	}

	ownerMember, err := entity.NewRoomMember("member-1", room.ID, senderID, roomtypes.RoomRoleOwner, now)
	if err != nil {
		t.Fatalf("build owner member: %v", err)
	}
	receiverMember, err := entity.NewRoomMember("member-2", room.ID, receiverID, roomtypes.RoomRoleMember, now)
	if err != nil {
		t.Fatalf("build receiver member: %v", err)
	}

	agg, err := aggregate.RestoreRoomStateAggregate(room, []*entity.RoomMemberEntity{ownerMember, receiverMember}, 1)
	if err != nil {
		t.Fatalf("build room aggregate: %v", err)
	}

	return agg
}
