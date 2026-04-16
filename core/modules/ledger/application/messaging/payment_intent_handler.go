package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ledgerservice "go-socket/core/modules/ledger/application/service"
	ledgerentity "go-socket/core/modules/ledger/domain/entity"
	sharedevents "go-socket/core/shared/contracts/events"
	sharedlock "go-socket/core/shared/infra/lock"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

func (h *messageHandler) handlePaymentOutboxEvent(ctx context.Context, value []byte) error {
	log := logging.FromContext(ctx).Named("LedgerPaymentEvent")

	var event outboxMessage
	if err := json.Unmarshal(value, &event); err != nil {
		return stackErr.Error(fmt.Errorf("unmarshal payment outbox event failed: %v", err))
	}

	log.Infow("handle payment outbox event",
		zap.String("event_name", event.EventName),
		zap.String("aggregate_id", event.AggregateID),
	)

	switch event.EventName {
	case sharedevents.EventPaymentSucceeded:
		payload, err := unmarshalPaymentSucceededPayload(event.EventData)
		if err != nil {
			return stackErr.Error(err)
		}
		payload.PaymentID = resolvePaymentSucceededID(event.AggregateID, payload)
		command := ledgerservice.RecordPaymentSucceededCommand{
			PaymentID:          payload.PaymentID,
			TransactionID:      payload.TransactionID,
			ClearingAccountKey: payload.ClearingAccountKey,
			CreditAccountID:    payload.CreditAccountID,
			Currency:           payload.Currency,
			Amount:             payload.Amount,
		}

		lockKeys, err := paymentSucceededAccountLockKeys(command)
		if err != nil {
			return stackErr.Error(h.ledgerService.RecordPaymentSucceeded(ctx, command))
		}

		opts := sharedlock.DefaultMultiLockOptions()
		opts.KeyPrefix = ledgerservice.LedgerAccountLockKeyPrefix

		_, err = sharedlock.WithLocks(ctx, h.locker, lockKeys, opts, func() (struct{}, error) {
			return struct{}{}, h.ledgerService.RecordPaymentSucceeded(ctx, command)
		})
		if err != nil {
			return stackErr.Error(err)
		}

		return nil
	default:
		return nil
	}
}

func unmarshalPaymentSucceededPayload(data json.RawMessage) (sharedevents.PaymentSucceededEvent, error) {
	var payload sharedevents.PaymentSucceededEvent
	if err := json.Unmarshal(data, &payload); err == nil {
		return payload, nil
	} else {
		var raw string
		if err2 := json.Unmarshal(data, &raw); err2 != nil {
			return sharedevents.PaymentSucceededEvent{}, stackErr.Error(fmt.Errorf("unmarshal payment succeeded payload failed: %v", err))
		}
		if err2 := json.Unmarshal([]byte(raw), &payload); err2 != nil {
			return sharedevents.PaymentSucceededEvent{}, stackErr.Error(fmt.Errorf("unmarshal inner payload failed: %v", err2))
		}
	}

	return payload, nil
}

func resolvePaymentSucceededID(aggregateID string, payload sharedevents.PaymentSucceededEvent) string {
	paymentID := strings.TrimSpace(payload.PaymentID)
	if paymentID != "" {
		return paymentID
	}

	paymentID = strings.TrimSpace(aggregateID)
	if paymentID != "" {
		return paymentID
	}

	return strings.TrimSpace(payload.TransactionID)
}

func paymentSucceededAccountLockKeys(command ledgerservice.RecordPaymentSucceededCommand) ([]string, error) {
	booking, err := ledgerentity.NewPaymentSucceededBooking(ledgerentity.PaymentSucceededBookingInput{
		PaymentID:          command.PaymentID,
		TransactionID:      command.TransactionID,
		ClearingAccountKey: command.ClearingAccountKey,
		CreditAccountID:    command.CreditAccountID,
		Currency:           command.Currency,
		Amount:             command.Amount,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return []string{booking.DebitAccountID, booking.CreditAccountID}, nil
}
