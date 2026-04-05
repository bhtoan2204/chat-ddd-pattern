package command

import (
	"context"

	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/application/service"
	"go-socket/core/modules/account/application/support"
	"go-socket/core/modules/account/domain/entity"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"
	"go-socket/core/shared/utils"

	"go.uber.org/zap"
)

type verifyEmailHandler struct {
	baseRepo            repos.Repos
	aggregateService    *service.AccountAggregateService
	verificationService *service.EmailVerificationService
}

func NewVerifyEmailHandler(baseRepo repos.Repos, aggregateService *service.AccountAggregateService, verificationService *service.EmailVerificationService) cqrs.Handler[*in.VerifyEmailRequest, *out.VerifyEmailResponse] {
	return &verifyEmailHandler{
		baseRepo:            baseRepo,
		aggregateService:    aggregateService,
		verificationService: verificationService,
	}
}

func (u *verifyEmailHandler) Handle(ctx context.Context, req *in.VerifyEmailRequest) (*out.VerifyEmailResponse, error) {
	_ = req
	log := logging.FromContext(ctx).Named("VerifyEmail")

	accountID, err := support.AccountIDFromCtx(ctx)
	if err != nil {
		log.Errorw("Failed to resolve account from context", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	accountEntity, err := u.baseRepo.AccountRepository().GetAccountByID(ctx, accountID)
	if err != nil {
		log.Errorw("Failed to get account by ID", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if accountEntity.EmailVerifiedAt != nil {
		return nil, stackErr.Error(entity.ErrAccountAlreadyVerified)
	}

	requestedAt := utils.NowUTC()
	token, _, err := u.verificationService.SendVerificationEmail(ctx, accountEntity, requestedAt)
	if err != nil {
		log.Errorw("Failed to send verification email", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	if txErr := u.baseRepo.WithTransaction(ctx, func(txRepos repos.Repos) error {
		return u.aggregateService.PublishEmailVerificationRequested(ctx, txRepos.AccountOutboxEventsRepository(), accountEntity, token, requestedAt)
	}); txErr != nil {
		log.Errorw("Failed to publish verification requested event", zap.Error(txErr))
		return nil, stackErr.Error(txErr)
	}

	return &out.VerifyEmailResponse{
		Message: "Verification email queued",
	}, nil
}
