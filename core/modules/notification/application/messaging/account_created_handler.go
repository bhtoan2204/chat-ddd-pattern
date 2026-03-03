package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	stackerr "go-socket/core/shared/pkg/stackErr"
)

func (h *messageHandler) handleAccountCreated(ctx context.Context, value []byte) error {
	var event accountOutboxMessage
	if err := json.Unmarshal(value, &event); err != nil {
		return stackerr.Error(fmt.Errorf("unmarshal account outbox event failed: %w", err))
	}

	if event.EventName != accountCreatedEventName {
		return nil
	}

	payload, err := parseAccountCreatedPayload(event.EventData)
	if err != nil {
		return stackerr.Error(fmt.Errorf("parse account created payload failed: %w", err))
	}

	subject := "Welcome to Go Socket"
	body := fmt.Sprintf("Welcome %s!", payload.AccountID)
	if err := h.emailSender.Send(ctx, payload.Email, subject, body); err != nil {
		return stackerr.Error(fmt.Errorf("send welcome email failed: %w", err))
	}

	return nil
}
