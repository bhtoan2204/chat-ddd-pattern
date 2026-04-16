package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-socket/core/modules/payment/application/dto/in"
	"go-socket/core/modules/payment/domain/entity"
	repos "go-socket/core/modules/payment/domain/repos"
	domainservice "go-socket/core/modules/payment/domain/service"
	sharedevents "go-socket/core/shared/contracts/events"
	sharedlock "go-socket/core/shared/infra/lock"
	"go-socket/core/shared/pkg/actorctx"
	eventpkg "go-socket/core/shared/pkg/event"

	"go.uber.org/mock/gomock"
)

func TestCreatePaymentRejectsUnauthorizedOrCrossAccountRequests(t *testing.T) {
	t.Run("rejects credit account that does not match authenticated actor", func(t *testing.T) {
		svc := &paymentCommandService{}

		_, err := svc.CreatePayment(
			actorctx.WithActor(context.Background(), actorctx.Actor{AccountID: "acc-user-1"}),
			&in.CreatePaymentRequest{
				Provider:        "stripe",
				Amount:          100,
				Currency:        "VND",
				CreditAccountID: "acc-user-2",
			},
		)
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("expected validation error, got %v", err)
		}
	})

	t.Run("requires authenticated actor for create payment", func(t *testing.T) {
		svc := &paymentCommandService{}

		_, err := svc.CreatePayment(context.Background(), &in.CreatePaymentRequest{
			Provider: "stripe",
			Amount:   100,
			Currency: "VND",
		})
		if !errors.Is(err, ErrPaymentUnauthorized) {
			t.Fatalf("expected unauthorized error, got %v", err)
		}
	})
}

