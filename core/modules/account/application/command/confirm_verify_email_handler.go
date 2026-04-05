package command

import (
	"context"
	"time"

	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/application/service"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"
	"go-socket/core/shared/utils"

	"go.uber.org/zap"
)

type confirmVerifyEmailHandler struct {
	baseRepo            repos.Repos
	aggregateService    *service.AccountAggregateService
	verificationService *service.EmailVerificationService
}

func NewConfirmVerifyEmailHandler(baseRepo repos.Repos, aggregateService *service.AccountAggregateService, verificationService *service.EmailVerificationService) cqrs.Handler[*in.ConfirmVerifyEmailRequest, *out.ConfirmVerifyEmailResponse] {
	return &confirmVerifyEmailHandler{
		baseRepo:            baseRepo,
		aggregateService:    aggregateService,
		verificationService: verificationService,
	}
}

func (u *confirmVerifyEmailHandler) Handle(ctx context.Context, req *in.ConfirmVerifyEmailRequest) (*out.ConfirmVerifyEmailResponse, error) {
	log := logging.FromContext(ctx).Named("ConfirmVerifyEmail")

	tokenPayload, err := u.verificationService.ConsumeVerificationToken(ctx, req.Token)
	if err != nil {
		log.Errorw("Failed to consume verification token", zap.Error(err))
		return nil, stackErr.Error(ErrInvalidVerificationToken)
	}

	accountEntity, err := u.baseRepo.AccountRepository().GetAccountByID(ctx, tokenPayload.AccountID)
	if err != nil {
		log.Errorw("Failed to get account by ID", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if accountEntity.Email.Value() != tokenPayload.Email {
		return nil, stackErr.Error(ErrInvalidVerificationToken)
	}

	updated, err := accountEntity.MarkEmailVerified(utils.NowUTC())
	if err != nil {
		log.Errorw("Failed to mark email verified", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if !updated || accountEntity.EmailVerifiedAt == nil {
		return nil, stackErr.Error(ErrInvalidVerificationToken)
	}

	if txErr := u.baseRepo.WithTransaction(ctx, func(txRepos repos.Repos) error {
		if err := txRepos.AccountRepository().UpdateAccount(ctx, accountEntity); err != nil {
			return stackErr.Error(err)
		}
		return u.aggregateService.PublishEmailVerified(ctx, txRepos.AccountOutboxEventsRepository(), accountEntity)
	}); txErr != nil {
		log.Errorw("Failed to persist verified email", zap.Error(txErr))
		return nil, stackErr.Error(txErr)
	}

	return &out.ConfirmVerifyEmailResponse{
		Message:    "Email verified successfully",
		VerifiedAt: accountEntity.EmailVerifiedAt.UTC().Format(time.RFC3339),
	}, nil
}
