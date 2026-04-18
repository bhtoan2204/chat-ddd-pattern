package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	reflect "reflect"
	"time"

	roomsupport "wechat-clone/core/modules/room/application/support"
	"wechat-clone/core/modules/room/domain/aggregate"
	"wechat-clone/core/modules/room/domain/entity"
	"wechat-clone/core/modules/room/domain/repos"
	"wechat-clone/core/modules/room/types"
	"wechat-clone/core/shared/pkg/stackErr"
)

func (h *messageHandler) handleLedgerAccountTransferredToAccount(ctx context.Context, raw json.RawMessage) error {
	transfer, err := decodeLedgerAccountTransferPayload(ctx, raw)
	if err != nil {
		return stackErr.Error(fmt.Errorf("decode ledger transfer payload failed: %w", err))
	}

	messageID := transferMessageID(transfer.TransactionID)
	existingMessage, err := h.baseRepo.MessageRepository().GetMessageByID(ctx, messageID)
	if err != nil {
		return stackErr.Error(fmt.Errorf("load transfer message failed: %w", err))
	}
	if existingMessage != nil {
		return nil
	}

	roomAgg, err := h.baseRepo.RoomAggregateRepository().LoadByDirectKey(ctx, entity.CanonicalDirectKey(transfer.SenderAccountID, transfer.ReceiverAccountID))
	if err != nil {
		return stackErr.Error(fmt.Errorf("load transfer room failed: %w", err))
	}
	if roomAgg == nil || roomAgg.Room() == nil {
		return stackErr.Error(fmt.Errorf("direct room not found for transfer participants"))
	}

	now := transfer.CreatedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}

	message, err := roomAgg.SendMessage(
		messageID,
		transfer.SenderAccountID,
		entity.MessageParams{
			Message:     formatTransferMessageBody(transfer.Currency, transfer.AmountMinor),
			MessageType: entity.MessageTypeTransfer,
		},
		aggregate.MessageSenderIdentity{},
		aggregate.MessageOutboxPayload{},
		now,
	)
	if err != nil {
		return stackErr.Error(err)
	}

	if err := h.baseRepo.WithTransaction(ctx, func(txRepos repos.Repos) error {
		return stackErr.Error(txRepos.RoomAggregateRepository().Save(ctx, roomAgg))
	}); err != nil {
		return stackErr.Error(err)
	}

	msg, err := roomsupport.BuildMessageResultFromState(ctx, h.baseRepo, transfer.SenderAccountID, message)
	if err != nil {
		return stackErr.Error(err)
	}
	out := roomsupport.ToMessageResponse(msg)
	if err := h.svc.EmitMessage(ctx, types.MessagePayload{
		RoomId:  out.RoomID,
		Type:    reflect.TypeOf(out).Elem().Name(),
		Payload: out,
	}); err != nil {
		return stackErr.Error(fmt.Errorf("failed to emit realtime message after handling ledger transfer event: %w", err))
	}

	return nil
}
