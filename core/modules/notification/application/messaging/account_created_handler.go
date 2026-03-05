package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"go-socket/core/shared/contracts/events"
	stackerr "go-socket/core/shared/pkg/stackErr"
)

func (h *messageHandler) handleAccountCreatedEvent(ctx context.Context, raw json.RawMessage) error {
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
	if err := h.emailSender.Send(ctx, payload.Email, subject, body); err != nil {
		return stackerr.Error(fmt.Errorf("send welcome email failed: %w", err))
	}

	return nil
}
