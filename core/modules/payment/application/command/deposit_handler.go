package command

import (
	"context"
	"errors"
	"reflect"
	"time"

	"go-socket/core/modules/payment/application/dto/in"
	"go-socket/core/modules/payment/application/dto/out"
	"go-socket/core/modules/payment/domain/aggregate"
	"go-socket/core/modules/payment/domain/repos"
	paymentrepos "go-socket/core/modules/payment/domain/repos"
	"go-socket/core/shared/infra/xpaseto"
	eventpkg "go-socket/core/shared/pkg/event"
	"go-socket/core/shared/pkg/logging"
	stackerr "go-socket/core/shared/pkg/stackErr"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type depositHandler struct {
	outboxRepo paymentrepos.PaymentOutboxEventsRepository
}

func NewDepositHandler(repos repos.Repos) DepositHandler {
	return &depositHandler{
		outboxRepo: repos.PaymentOutboxEventsRepository(),
	}
}

func (h *depositHandler) Handle(ctx context.Context, req *in.DepositRequest) (*out.DepositResponse, error) {
	log := logging.FromContext(ctx).Named("Deposit")

	account, ok := ctx.Value("account").(*xpaseto.PasetoPayload)
	if !ok || account == nil || account.AccountID == "" {
		log.Errorw("Account not found")
		return nil, stackerr.Error(errors.New("account not found"))
	}

	now := time.Now().UTC()
	transactionID := uuid.NewString()

	agg := &aggregate.PaymentTransactionAggregate{}
	aggType := reflect.TypeOf(agg).Elem().Name()
	agg.SetAggregateType(aggType)
	if err := agg.SetID(transactionID); err != nil {
		log.Errorw("Failed to set aggregate id", zap.Error(err))
		return nil, stackerr.Error(errors.New("failed to set aggregate id"))
	}

	if err := agg.ApplyChange(agg, &aggregate.EventPaymentTransactionDeposited{
		PaymentTransactionID:         transactionID,
		PaymentTransactionAmount:     req.Amount,
		PaymentTransactionReceiverID: account.AccountID,
		PaymentTransactionCreatedAt:  now,
		PaymentTransactionUpdatedAt:  now,
	}); err != nil {
		log.Errorw("Failed to apply deposit event", zap.Error(err))
		return nil, stackerr.Error(err)
	}

	publisher := eventpkg.NewPublisher(h.outboxRepo)
	if err := publisher.PublishAggregate(ctx, agg); err != nil {
		log.Errorw("Failed to publish deposit event", zap.Error(err))
		return nil, stackerr.Error(err)
	}

	return &out.DepositResponse{
		Message: "Deposit successful",
	}, nil
}
