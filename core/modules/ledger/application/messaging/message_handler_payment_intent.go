package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ledgerservice "wechat-clone/core/modules/ledger/application/service"
	ledgeraggregate "wechat-clone/core/modules/ledger/domain/aggregate"
	ledgerentity "wechat-clone/core/modules/ledger/domain/entity"
	valueobject "wechat-clone/core/modules/ledger/domain/value_object"
	paymententity "wechat-clone/core/modules/payment/domain/entity"
	"wechat-clone/core/shared/contracts"
	sharedevents "wechat-clone/core/shared/contracts/events"
	sharedlock "wechat-clone/core/shared/infra/lock"
	eventpkg "wechat-clone/core/shared/pkg/event"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

func (h *messageHandler) handlePaymentOutboxEvent(ctx context.Context, value []byte) error {
	log := logging.FromContext(ctx).Named("LedgerPaymentEvent")

	var event contracts.OutboxMessage
	if err := json.Unmarshal(value, &event); err != nil {
		return stackErr.Error(fmt.Errorf("unmarshal payment outbox event failed: %w", err))
	}

	log.Infow("handle payment outbox event",
		zap.String("event_name", event.EventName),
		zap.String("aggregate_id", event.AggregateID),
	)

	var (
		events []eventpkg.Event
		err    error
	)

	switch event.EventName {
	case sharedevents.EventPaymentCreated:
		payload, decodeErr := unmarshalPaymentCreatedPayload(event.EventData)
		if decodeErr != nil {
			return stackErr.Error(decodeErr)
		}
		payload.PaymentID = resolvePaymentCreatedID(event.AggregateID, payload)
		events, err = h.paymentCreatedLedgerEvents(payload)
	case sharedevents.EventPaymentSucceeded:
		payload, decodeErr := unmarshalPaymentSucceededPayload(event.EventData)
		if decodeErr != nil {
			return stackErr.Error(decodeErr)
		}
		payload.PaymentID = resolvePaymentSucceededID(event.AggregateID, payload)
		events, err = h.paymentSucceededLedgerEvents(payload)
	case sharedevents.EventPaymentFailed:
		payload, decodeErr := unmarshalPaymentFailedPayload(event.EventData)
		if decodeErr != nil {
			return stackErr.Error(decodeErr)
		}
		payload.PaymentID = resolvePaymentFailedID(event.AggregateID, payload)
		events, err = h.paymentFailedLedgerEvents(payload)
	case sharedevents.EventPaymentRefunded:
		payload, decodeErr := unmarshalPaymentRefundedPayload(event.EventData)
		if decodeErr != nil {
			return stackErr.Error(decodeErr)
		}
		payload.PaymentID = resolvePaymentRefundedID(event.AggregateID, payload)
		events, err = h.paymentReversedLedgerEvents(
			payload.PaymentID,
			payload.TransactionID,
			payload.ClearingAccountKey,
			payload.CreditAccountID,
			payload.Currency,
			payload.Amount,
			payload.FeeAmount,
			sharedevents.EventPaymentRefunded,
			payload.RefundedAt,
		)
	case sharedevents.EventPaymentChargeback:
		payload, decodeErr := unmarshalPaymentChargebackPayload(event.EventData)
		if decodeErr != nil {
			return stackErr.Error(decodeErr)
		}
		payload.PaymentID = resolvePaymentChargebackID(event.AggregateID, payload)
		events, err = h.paymentReversedLedgerEvents(
			payload.PaymentID,
			payload.TransactionID,
			payload.ClearingAccountKey,
			payload.CreditAccountID,
			payload.Currency,
			payload.Amount,
			payload.FeeAmount,
			sharedevents.EventPaymentChargeback,
			payload.ChargedBackAt,
		)
	default:
		return nil
	}
	if err != nil {
		return stackErr.Error(err)
	}
	if len(events) == 0 {
		return nil
	}

	return stackErr.Error(h.withLedgerAccountLocks(ctx, ledgerEventLockKeys(events), func() error {
		return h.ledgerService.RecordLedgerEvents(ctx, ledgerservice.RecordLedgerEventsCommand{Events: events})
	}))
}

