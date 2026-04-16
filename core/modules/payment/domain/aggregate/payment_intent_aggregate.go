package aggregate

import (
	"time"

	paymententity "go-socket/core/modules/payment/domain/entity"
	sharedevents "go-socket/core/shared/contracts/events"
	eventpkg "go-socket/core/shared/pkg/event"
	"go-socket/core/shared/pkg/stackErr"
)

const AggregateTypePaymentIntent = "PaymentIntentAggregate"

type PaymentIntentMutation struct {
	Duplicate bool
	Persist   bool
}

type PaymentIntentAggregate struct {
	eventpkg.Aggregate

	intent          *paymententity.PaymentIntent
	processedEvents []*paymententity.ProcessedPaymentEvent
	outboxEvents    []eventpkg.Event
	version         int
}

func NewProviderTopUpAggregate(
	transactionID,
	provider string,
	amount int64,
	currency,
	creditAccountID string,
	metadata map[string]string,
	now time.Time,
) (*PaymentIntentAggregate, error) {
	intent, err := paymententity.NewProviderTopUpIntent(
		transactionID,
		provider,
		amount,
		currency,
		creditAccountID,
		now,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	agg := &PaymentIntentAggregate{intent: intent}
	agg.recordOutboxEvent(sharedevents.EventPaymentCreated, intent.BuildCreatedEventData(metadata, now), now)
	return agg, nil
}

func RehydratePaymentIntentAggregate(intent *paymententity.PaymentIntent) (*PaymentIntentAggregate, error) {
	if intent == nil {
		return nil, stackErr.Error(paymententity.ErrPaymentTransactionIDRequired)
	}
	clone := *intent
	if err := clone.ApplyProviderResult(clone.CurrentProviderResult(paymententity.PaymentProviderResult{}), clone.UpdatedAt); err != nil {
		return nil, stackErr.Error(err)
	}
	clone.UpdatedAt = intent.UpdatedAt.UTC()
	clone.CreatedAt = intent.CreatedAt.UTC()
	return &PaymentIntentAggregate{intent: &clone}, nil
}

func (a *PaymentIntentAggregate) Snapshot() *paymententity.PaymentIntent {
	if a == nil || a.intent == nil {
		return nil
	}
	clone := *a.intent
	return &clone
}

func (a *PaymentIntentAggregate) TransactionID() string {
	if a == nil || a.intent == nil {
		return ""
	}
	return a.intent.TransactionID
}

func (a *PaymentIntentAggregate) Provider() string {
	if a == nil || a.intent == nil {
		return ""
	}
	return a.intent.Provider
}

func (a *PaymentIntentAggregate) ExternalRef() string {
	if a == nil || a.intent == nil {
		return ""
	}
	return a.intent.ExternalRef
}

func (a *PaymentIntentAggregate) Status() string {
	if a == nil || a.intent == nil {
		return ""
	}
	return a.intent.Status
}

func (a *PaymentIntentAggregate) ValidateProviderResultForStatus(status string, amount int64, currency string) error {
	if a == nil || a.intent == nil {
		return stackErr.Error(paymententity.ErrPaymentTransactionIDRequired)
	}
	if err := a.intent.ValidateProviderResultForStatus(status, amount, currency); err != nil {
		return stackErr.Error(err)
	}
	return nil
}

func (a *PaymentIntentAggregate) PendingProcessedEvents() []*paymententity.ProcessedPaymentEvent {
	if len(a.processedEvents) == 0 {
		return nil
	}
	items := make([]*paymententity.ProcessedPaymentEvent, 0, len(a.processedEvents))
	for _, item := range a.processedEvents {
		if item == nil {
			continue
		}
		clone := *item
		items = append(items, &clone)
	}
	return items
}

func (a *PaymentIntentAggregate) PendingOutboxEvents() []eventpkg.Event {
	if len(a.outboxEvents) == 0 {
		return nil
	}
	items := make([]eventpkg.Event, len(a.outboxEvents))
	copy(items, a.outboxEvents)
	return items
}

func (a *PaymentIntentAggregate) MarkPersisted() {
	if a == nil {
		return
	}
	a.processedEvents = nil
	a.outboxEvents = nil
}

func (a *PaymentIntentAggregate) ApplyProviderOutcome(
	result paymententity.PaymentProviderResult,
	checkoutURL string,
	emitCheckoutEvent bool,
	occurredAt time.Time,
) (PaymentIntentMutation, error) {
	if a == nil || a.intent == nil {
		return PaymentIntentMutation{}, stackErr.Error(paymententity.ErrPaymentTransactionIDRequired)
	}

	switch paymententity.NormalizePaymentStatus(result.Status) {
	case paymententity.PaymentStatusSuccess:
		return a.applySuccessfulOutcome(result, checkoutURL, emitCheckoutEvent, occurredAt)
	case paymententity.PaymentStatusRefunded, paymententity.PaymentStatusChargeback:
		return a.applyReversedOutcome(result, checkoutURL, emitCheckoutEvent, occurredAt)
	default:
		return a.applyNonFinalOutcome(result, checkoutURL, emitCheckoutEvent, occurredAt)
	}
}

func (a *PaymentIntentAggregate) MarkCreateFailed(occurredAt time.Time) (PaymentIntentMutation, error) {
	if a == nil || a.intent == nil {
		return PaymentIntentMutation{}, stackErr.Error(paymententity.ErrPaymentTransactionIDRequired)
	}

	transition, err := a.intent.MarkCreateFailed(occurredAt)
	if err != nil {
		return PaymentIntentMutation{}, stackErr.Error(err)
	}
	if transition.Ignored || !transition.StateChanged {
		return PaymentIntentMutation{Duplicate: true}, nil
	}

	failedResult := a.intent.CurrentProviderResult(paymententity.PaymentProviderResult{Status: paymententity.PaymentStatusFailed})
	a.recordOutboxEvent(sharedevents.EventPaymentFailed, a.intent.BuildFailedEventData(failedResult, a.intent.UpdatedAt), a.intent.UpdatedAt)
	return PaymentIntentMutation{Persist: true}, nil
}

func (a *PaymentIntentAggregate) applyNonFinalOutcome(
	result paymententity.PaymentProviderResult,
	checkoutURL string,
	emitCheckoutEvent bool,
	occurredAt time.Time,
) (PaymentIntentMutation, error) {
	transition, err := a.intent.TransitionProviderResult(result, occurredAt)
	if err != nil {
		return PaymentIntentMutation{}, stackErr.Error(err)
	}
	if transition.Ignored || (!transition.StateChanged && !transition.ExternalRefChanged) {
		return PaymentIntentMutation{Duplicate: true}, nil
	}

	a.recordCheckoutSessionEvent(checkoutURL, emitCheckoutEvent, occurredAt)
	if transition.Type == paymententity.PaymentTransitionFailed {
		a.recordOutboxEvent(
			sharedevents.EventPaymentFailed,
			a.intent.BuildFailedEventData(a.intent.CurrentProviderResult(result), occurredAt),
			occurredAt,
		)
	}

	return PaymentIntentMutation{Persist: true}, nil
}

func (a *PaymentIntentAggregate) applySuccessfulOutcome(
	result paymententity.PaymentProviderResult,
	checkoutURL string,
	emitCheckoutEvent bool,
	occurredAt time.Time,
) (PaymentIntentMutation, error) {
	transition, err := a.intent.TransitionProviderResult(result, occurredAt)
	if err != nil {
		return PaymentIntentMutation{}, stackErr.Error(err)
	}
	if transition.Ignored || transition.Type == paymententity.PaymentTransitionNone {
		if transition.ExternalRefChanged {
			a.recordCheckoutSessionEvent(checkoutURL, emitCheckoutEvent, occurredAt)
			return PaymentIntentMutation{Duplicate: true, Persist: true}, nil
		}
		return PaymentIntentMutation{Duplicate: true}, nil
	}
	if transition.Type != paymententity.PaymentTransitionSucceeded {
		return PaymentIntentMutation{Duplicate: true}, nil
	}

	processedEvent, err := a.intent.NewProcessedTransitionEvent(sharedevents.EventPaymentSucceeded, occurredAt)
	if err != nil {
		return PaymentIntentMutation{}, stackErr.Error(err)
	}

	a.recordProcessedEvent(processedEvent)
	a.recordOutboxEvent(
		sharedevents.EventPaymentSucceeded,
		a.intent.BuildSucceededEventData(a.intent.CurrentProviderResult(result), occurredAt),
		occurredAt,
	)
	a.recordCheckoutSessionEvent(checkoutURL, emitCheckoutEvent, occurredAt)
	return PaymentIntentMutation{Persist: true}, nil
}

func (a *PaymentIntentAggregate) applyReversedOutcome(
	result paymententity.PaymentProviderResult,
	checkoutURL string,
	emitCheckoutEvent bool,
	occurredAt time.Time,
) (PaymentIntentMutation, error) {
	transition, err := a.intent.TransitionProviderResult(result, occurredAt)
	if err != nil {
		return PaymentIntentMutation{}, stackErr.Error(err)
	}
	if transition.Ignored || transition.Type == paymententity.PaymentTransitionNone {
		return PaymentIntentMutation{Duplicate: true}, nil
	}

	var (
		processedEvent *paymententity.ProcessedPaymentEvent
		reversalData   interface{}
		reversalName   string
	)
	switch transition.Type {
	case paymententity.PaymentTransitionRefunded:
		processedEvent, err = a.intent.NewProcessedTransitionEvent(sharedevents.EventPaymentRefunded, occurredAt)
		if err != nil {
			return PaymentIntentMutation{}, stackErr.Error(err)
		}
		reversalName = sharedevents.EventPaymentRefunded
		reversalData = a.intent.BuildRefundedEventData(a.intent.CurrentProviderResult(result), occurredAt)
	case paymententity.PaymentTransitionChargeback:
		processedEvent, err = a.intent.NewProcessedTransitionEvent(sharedevents.EventPaymentChargeback, occurredAt)
		if err != nil {
			return PaymentIntentMutation{}, stackErr.Error(err)
		}
		reversalName = sharedevents.EventPaymentChargeback
		reversalData = a.intent.BuildChargebackEventData(a.intent.CurrentProviderResult(result), occurredAt)
	default:
		return PaymentIntentMutation{Duplicate: true}, nil
	}

	a.recordProcessedEvent(processedEvent)
	a.recordOutboxEvent(reversalName, reversalData, occurredAt)
	a.recordCheckoutSessionEvent(checkoutURL, emitCheckoutEvent, occurredAt)
	return PaymentIntentMutation{Persist: true}, nil
}

func (a *PaymentIntentAggregate) recordProcessedEvent(evt *paymententity.ProcessedPaymentEvent) {
	if a == nil || evt == nil {
		return
	}
	clone := *evt
	a.processedEvents = append(a.processedEvents, &clone)
}

// Aggregate root owns the outbox envelope; entities only supply payload facts.
func (a *PaymentIntentAggregate) recordOutboxEvent(eventName string, eventData interface{}, occurredAt time.Time) {
	if a == nil || a.intent == nil {
		return
	}
	eventTime := normalizeAggregateEventTime(occurredAt)
	a.version++
	a.outboxEvents = append(a.outboxEvents, eventpkg.Event{
		AggregateID:   a.intent.TransactionID,
		AggregateType: AggregateTypePaymentIntent,
		Version:       a.version,
		EventName:     eventName,
		EventData:     eventData,
		CreatedAt:     eventTime.Unix(),
	})
}

func (a *PaymentIntentAggregate) recordCheckoutSessionEvent(checkoutURL string, emitCheckoutEvent bool, occurredAt time.Time) {
	if a == nil || a.intent == nil {
		return
	}
	if !emitCheckoutEvent || !a.intent.ShouldEmitCheckoutSessionCreated(checkoutURL) {
		return
	}
	a.recordOutboxEvent(
		sharedevents.EventPaymentCheckoutSessionCreated,
		a.intent.BuildCheckoutSessionCreatedEventData(checkoutURL, occurredAt),
		occurredAt,
	)
}

func normalizeAggregateEventTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now().UTC()
	}
	return value.UTC()
}
