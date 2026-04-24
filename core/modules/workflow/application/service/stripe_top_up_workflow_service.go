package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"wechat-clone/core/modules/workflow/application/dto/in"
	"wechat-clone/core/modules/workflow/application/dto/out"
	"wechat-clone/core/shared/pkg/actorctx"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

var (
	ErrWorkflowUnavailable  = errors.New("workflow runner is unavailable")
	ErrWorkflowUnauthorized = errors.New("unauthorized")
	ErrWorkflowValidation   = errors.New("validation failed")
)

type WorkflowActor struct {
	AccountID string
	Email     string
	Role      string
}

type PaymentIntegrationEvent struct {
	Name     string
	DataJSON string
}

type CreateStripeTopUpWorkflowInput struct {
	Actor    WorkflowActor
	Amount   int64
	Currency string
	Metadata map[string]string
}

type StripeTopUpWorkflowResult struct {
	Provider       string
	Workflow       string
	TransactionID  string
	ExternalRef    string
	Amount         int64
	FeeAmount      int64
	ProviderAmount int64
	Status         string
	CheckoutURL    string
}

type ProcessStripeWebhookWorkflowInput struct {
	Signature string
	Payload   string
}

type StripeWebhookWorkflowResult struct {
	Provider      string
	TransactionID string
	ExternalRef   string
	Status        string
	Duplicate     bool
	LedgerPosted  bool
	Events        []PaymentIntegrationEvent
}

type StripeTopUpWorkflowRunner interface {
	CreateStripeTopUp(ctx context.Context, input CreateStripeTopUpWorkflowInput) (*StripeTopUpWorkflowResult, error)
	ProcessStripeWebhook(ctx context.Context, input ProcessStripeWebhookWorkflowInput) (*StripeWebhookWorkflowResult, error)
}

type StripeTopUpWorkflowService interface {
	CreateStripeTopUp(ctx context.Context, req *in.CreateStripeTopUpRequest) (*out.StripeTopUpResponse, error)
	ProcessStripeWebhook(ctx context.Context, req *in.ProcessStripeWebhookRequest) (*out.StripeWebhookResponse, error)
}

type stripeTopUpWorkflowService struct {
	runner StripeTopUpWorkflowRunner
}

func NewStripeTopUpWorkflowService(runner StripeTopUpWorkflowRunner) StripeTopUpWorkflowService {
	return &stripeTopUpWorkflowService{
		runner: runner,
	}
}

func (s *stripeTopUpWorkflowService) CreateStripeTopUp(ctx context.Context, req *in.CreateStripeTopUpRequest) (*out.StripeTopUpResponse, error) {
	log := logging.FromContext(ctx)
	if s == nil || s.runner == nil {
		return nil, stackErr.Error(ErrWorkflowUnavailable)
	}

	actor, ok := actorctx.FromContext(ctx)
	if !ok {
		return nil, stackErr.Error(ErrWorkflowUnauthorized)
	}
	if req == nil {
		return nil, stackErr.Error(fmt.Errorf("%w: request is required", ErrWorkflowValidation))
	}

	result, err := s.runner.CreateStripeTopUp(ctx, CreateStripeTopUpWorkflowInput{
		Actor: WorkflowActor{
			AccountID: actor.AccountID,
			Email:     actor.Email,
			Role:      actor.Role,
		},
		Amount:   req.Amount,
		Currency: strings.TrimSpace(req.Currency),
		Metadata: cloneMetadata(req.Metadata),
	})
	if err != nil {
		log.Errorw("CreateStripeTopUp", zap.Any("payload", CreateStripeTopUpWorkflowInput{
			Actor: WorkflowActor{
				AccountID: actor.AccountID,
				Email:     actor.Email,
				Role:      actor.Role,
			},
			Amount:   req.Amount,
			Currency: strings.TrimSpace(req.Currency),
			Metadata: cloneMetadata(req.Metadata),
		}))
		return nil, stackErr.Error(err)
	}

	return &out.StripeTopUpResponse{
		Provider:       result.Provider,
		Workflow:       result.Workflow,
		TransactionID:  result.TransactionID,
		ExternalRef:    result.ExternalRef,
		Amount:         result.Amount,
		FeeAmount:      result.FeeAmount,
		ProviderAmount: result.ProviderAmount,
		Status:         result.Status,
		CheckoutURL:    result.CheckoutURL,
	}, nil
}

func (s *stripeTopUpWorkflowService) ProcessStripeWebhook(ctx context.Context, req *in.ProcessStripeWebhookRequest) (*out.StripeWebhookResponse, error) {
	if s == nil || s.runner == nil {
		return nil, stackErr.Error(ErrWorkflowUnavailable)
	}
	if req == nil {
		return nil, stackErr.Error(fmt.Errorf("%w: request is required", ErrWorkflowValidation))
	}

	result, err := s.runner.ProcessStripeWebhook(ctx, ProcessStripeWebhookWorkflowInput{
		Signature: strings.TrimSpace(req.Signature),
		Payload:   req.Payload,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return &out.StripeWebhookResponse{
		Provider:      result.Provider,
		TransactionID: result.TransactionID,
		ExternalRef:   result.ExternalRef,
		Status:        result.Status,
		Duplicate:     result.Duplicate,
		LedgerPosted:  result.LedgerPosted,
	}, nil
}

func cloneMetadata(metadata map[string]string) map[string]string {
	if len(metadata) == 0 {
		return map[string]string{}
	}

	cloned := make(map[string]string, len(metadata))
	for key, value := range metadata {
		cloned[key] = strings.TrimSpace(value)
	}

	return cloned
}
