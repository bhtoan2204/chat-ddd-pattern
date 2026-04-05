package command

import (
	"context"

	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/application/service"
	"go-socket/core/modules/account/application/support"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"
	"go-socket/core/shared/utils"

	"go.uber.org/zap"
)

type updateProfileHandler struct {
	baseRepo         repos.Repos
	aggregateService *service.AccountAggregateService
}

func NewUpdateProfileHandler(baseRepo repos.Repos, aggregateService *service.AccountAggregateService) cqrs.Handler[*in.UpdateProfileRequest, *out.UpdateProfileResponse] {
	return &updateProfileHandler{
		baseRepo:         baseRepo,
		aggregateService: aggregateService,
	}
}

func (u *updateProfileHandler) Handle(ctx context.Context, req *in.UpdateProfileRequest) (*out.UpdateProfileResponse, error) {
	log := logging.FromContext(ctx).Named("UpdateProfile")

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

	updated, err := accountEntity.UpdateProfile(req.DisplayName, req.Username, req.AvatarObjectKey, utils.NowUTC())
	if err != nil {
		log.Errorw("Failed to update account profile", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if !updated {
		return support.ToUpdateProfileResponse(accountEntity), nil
	}

	if txErr := u.baseRepo.WithTransaction(ctx, func(txRepos repos.Repos) error {
		if err := txRepos.AccountRepository().UpdateAccount(ctx, accountEntity); err != nil {
			return stackErr.Error(err)
		}
		return u.aggregateService.PublishProfileUpdated(ctx, txRepos.AccountOutboxEventsRepository(), accountEntity)
	}); txErr != nil {
		log.Errorw("Failed to persist updated profile", zap.Error(txErr))
		return nil, stackErr.Error(txErr)
	}

	return support.ToUpdateProfileResponse(accountEntity), nil
}
