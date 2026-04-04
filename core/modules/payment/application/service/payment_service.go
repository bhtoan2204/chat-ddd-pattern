package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	paymentin "go-socket/core/modules/payment/application/dto/in"
	paymentout "go-socket/core/modules/payment/application/dto/out"
	"go-socket/core/modules/payment/domain/entity"
	paymentrepos "go-socket/core/modules/payment/domain/repos"
	"go-socket/core/modules/payment/providers"
	sharedevents "go-socket/core/shared/contracts/events"
	eventpkg "go-socket/core/shared/pkg/event"
	"go-socket/core/shared/pkg/logging"
)

const paymentAggregateType = "payment"

type PaymentService struct {
	intentStore      PaymentIntentStore
	providerRegistry *providers.ProviderRegistry
}

func NewPaymentService(intentStore PaymentIntentStore, providerRegistry *providers.ProviderRegistry) *PaymentService {
	return &PaymentService{
		intentStore:      intentStore,
		providerRegistry: providerRegistry,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *paymentin.CreatePaymentRequest) (*paymentout.CreatePaymentResponse, error) {
	req.Normalize()
	if err := wrapValidation(req.Validate()); err != nil {
		return nil, err
	}

	provider, err := s.providerRegistry.Get(req.Provider)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	intent := &entity.PaymentIntent{
		TransactionID:   req.TransactionID,
		Provider:        req.Provider,
		Amount:          req.Amount,
		Currency:        req.Currency,
		DebitAccountID:  req.DebitAccountID,
		CreditAccountID: req.CreditAccountID,
		Status:          entity.PaymentStatusCreating,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.intentStore.WithTransaction(ctx, func(store PaymentIntentStore) error {
		if err := store.CreateIntent(ctx, intent); err != nil {
			return err
		}
		return store.AppendOutboxEvent(ctx, newPaymentCreatedEvent(intent, req.Metadata))
	}); err != nil {
		if errors.Is(err, paymentrepos.ErrProviderPaymentDuplicateIntent) {
			return nil, fmt.Errorf("%v: %s", ErrDuplicatePayment, req.TransactionID)
		}
		return nil, err
	}

	response, err := provider.CreatePayment(ctx, providers.CreatePaymentRequest{
		TransactionID:   req.TransactionID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		DebitAccountID:  req.DebitAccountID,
		CreditAccountID: req.CreditAccountID,
		Metadata:        req.Metadata,
	})
	if err != nil {
		logging.FromContext(ctx).Errorw("provider create payment failed",
			"provider", provider.Name(),
			"transaction_id", req.TransactionID,
			"error", err,
		)
		_ = s.updateIntentStatus(ctx, req.TransactionID, entity.PaymentStatusFailed)
		return nil, err
	}

	targetStatus := normalizePaymentStatus(response.Status)
	if targetStatus == "" {
		targetStatus = entity.PaymentStatusPending
	}

	persistedIntent, err := s.intentStore.GetIntentByTransactionID(ctx, req.TransactionID)
	if err != nil {
		if errors.Is(err, paymentrepos.ErrProviderPaymentNotFound) {
			return nil, fmt.Errorf("%v: %s", ErrPaymentIntentNotFound, req.TransactionID)
		}
		return nil, err
	}
	if err := s.intentStore.WithTransaction(ctx, func(store PaymentIntentStore) error {
		if err := store.UpdateIntentProviderState(ctx, persistedIntent.TransactionID, response.ExternalRef, targetStatus); err != nil {
			return err
		}

		if response.CheckoutURL != "" || response.ExternalRef != "" {
			if err := store.AppendOutboxEvent(ctx, newPaymentCheckoutSessionCreatedEvent(persistedIntent, response, targetStatus)); err != nil {
				return err
			}
		}

		if targetStatus == entity.PaymentStatusSuccess {
			return s.finalizeSuccessfulPaymentTx(ctx, store, persistedIntent, &providers.PaymentResult{
				TransactionID: response.TransactionID,
				Status:        targetStatus,
				Amount:        persistedIntent.Amount,
				Currency:      persistedIntent.Currency,
				ExternalRef:   response.ExternalRef,
			})
		}
		if targetStatus == entity.PaymentStatusFailed {
			return store.AppendOutboxEvent(ctx, newPaymentFailedEvent(persistedIntent, &providers.PaymentResult{
				TransactionID: response.TransactionID,
				Status:        targetStatus,
				Amount:        persistedIntent.Amount,
				Currency:      persistedIntent.Currency,
				ExternalRef:   response.ExternalRef,
			}))
		}
		return nil
	}); err != nil {
		return nil, err
	}

	logging.FromContext(ctx).Infow("payment created",
		"provider", provider.Name(),
		"transaction_id", response.TransactionID,
		"status", targetStatus,
		"external_ref", response.ExternalRef,
	)

	return &paymentout.CreatePaymentResponse{
		Provider:      strings.ToLower(provider.Name()),
		TransactionID: response.TransactionID,
		ExternalRef:   response.ExternalRef,
		Status:        targetStatus,
		CheckoutURL:   response.CheckoutURL,
	}, nil
}

func (s *PaymentService) HandleWebhook(ctx context.Context, providerName string, payload []byte, signature string) (*paymentout.ProcessWebhookResponse, error) {
	provider, err := s.providerRegistry.Get(providerName)
	if err != nil {
		return nil, err
	}

	event, err := provider.VerifyWebhook(ctx, payload, signature)
	if err != nil {
		return nil, err
	}

	result, err := provider.ParseEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	result.Status = normalizePaymentStatus(result.Status)

	intent, err := s.findIntent(ctx, strings.ToLower(provider.Name()), result)
	if err != nil {
		return nil, err
	}
	if err := validateResultAgainstIntent(intent, result); err != nil {
		return nil, err
	}

	if result.Status != entity.PaymentStatusSuccess {
		if err := s.intentStore.WithTransaction(ctx, func(store PaymentIntentStore) error {
			if err := store.UpdateIntentProviderState(ctx, intent.TransactionID, result.ExternalRef, result.Status); err != nil {
				return err
			}
			if result.Status == entity.PaymentStatusFailed {
				return store.AppendOutboxEvent(ctx, newPaymentFailedEvent(intent, result))
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return &paymentout.ProcessWebhookResponse{
			Provider:      intent.Provider,
			TransactionID: intent.TransactionID,
			ExternalRef:   coalesce(result.ExternalRef, intent.ExternalRef),
			Status:        result.Status,
		}, nil
	}

	duplicate, err := s.finalizeSuccessfulPayment(ctx, intent, result)
	if err != nil {
		return nil, err
	}

	return &paymentout.ProcessWebhookResponse{
		Provider:      intent.Provider,
		TransactionID: intent.TransactionID,
		ExternalRef:   coalesce(result.ExternalRef, intent.ExternalRef),
		Status:        entity.PaymentStatusSuccess,
		Duplicate:     duplicate,
		LedgerPosted:  false,
	}, nil
}

func (s *PaymentService) finalizeSuccessfulPayment(ctx context.Context, intent *entity.PaymentIntent, result *providers.PaymentResult) (bool, error) {
	idempotencyKey := paymentIdempotencyKey(intent, result)
	processed, err := s.intentStore.IsProcessed(ctx, intent.Provider, idempotencyKey)
	if err != nil {
		return false, err
	}
	if processed {
		if err := s.intentStore.UpdateIntentProviderState(ctx, intent.TransactionID, result.ExternalRef, entity.PaymentStatusSuccess); err != nil {
			return false, err
		}
		return true, nil
	}

	if err := s.intentStore.WithTransaction(ctx, func(store PaymentIntentStore) error {
		return s.finalizeSuccessfulPaymentTx(ctx, store, intent, result)
	}); err != nil {
		if errors.Is(err, paymentrepos.ErrProviderPaymentDuplicateProcessed) {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func (s *PaymentService) finalizeSuccessfulPaymentTx(ctx context.Context, store PaymentIntentStore, intent *entity.PaymentIntent, result *providers.PaymentResult) error {
	idempotencyKey := paymentIdempotencyKey(intent, result)

	if err := store.MarkProcessed(ctx, &entity.ProcessedPaymentEvent{
		Provider:       intent.Provider,
		IdempotencyKey: idempotencyKey,
		TransactionID:  intent.TransactionID,
		CreatedAt:      time.Now().UTC(),
	}); err != nil {
		if errors.Is(err, paymentrepos.ErrProviderPaymentDuplicateProcessed) {
			return err
		}
		return err
	}

	if err := store.UpdateIntentProviderState(ctx, intent.TransactionID, result.ExternalRef, entity.PaymentStatusSuccess); err != nil {
		return err
	}

	return store.AppendOutboxEvent(ctx, newPaymentSucceededEvent(intent, result))
}

func (s *PaymentService) findIntent(ctx context.Context, provider string, result *providers.PaymentResult) (*entity.PaymentIntent, error) {
	if strings.TrimSpace(result.TransactionID) != "" {
		intent, err := s.intentStore.GetIntentByTransactionID(ctx, result.TransactionID)
		if err == nil {
			return intent, nil
		}
		if !errors.Is(err, paymentrepos.ErrProviderPaymentNotFound) {
			return nil, err
		}
	}

	if strings.TrimSpace(result.ExternalRef) != "" {
		intent, err := s.intentStore.GetIntentByExternalRef(ctx, provider, result.ExternalRef)
		if err == nil {
			return intent, nil
		}
		if !errors.Is(err, paymentrepos.ErrProviderPaymentNotFound) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("%v: transaction_id=%s external_ref=%s", ErrPaymentIntentNotFound, result.TransactionID, result.ExternalRef)
}

func (s *PaymentService) updateIntentStatus(ctx context.Context, transactionID, status string) error {
	return s.intentStore.UpdateIntentStatus(ctx, transactionID, status)
}

func validateResultAgainstIntent(intent *entity.PaymentIntent, result *providers.PaymentResult) error {
	if result.Amount != 0 && result.Amount != intent.Amount {
		return fmt.Errorf("%v: provider amount does not match reserved payment", ErrValidation)
	}
	if currency := strings.TrimSpace(result.Currency); currency != "" && !strings.EqualFold(currency, intent.Currency) {
		return fmt.Errorf("%v: provider currency does not match reserved payment", ErrValidation)
	}
	return nil
}

func normalizePaymentStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case entity.PaymentStatusSuccess:
		return entity.PaymentStatusSuccess
	case entity.PaymentStatusFailed:
		return entity.PaymentStatusFailed
	case entity.PaymentStatusCreating:
		return entity.PaymentStatusCreating
	case entity.PaymentStatusPending:
		return entity.PaymentStatusPending
	default:
		return entity.PaymentStatusPending
	}
}

func paymentIdempotencyKey(intent *entity.PaymentIntent, result *providers.PaymentResult) string {
	if strings.TrimSpace(result.EventID) != "" {
		return result.EventID
	}
	if strings.TrimSpace(result.ExternalRef) != "" {
		return result.ExternalRef
	}
	if strings.TrimSpace(intent.ExternalRef) != "" {
		return intent.ExternalRef
	}
	return intent.TransactionID
}

func coalesce(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func wrapValidation(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%v: %s", ErrValidation, err.Error())
}

func newPaymentCreatedEvent(intent *entity.PaymentIntent, metadata map[string]string) eventpkg.Event {
	return eventpkg.Event{
		AggregateID:   intent.TransactionID,
		AggregateType: paymentAggregateType,
		Version:       1,
		EventName:     sharedevents.EventPaymentCreated,
		EventData: sharedevents.PaymentCreatedEvent{
			PaymentID:       intent.TransactionID,
			TransactionID:   intent.TransactionID,
			Provider:        intent.Provider,
			Amount:          intent.Amount,
			Currency:        intent.Currency,
			DebitAccountID:  intent.DebitAccountID,
			CreditAccountID: intent.CreditAccountID,
			Status:          intent.Status,
			Metadata:        metadata,
			CreatedAt:       intent.CreatedAt,
		},
		CreatedAt: time.Now().Unix(),
	}
}

func newPaymentCheckoutSessionCreatedEvent(intent *entity.PaymentIntent, response *providers.CreatePaymentResponse, status string) eventpkg.Event {
	return eventpkg.Event{
		AggregateID:   intent.TransactionID,
		AggregateType: paymentAggregateType,
		Version:       1,
		EventName:     sharedevents.EventPaymentCheckoutSessionCreated,
		EventData: sharedevents.PaymentCheckoutSessionCreatedEvent{
			PaymentID:          intent.TransactionID,
			TransactionID:      intent.TransactionID,
			Provider:           intent.Provider,
			ProviderPaymentRef: response.ExternalRef,
			CheckoutURL:        response.CheckoutURL,
			Amount:             intent.Amount,
			Currency:           intent.Currency,
			Status:             status,
			OccurredAt:         time.Now().UTC(),
		},
		CreatedAt: time.Now().Unix(),
	}
}

func newPaymentSucceededEvent(intent *entity.PaymentIntent, result *providers.PaymentResult) eventpkg.Event {
	return eventpkg.Event{
		AggregateID:   intent.TransactionID,
		AggregateType: paymentAggregateType,
		Version:       1,
		EventName:     sharedevents.EventPaymentSucceeded,
		EventData: sharedevents.PaymentSucceededEvent{
			PaymentID:          intent.TransactionID,
			TransactionID:      intent.TransactionID,
			Provider:           intent.Provider,
			ProviderEventID:    result.EventID,
			ProviderEventType:  result.EventType,
			ProviderPaymentRef: coalesce(result.ExternalRef, intent.ExternalRef),
			Amount:             intent.Amount,
			Currency:           intent.Currency,
			DebitAccountID:     intent.DebitAccountID,
			CreditAccountID:    intent.CreditAccountID,
			IdempotencyKey:     fmt.Sprintf("%s:%s", sharedevents.EventPaymentSucceeded, intent.TransactionID),
			SucceededAt:        time.Now().UTC(),
		},
		CreatedAt: time.Now().Unix(),
	}
}

func newPaymentFailedEvent(intent *entity.PaymentIntent, result *providers.PaymentResult) eventpkg.Event {
	return eventpkg.Event{
		AggregateID:   intent.TransactionID,
		AggregateType: paymentAggregateType,
		Version:       1,
		EventName:     sharedevents.EventPaymentFailed,
		EventData: sharedevents.PaymentFailedEvent{
			PaymentID:          intent.TransactionID,
			TransactionID:      intent.TransactionID,
			Provider:           intent.Provider,
			ProviderEventID:    result.EventID,
			ProviderEventType:  result.EventType,
			ProviderPaymentRef: coalesce(result.ExternalRef, intent.ExternalRef),
			Amount:             intent.Amount,
			Currency:           intent.Currency,
			Status:             result.Status,
			OccurredAt:         time.Now().UTC(),
		},
		CreatedAt: time.Now().Unix(),
	}
}
