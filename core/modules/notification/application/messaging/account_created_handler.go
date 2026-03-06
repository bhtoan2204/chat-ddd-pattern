package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"go-socket/core/modules/notification/domain/entity"
	"go-socket/core/modules/notification/types"
	"go-socket/core/shared/contracts/events"
	"go-socket/core/shared/pkg/logging"
	stackerr "go-socket/core/shared/pkg/stackErr"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *messageHandler) handleAccountCreatedEvent(ctx context.Context, raw json.RawMessage) error {
	log := logging.FromContext(ctx).Named("handleAccountCreatedEvent")
	payloadAny, err := decodeEventPayload(ctx, events.AccountCreatedEventName, raw)
	if err != nil {
		return stackerr.Error(fmt.Errorf("decode event payload failed: %w", err))
	}

	payload, ok := payloadAny.(*events.AccountCreatedEvent)
	if !ok {
		return stackerr.Error(fmt.Errorf("invalid payload type for event %s", events.AccountCreatedEventName))
	}

	subject := "Welcome to Go Socket"
	body := fmt.Sprintf("Welcome %s!", payload.AccountID)

	notification := &entity.NotificationEntity{
		ID:        uuid.New().String(),
		AccountID: payload.AccountID,
		Type:      types.NotificationTypeAccountCreated,
		Subject:   subject,
		Body:      body,
		CreatedAt: payload.CreatedAt,
	}
	if err := h.notificationRepo.CreateNotification(ctx, notification); err != nil {
		log.Errorw("create notification failed", zap.Error(err))
		return stackerr.Error(fmt.Errorf("create notification failed: %w", err))
	}

	return h.emailSender.Send(ctx, payload.Email, subject, body)
}
