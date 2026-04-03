package service

import (
	"context"
	"testing"
	"time"

	ledgerin "go-socket/core/modules/ledger/application/dto/in"
	"go-socket/core/modules/ledger/domain/entity"
	ledgerrepos "go-socket/core/modules/ledger/domain/repos"
	ledgerrepo "go-socket/core/modules/ledger/infra/persistent/repository"
	"go-socket/core/modules/ledger/providers"
)

func TestPaymentServiceCreatePaymentDoesNotRefetchIntent(t *testing.T) {
	repos := &paymentServiceTestRepos{
		paymentRepo: &paymentServiceTestPaymentRepo{},
	}

	registry := providers.NewProviderRegistry()
	registry.Register(&paymentServiceTestProvider{
		name: "stripe",
		response: &providers.CreatePaymentResponse{
			Provider:      "stripe",
			TransactionID: "tx-001",
			ExternalRef:   "cs_test_123",
			Status:        entity.PaymentStatusPending,
			CheckoutURL:   "https://checkout.stripe.test/cs_test_123",
		},
	})

	service := NewPaymentService(repos, nil, registry)

	response, err := service.CreatePayment(context.Background(), &ledgerin.CreatePaymentRequest{
		Provider:        "stripe",
		TransactionID:   "tx-001",
		Amount:          10000,
		Currency:        "USD",
		DebitAccountID:  "acc-debit",
		CreditAccountID: "acc-credit",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.ExternalRef != "cs_test_123" {
		t.Fatalf("expected external ref cs_test_123, got %s", response.ExternalRef)
	}
	if repos.paymentRepo.getIntentCalls != 0 {
		t.Fatalf("expected GetIntentByTransactionID to not be called, got %d calls", repos.paymentRepo.getIntentCalls)
	}
	if repos.paymentRepo.updatedTransactionID != "tx-001" {
		t.Fatalf("expected updated transaction id tx-001, got %s", repos.paymentRepo.updatedTransactionID)
	}
	if repos.paymentRepo.updatedExternalRef != "cs_test_123" {
		t.Fatalf("expected updated external ref cs_test_123, got %s", repos.paymentRepo.updatedExternalRef)
	}
	if repos.paymentRepo.updatedStatus != entity.PaymentStatusPending {
		t.Fatalf("expected updated status PENDING, got %s", repos.paymentRepo.updatedStatus)
	}
}

type paymentServiceTestRepos struct {
	paymentRepo *paymentServiceTestPaymentRepo
}

func (r *paymentServiceTestRepos) LedgerRepository() ledgerrepos.LedgerRepository {
	return nil
}

func (r *paymentServiceTestRepos) PaymentRepository() ledgerrepos.PaymentRepository {
	return r.paymentRepo
}

func (r *paymentServiceTestRepos) WithTransaction(_ context.Context, fn func(ledgerrepos.Repos) error) error {
	return fn(r)
}

type paymentServiceTestPaymentRepo struct {
	intent               *entity.PaymentIntent
	getIntentCalls       int
	updatedTransactionID string
	updatedExternalRef   string
	updatedStatus        string
}

func (r *paymentServiceTestPaymentRepo) CreateIntent(_ context.Context, intent *entity.PaymentIntent) error {
	copied := *intent
	r.intent = &copied
	return nil
}

func (r *paymentServiceTestPaymentRepo) GetIntentByTransactionID(_ context.Context, _ string) (*entity.PaymentIntent, error) {
	r.getIntentCalls++
	return nil, ledgerrepo.ErrNotFound
}

func (r *paymentServiceTestPaymentRepo) GetIntentByExternalRef(_ context.Context, _, _ string) (*entity.PaymentIntent, error) {
	return nil, ledgerrepo.ErrNotFound
}

func (r *paymentServiceTestPaymentRepo) UpdateIntentProviderState(_ context.Context, transactionID, externalRef, status string) error {
	if r.intent == nil || r.intent.TransactionID != transactionID {
		return ledgerrepo.ErrNotFound
	}
	r.updatedTransactionID = transactionID
	r.updatedExternalRef = externalRef
	r.updatedStatus = status
	r.intent.ExternalRef = externalRef
	r.intent.Status = status
	r.intent.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *paymentServiceTestPaymentRepo) UpdateIntentStatus(_ context.Context, transactionID, status string) error {
	if r.intent == nil || r.intent.TransactionID != transactionID {
		return ledgerrepo.ErrNotFound
	}
	r.intent.Status = status
	r.intent.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *paymentServiceTestPaymentRepo) IsProcessed(_ context.Context, _, _ string) (bool, error) {
	return false, nil
}

func (r *paymentServiceTestPaymentRepo) MarkProcessed(_ context.Context, _ *entity.ProcessedPaymentEvent) error {
	return nil
}

type paymentServiceTestProvider struct {
	name     string
	response *providers.CreatePaymentResponse
	err      error
}

func (p *paymentServiceTestProvider) Name() string {
	return p.name
}

func (p *paymentServiceTestProvider) CreatePayment(_ context.Context, _ providers.CreatePaymentRequest) (*providers.CreatePaymentResponse, error) {
	return p.response, p.err
}

func (p *paymentServiceTestProvider) VerifyWebhook(_ context.Context, _ []byte, _ string) (*providers.WebhookEvent, error) {
	return nil, nil
}

func (p *paymentServiceTestProvider) ParseEvent(_ context.Context, _ *providers.WebhookEvent) (*providers.PaymentResult, error) {
	return nil, nil
}