func TestProcessWebhookLocksByTransactionIDAndFinalizesSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	baseRepo := repos.NewMockRepos(ctrl)
	txRepos := repos.NewMockRepos(ctrl)
	providerRepo := repos.NewMockProviderPaymentRepository(ctrl)
	providerRegistry := domainservice.NewMockPaymentProviderRegistry(ctrl)
	provider := domainservice.NewMockPaymentProvider(ctrl)
	locker := sharedlock.NewMockLock(ctrl)

	intent, err := entity.NewProviderTopUpIntent("txn-1", "stripe", 100, "VND", "wallet:available", time.Now().UTC())
	if err != nil {
		t.Fatalf("new provider top up intent: %v", err)
	}

	providerRegistry.EXPECT().Get("stripe").Return(provider, nil)
	provider.EXPECT().ParseWebhook(gomock.Any(), []byte("{}"), "sig-1").Return(&domainservice.PaymentWebhook{
		Provider: "stripe",
		Result: entity.PaymentProviderResult{
			TransactionID: "txn-1",
			EventID:       "evt-1",
			EventType:     "checkout.session.completed",
			Status:        entity.PaymentStatusSuccess,
			Amount:        100,
			Currency:      "VND",
			ExternalRef:   "cs-1",
		},
	}, nil)

	baseRepo.EXPECT().ProviderPaymentRepository().Return(providerRepo).AnyTimes()
	providerRepo.EXPECT().GetIntentByTransactionID(gomock.Any(), "txn-1").Return(intent, nil).Times(2)
	locker.EXPECT().AcquireLock(gomock.Any(), "payment:txn-1", gomock.Any(), 30*time.Second, 100*time.Millisecond, 3*time.Second).Return(true, nil)
	locker.EXPECT().ReleaseLock(gomock.Any(), "payment:txn-1", gomock.Any()).Return(true, nil)
	baseRepo.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(repos.Repos) error) error {
		return fn(txRepos)
	})
	txRepos.EXPECT().ProviderPaymentRepository().Return(providerRepo).AnyTimes()
	providerRepo.EXPECT().FinalizeSuccessfulPayment(
		gomock.Any(),
		intent,
		gomock.AssignableToTypeOf(&entity.ProcessedPaymentEvent{}),
		gomock.AssignableToTypeOf(eventpkg.Event{}),
	).DoAndReturn(func(_ context.Context, savedIntent *entity.PaymentIntent, processedEvent *entity.ProcessedPaymentEvent, successEvent eventpkg.Event, _ ...eventpkg.Event) error {
		if savedIntent.Status != entity.PaymentStatusSuccess {
			t.Fatalf("expected saved intent status success, got %s", savedIntent.Status)
		}
		if processedEvent.IdempotencyKey != "payment.succeeded:txn-1" {
			t.Fatalf("unexpected processed event idempotency key: %s", processedEvent.IdempotencyKey)
		}
		if successEvent.EventName != sharedevents.EventPaymentSucceeded {
			t.Fatalf("unexpected success event name: %s", successEvent.EventName)
		}
		payload, ok := successEvent.EventData.(sharedevents.PaymentSucceededEvent)
		if !ok {
			t.Fatalf("unexpected success payload type: %T", successEvent.EventData)
		}
		if payload.IdempotencyKey != "payment.succeeded:txn-1" {
			t.Fatalf("unexpected success payload idempotency key: %s", payload.IdempotencyKey)
		}
		return nil
	})

	svc := &paymentCommandService{
		baseRepo:         baseRepo,
		locker:           locker,
		providerRegistry: providerRegistry,
	}

	response, err := svc.ProcessWebhook(context.Background(), &in.ProcessWebhookRequest{
		Provider:  "stripe",
		Signature: "sig-1",
		Payload:   "{}",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response.TransactionID != "txn-1" {
		t.Fatalf("unexpected transaction id: %s", response.TransactionID)
	}
	if response.Status != entity.PaymentStatusSuccess {
		t.Fatalf("unexpected status: %s", response.Status)
	}
	if response.Duplicate {
		t.Fatalf("expected non-duplicate webhook processing")
	}
}

func TestApplyProviderOutcomeFinalizesSuccessOnlyOncePerPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	baseRepo := repos.NewMockRepos(ctrl)
	txRepos := repos.NewMockRepos(ctrl)
	providerRepo := repos.NewMockProviderPaymentRepository(ctrl)

	intent, err := entity.NewProviderTopUpIntent("txn-1", "stripe", 100, "VND", "wallet:available", time.Now().UTC())
	if err != nil {
		t.Fatalf("new provider top up intent: %v", err)
	}

	baseRepo.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(repos.Repos) error) error {
		return fn(txRepos)
	}).Times(2)
	txRepos.EXPECT().ProviderPaymentRepository().Return(providerRepo).AnyTimes()
	providerRepo.EXPECT().FinalizeSuccessfulPayment(
		gomock.Any(),
		intent,
		gomock.AssignableToTypeOf(&entity.ProcessedPaymentEvent{}),
		gomock.AssignableToTypeOf(eventpkg.Event{}),
	).DoAndReturn(func(_ context.Context, _ *entity.PaymentIntent, processedEvent *entity.ProcessedPaymentEvent, successEvent eventpkg.Event, _ ...eventpkg.Event) error {
		if processedEvent.IdempotencyKey != "payment.succeeded:txn-1" {
			t.Fatalf("unexpected processed event idempotency key: %s", processedEvent.IdempotencyKey)
		}
		if successEvent.EventName != sharedevents.EventPaymentSucceeded {
			t.Fatalf("unexpected success event name: %s", successEvent.EventName)
		}
		return nil
	}).Times(1)

	svc := &paymentCommandService{baseRepo: baseRepo}
	duplicate, err := svc.applyProviderOutcome(context.Background(), intent, entity.PaymentProviderResult{
		TransactionID: "txn-1",
		EventID:       "evt-checkout-completed",
		EventType:     "checkout.session.completed",
		Status:        entity.PaymentStatusSuccess,
		Amount:        100,
		Currency:      "VND",
		ExternalRef:   "cs-1",
	}, "", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if duplicate {
		t.Fatalf("expected first success to finalize payment")
	}
	if intent.Status != entity.PaymentStatusSuccess {
		t.Fatalf("expected success status, got %s", intent.Status)
	}

	duplicate, err = svc.applyProviderOutcome(context.Background(), intent, entity.PaymentProviderResult{
		TransactionID: "txn-1",
		EventID:       "evt-charge-succeeded",
		EventType:     "charge.succeeded",
		Status:        entity.PaymentStatusSuccess,
		Amount:        100,
		Currency:      "VND",
		ExternalRef:   "cs-1",
	}, "", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !duplicate {
		t.Fatalf("expected second success to be treated as duplicate")
	}
	if intent.Status != entity.PaymentStatusSuccess {
		t.Fatalf("expected status to stay success, got %s", intent.Status)
	}
}

func TestApplyProviderOutcomeIgnoresFailAfterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	baseRepo := repos.NewMockRepos(ctrl)

	intent, err := entity.NewProviderTopUpIntent("txn-1", "stripe", 100, "VND", "wallet:available", time.Now().UTC())
	if err != nil {
		t.Fatalf("new provider top up intent: %v", err)
	}
	if err := intent.ApplyProviderResult(entity.PaymentProviderResult{
		ExternalRef: "cs-1",
		Status:      entity.PaymentStatusSuccess,
		Amount:      100,
		Currency:    "VND",
	}, time.Now().UTC()); err != nil {
		t.Fatalf("apply initial success: %v", err)
	}

	baseRepo.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(repos.Repos) error) error {
		return fn(repos.NewMockRepos(ctrl))
	})

	svc := &paymentCommandService{baseRepo: baseRepo}
	duplicate, err := svc.applyProviderOutcome(context.Background(), intent, entity.PaymentProviderResult{
		TransactionID: "txn-1",
		EventID:       "evt-payment-failed",
		EventType:     "payment_intent.payment_failed",
		Status:        entity.PaymentStatusFailed,
		Amount:        100,
		Currency:      "VND",
		ExternalRef:   "cs-1",
	}, "", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !duplicate {
		t.Fatalf("expected fail-after-success to be ignored as duplicate")
	}
	if intent.Status != entity.PaymentStatusSuccess {
		t.Fatalf("expected status to stay success, got %s", intent.Status)
	}
}

