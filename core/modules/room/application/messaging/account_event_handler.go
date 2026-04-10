package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"go-socket/core/modules/account/domain/aggregate"
	"go-socket/core/modules/room/domain/entity"
	"go-socket/core/shared/pkg/stackErr"
)

func (h *messageHandler) handleAccountCreatedEvent(ctx context.Context, raw json.RawMessage) error {
	payloadAny, err := decodeEventPayload(ctx, "EventAccountCreated", raw)
	if err != nil {
		return stackErr.Error(fmt.Errorf("decode event payload failed: %v", err))
	}

	payload, ok := payloadAny.(*aggregate.EventAccountCreated)
	if !ok {
		return stackErr.Error(fmt.Errorf("invalid payload type for event %s", "EventAccountCreated"))
	}

	if err := h.accountRepo.ProjectAccount(ctx, &entity.AccountEntity{
		AccountID:   payload.AccountID,
		DisplayName: payload.DisplayName,
		CreatedAt:   payload.CreatedAt,
	}); err != nil {
		return stackErr.Error(err)
	}

	return nil
}

func (h *messageHandler) handleAccountUpdatedEvent(ctx context.Context, raw json.RawMessage) error {
	payloadAny, err := decodeEventPayload(ctx, "EventAccountProfileUpdated", raw)
	if err != nil {
		return stackErr.Error(fmt.Errorf("decode event payload failed: %v", err))
	}

	payload, ok := payloadAny.(*aggregate.EventAccountProfileUpdated)
	if !ok {
		return stackErr.Error(fmt.Errorf("invalid payload type for event %s", "EventAccountProfileUpdated"))
	}

	if err := h.accountRepo.ProjectAccount(ctx, &entity.AccountEntity{
		AccountID:   payload.AccountID,
		DisplayName: payload.DisplayName,
		UpdatedAt:   payload.UpdatedAt,
		AvatarObjectKey: func(data *string) string {
			if data != nil {
				return *payload.AvatarObjectKey
			}
			return ""
		}(payload.AvatarObjectKey),
		Username: func(data *string) string {
			if data != nil {
				return *payload.AvatarObjectKey
			}
			return ""
		}(payload.Username),
	}); err != nil {
		return stackErr.Error(err)
	}

	return nil
}