func (h *messageHandler) withLedgerAccountLocks(ctx context.Context, lockKeys []string, fn func() error) error {
	opts := sharedlock.DefaultMultiLockOptions()
	opts.KeyPrefix = ledgerservice.LedgerAccountLockKeyPrefix

	_, err := sharedlock.WithLocks(ctx, h.locker, lockKeys, opts, func() (struct{}, error) {
		return struct{}{}, fn()
	})
	if err != nil {
		return stackErr.Error(err)
	}
	return nil
}

func unmarshalPaymentCreatedPayload(data json.RawMessage) (sharedevents.PaymentCreatedEvent, error) {
	var payload sharedevents.PaymentCreatedEvent
	if err := contracts.UnmarshalEventData(data, &payload); err != nil {
		return sharedevents.PaymentCreatedEvent{}, stackErr.Error(fmt.Errorf("unmarshal payment created payload failed: %w", err))
	}
	return payload, nil
}

func unmarshalPaymentSucceededPayload(data json.RawMessage) (sharedevents.PaymentSucceededEvent, error) {
	var payload sharedevents.PaymentSucceededEvent
	if err := contracts.UnmarshalEventData(data, &payload); err != nil {
		return sharedevents.PaymentSucceededEvent{}, stackErr.Error(fmt.Errorf("unmarshal payment succeeded payload failed: %w", err))
	}
	return payload, nil
}

func unmarshalPaymentFailedPayload(data json.RawMessage) (sharedevents.PaymentFailedEvent, error) {
	var payload sharedevents.PaymentFailedEvent
	if err := contracts.UnmarshalEventData(data, &payload); err != nil {
		return sharedevents.PaymentFailedEvent{}, stackErr.Error(fmt.Errorf("unmarshal payment failed payload failed: %w", err))
	}
	return payload, nil
}

func unmarshalPaymentRefundedPayload(data json.RawMessage) (sharedevents.PaymentRefundedEvent, error) {
	var payload sharedevents.PaymentRefundedEvent
	if err := contracts.UnmarshalEventData(data, &payload); err != nil {
		return sharedevents.PaymentRefundedEvent{}, stackErr.Error(fmt.Errorf("unmarshal payment refunded payload failed: %w", err))
	}
	return payload, nil
}

func unmarshalPaymentChargebackPayload(data json.RawMessage) (sharedevents.PaymentChargebackEvent, error) {
	var payload sharedevents.PaymentChargebackEvent
	if err := contracts.UnmarshalEventData(data, &payload); err != nil {
		return sharedevents.PaymentChargebackEvent{}, stackErr.Error(fmt.Errorf("unmarshal payment chargeback payload failed: %w", err))
	}
	return payload, nil
}

func resolvePaymentCreatedID(aggregateID string, payload sharedevents.PaymentCreatedEvent) string {
	return firstNonEmpty(strings.TrimSpace(payload.PaymentID), strings.TrimSpace(aggregateID), strings.TrimSpace(payload.TransactionID))
}

func resolvePaymentSucceededID(aggregateID string, payload sharedevents.PaymentSucceededEvent) string {
	return firstNonEmpty(strings.TrimSpace(payload.PaymentID), strings.TrimSpace(aggregateID), strings.TrimSpace(payload.TransactionID))
}

func resolvePaymentFailedID(aggregateID string, payload sharedevents.PaymentFailedEvent) string {
	return firstNonEmpty(strings.TrimSpace(payload.PaymentID), strings.TrimSpace(aggregateID), strings.TrimSpace(payload.TransactionID))
}

func resolvePaymentRefundedID(aggregateID string, payload sharedevents.PaymentRefundedEvent) string {
	return firstNonEmpty(strings.TrimSpace(payload.PaymentID), strings.TrimSpace(aggregateID), strings.TrimSpace(payload.TransactionID))
}