func TestApplyProviderOutcomeFinalizesRefundAsReversal(t *testing.T) {
	ctrl := gomock.NewController(t)
	baseRepo := repos.NewMockRepos(ctrl)
	txRepos := repos.NewMockRepos(ctrl)
	providerRepo := repos.NewMockProviderPaymentRepository(ctrl)

	intent, err := entity.NewProviderTopUpIntent("txn-1", "stripe", 100, "VND", "wallet:available", time.Now().UTC())
	if err != nil {
		t.Fatalf("new provider top up intent: %v", err)
	}
	if err := intent.ApplyProviderResult(entity.PaymentProviderResult{
		ExternalRef: "cs-1",
		Status:      entity.PaymentStatusSuccess,
		Amount:      100,
		Currency:    "VND",
	}, time.Now().UTC()); err != nil {
		t.Fatalf("apply initial success: %v", err)
	}

	baseRepo.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(repos.Repos) error) error {
		return fn(txRepos)
	})
	txRepos.EXPECT().ProviderPaymentRepository().Return(providerRepo).AnyTimes()
	providerRepo.EXPECT().FinalizeReversedPayment(
		gomock.Any(),
		intent,
		gomock.AssignableToTypeOf(&entity.ProcessedPaymentEvent{}),
		gomock.AssignableToTypeOf(eventpkg.Event{}),
	).DoAndReturn(func(_ context.Context, savedIntent *entity.PaymentIntent, processedEvent *entity.ProcessedPaymentEvent, reversalEvent eventpkg.Event, _ ...eventpkg.Event) error {
		if savedIntent.Status != entity.PaymentStatusRefunded {
			t.Fatalf("expected refunded status, got %s", savedIntent.Status)
		}
		if processedEvent.IdempotencyKey != "payment.refunded:txn-1" {
			t.Fatalf("unexpected processed event idempotency key: %s", processedEvent.IdempotencyKey)
		}
		if reversalEvent.EventName != sharedevents.EventPaymentRefunded {
			t.Fatalf("unexpected reversal event name: %s", reversalEvent.EventName)
		}
		payload, ok := reversalEvent.EventData.(sharedevents.PaymentRefundedEvent)
		if !ok {
			t.Fatalf("unexpected reversal payload type: %T", reversalEvent.EventData)
		}
		if payload.IdempotencyKey != "payment.refunded:txn-1" {
			t.Fatalf("unexpected reversal payload idempotency key: %s", payload.IdempotencyKey)
		}
		return nil
	})

	svc := &paymentCommandService{baseRepo: baseRepo}
	duplicate, err := svc.applyProviderOutcome(context.Background(), intent, entity.PaymentProviderResult{
		TransactionID: "txn-1",
		EventID:       "evt-charge-refunded",
		EventType:     "charge.refunded",
		Status:        entity.PaymentStatusRefunded,
		Amount:        100,
		Currency:      "VND",
		ExternalRef:   "cs-1",
	}, "", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if duplicate {
		t.Fatalf("expected first refund to finalize reversal")
	}
	if intent.Status != entity.PaymentStatusRefunded {
		t.Fatalf("expected refunded status, got %s", intent.Status)
	}
}