func resolvePaymentChargebackID(aggregateID string, payload sharedevents.PaymentChargebackEvent) string {
	return firstNonEmpty(strings.TrimSpace(payload.PaymentID), strings.TrimSpace(aggregateID), strings.TrimSpace(payload.TransactionID))
}

func (h *messageHandler) paymentCreatedLedgerEvents(payload sharedevents.PaymentCreatedEvent) ([]eventpkg.Event, error) {
	if paymententity.NormalizePaymentWorkflow(payload.Workflow) != paymententity.PaymentWorkflowWithdrawal {
		return nil, nil
	}

	events := make([]eventpkg.Event, 0, 4)
	principalEvents, err := transferLedgerEvents(
		fmt.Sprintf("payment:%s:withdrawal:principal", strings.TrimSpace(payload.PaymentID)),
		strings.TrimSpace(payload.DebitAccountID),
		ledgerClearingAccountID(payload.ClearingAccountKey),
		payload.Currency,
		payload.Amount,
		payload.CreatedAt,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	events = append(events, principalEvents...)

	if payload.FeeAmount > 0 && strings.TrimSpace(h.feeAccountID) != "" {
		feeEvents, err := transferLedgerEvents(
			fmt.Sprintf("payment:%s:withdrawal:fee", strings.TrimSpace(payload.PaymentID)),
			strings.TrimSpace(payload.DebitAccountID),
			strings.TrimSpace(h.feeAccountID),
			payload.Currency,
			payload.FeeAmount,
			payload.CreatedAt,
		)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		events = append(events, feeEvents...)
	}

	return events, nil
}

func (h *messageHandler) paymentSucceededLedgerEvents(payload sharedevents.PaymentSucceededEvent) ([]eventpkg.Event, error) {
	if paymententity.NormalizePaymentWorkflow(payload.Workflow) == paymententity.PaymentWorkflowWithdrawal {
		return nil, nil
	}

	booking, err := ledgerentity.NewPaymentSucceededBooking(ledgerentity.PaymentSucceededBookingInput{
		PaymentID:          payload.PaymentID,
		TransactionID:      payload.TransactionID,
		ClearingAccountKey: payload.ClearingAccountKey,
		CreditAccountID:    payload.CreditAccountID,
		Currency:           payload.Currency,
		Amount:             payload.Amount,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	events, err := paymentLedgerEventsFromBooking(
		booking.LedgerTransactionID(),
		booking.PaymentID,
		booking.Currency,
		booking.Amount,
		booking.DebitAccountID,
		booking.CreditAccountID,
		ledgeraggregate.EventNameLedgerAccountWithdrawFromIntent,
		ledgeraggregate.EventNameLedgerAccountDepositFromIntent,
		payload.SucceededAt,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	if payload.FeeAmount > 0 && strings.TrimSpace(h.feeAccountID) != "" {
		feeEvents, err := paymentLedgerEventsFromBooking(
			fmt.Sprintf("payment:%s:succeeded:fee", strings.TrimSpace(booking.PaymentID)),
			booking.PaymentID,
			booking.Currency,
			payload.FeeAmount,
			ledgerClearingAccountID(payload.ClearingAccountKey),
			strings.TrimSpace(h.feeAccountID),
			ledgeraggregate.EventNameLedgerAccountWithdrawFromIntent,
			ledgeraggregate.EventNameLedgerAccountDepositFromIntent,
			payload.SucceededAt,
		)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		events = append(events, feeEvents...)
	}

	return events, nil
}

func (h *messageHandler) paymentFailedLedgerEvents(payload sharedevents.PaymentFailedEvent) ([]eventpkg.Event, error) {
	if paymententity.NormalizePaymentWorkflow(payload.Workflow) != paymententity.PaymentWorkflowWithdrawal {
		return nil, nil
	}

	events := make([]eventpkg.Event, 0, 4)
	principalEvents, err := transferLedgerEvents(
		fmt.Sprintf("payment:%s:withdrawal:principal:failed", strings.TrimSpace(payload.PaymentID)),
		ledgerClearingAccountID(payload.ClearingAccountKey),
		strings.TrimSpace(payload.DebitAccountID),
		payload.Currency,
		payload.Amount,
		payload.OccurredAt,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	events = append(events, principalEvents...)

	if payload.FeeAmount > 0 && strings.TrimSpace(h.feeAccountID) != "" {
		feeEvents, err := transferLedgerEvents(
			fmt.Sprintf("payment:%s:withdrawal:fee:failed", strings.TrimSpace(payload.PaymentID)),
			strings.TrimSpace(h.feeAccountID),
			strings.TrimSpace(payload.DebitAccountID),
			payload.Currency,
			payload.FeeAmount,
			payload.OccurredAt,
		)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		events = append(events, feeEvents...)
	}

	return events, nil
}

func (h *messageHandler) paymentReversedLedgerEvents(paymentID, transactionID, clearingAccountKey, creditAccountID, currency string, amount int64, feeAmount int64, reversalType string, bookedAt time.Time) ([]eventpkg.Event, error) {
	booking, err := ledgerentity.NewPaymentReversalBooking(ledgerentity.PaymentReversalBookingInput{
		PaymentID:          paymentID,
		TransactionID:      transactionID,
		ClearingAccountKey: clearingAccountKey,
		CreditAccountID:    creditAccountID,
		Currency:           currency,
		Amount:             amount,
		ReversalType:       reversalType,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	events, err := paymentLedgerEventsFromBooking(
		booking.LedgerTransactionID(),
		booking.PaymentID,
		booking.Currency,
		booking.Amount,
		booking.DebitAccountID,
		booking.CreditAccountID,
		debitLedgerEventNameForReversal(booking.ReversalType),
		creditLedgerEventNameForReversal(booking.ReversalType),
		bookedAt,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	if feeAmount > 0 && strings.TrimSpace(h.feeAccountID) != "" {
		feeEvents, err := paymentLedgerEventsFromBooking(
			fmt.Sprintf("payment:%s:%s:fee", strings.TrimSpace(booking.PaymentID), reversalSuffix(reversalType)),
			booking.PaymentID,
			booking.Currency,
			feeAmount,
			strings.TrimSpace(h.feeAccountID),
			ledgerClearingAccountID(clearingAccountKey),
			debitLedgerEventNameForReversal(reversalType),
			creditLedgerEventNameForReversal(reversalType),
			bookedAt,
		)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		events = append(events, feeEvents...)
	}

	return events, nil
}

func paymentLedgerEventsFromBooking(transactionID, paymentID, currency string, amount int64, debitAccountID, creditAccountID, debitEventName, creditEventName string, bookedAt time.Time) ([]eventpkg.Event, error) {
	if bookedAt.IsZero() {
		bookedAt = time.Now().UTC()
	}

	debitPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		valueobject.LedgerAccountPostingInput{
			AccountID:             debitAccountID,
			TransactionID:         transactionID,
			ReferenceType:         debitEventName,
			ReferenceID:           paymentID,
			CounterpartyAccountID: creditAccountID,
			Currency:              currency,
			AmountDelta:           -amount,
			BookedAt:              bookedAt,
		},
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	creditPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		valueobject.LedgerAccountPostingInput{
			AccountID:             creditAccountID,
			TransactionID:         transactionID,
			ReferenceType:         creditEventName,
			ReferenceID:           paymentID,
			CounterpartyAccountID: debitAccountID,
			Currency:              currency,
			AmountDelta:           amount,
			BookedAt:              bookedAt,
		},
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	debitEvent, ok, err := ledgeraggregate.NewLedgerAccountEventFromPosting(debitAccountID, debitPosting)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if !ok {
		return nil, stackErr.Error(fmt.Errorf("unsupported debit ledger event reference_type=%s", debitEventName))
	}
	creditEvent, ok, err := ledgeraggregate.NewLedgerAccountEventFromPosting(creditAccountID, creditPosting)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if !ok {
		return nil, stackErr.Error(fmt.Errorf("unsupported credit ledger event reference_type=%s", creditEventName))
	}

	return []eventpkg.Event{debitEvent, creditEvent}, nil
}

func transferLedgerEvents(transactionID, fromAccountID, toAccountID, currency string, amount int64, bookedAt time.Time) ([]eventpkg.Event, error) {
	if bookedAt.IsZero() {
		bookedAt = time.Now().UTC()
	}

	if _, err := ledgeraggregate.NewLedgerAccountTransferOutPosting(
		valueobject.LedgerAccountTransferPostingInput{
			AccountID:             fromAccountID,
			TransactionID:         transactionID,
			CounterpartyAccountID: toAccountID,
			Currency:              currency,
			Amount:                amount,
			BookedAt:              bookedAt,
		},
	); err != nil {
		return nil, stackErr.Error(err)
	}
	if _, err := ledgeraggregate.NewLedgerAccountTransferInPosting(
		valueobject.LedgerAccountTransferPostingInput{
			AccountID:             toAccountID,
			TransactionID:         transactionID,
			CounterpartyAccountID: fromAccountID,
			Currency:              currency,
			Amount:                amount,
			BookedAt:              bookedAt,
		},
	); err != nil {
		return nil, stackErr.Error(err)
	}

	debitEvent, err := ledgeraggregate.NewLedgerAccountEvent(
		fromAccountID,
		eventpkg.AggregateTypeName(&ledgeraggregate.LedgerAccountAggregate{}),
		&ledgeraggregate.EventLedgerAccountTransferredToAccount{
			TransactionID: transactionID,
			ToAccountID:   toAccountID,
			Currency:      currency,
			Amount:        amount,
			BookedAt:      bookedAt,
		},
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	creditEvent, err := ledgeraggregate.NewLedgerAccountEvent(
		toAccountID,
		eventpkg.AggregateTypeName(&ledgeraggregate.LedgerAccountAggregate{}),
		&ledgeraggregate.EventLedgerAccountReceivedTransfer{
			TransactionID: transactionID,
			FromAccountID: fromAccountID,
			Currency:      currency,
			Amount:        amount,
			BookedAt:      bookedAt,
		},
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return []eventpkg.Event{debitEvent, creditEvent}, nil
}

func debitLedgerEventNameForReversal(paymentEventName string) string {
	switch strings.TrimSpace(paymentEventName) {
	case sharedevents.EventPaymentRefunded:
		return ledgeraggregate.EventNameLedgerAccountWithdrawFromRefund
	case sharedevents.EventPaymentChargeback:
		return ledgeraggregate.EventNameLedgerAccountWithdrawFromChargeback
	default:
		return ""
	}
}

func creditLedgerEventNameForReversal(paymentEventName string) string {
	switch strings.TrimSpace(paymentEventName) {
	case sharedevents.EventPaymentRefunded:
		return ledgeraggregate.EventNameLedgerAccountDepositFromRefund
	case sharedevents.EventPaymentChargeback:
		return ledgeraggregate.EventNameLedgerAccountDepositFromChargeback
	default:
		return ""
	}
}

func ledgerEventLockKeys(events []eventpkg.Event) []string {
	keys := make([]string, 0, len(events))
	seen := make(map[string]struct{}, len(events))
	for _, evt := range events {
		key := strings.TrimSpace(evt.AggregateID)
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	return keys
}

func ledgerClearingAccountID(clearingAccountKey string) string {
	return fmt.Sprintf("ledger:clearing:%s", strings.ToLower(strings.TrimSpace(clearingAccountKey)))
}

func reversalSuffix(reversalType string) string {
	switch strings.TrimSpace(reversalType) {
	case sharedevents.EventPaymentRefunded:
		return "refunded"
	case sharedevents.EventPaymentChargeback:
		return "chargeback"
	default:
		return "reversed"
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}
