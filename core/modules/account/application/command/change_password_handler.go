package command

import (
	"context"

	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/application/service"
	"go-socket/core/modules/account/application/support"
	"go-socket/core/modules/account/domain/entity"
	repos "go-socket/core/modules/account/domain/repos"
	valueobject "go-socket/core/modules/account/domain/value_object"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/hasher"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"
	"go-socket/core/shared/utils"

	"go.uber.org/zap"
)

type changePasswordHandler struct {
	baseRepo         repos.Repos
	hasher           hasher.Hasher
	aggregateService *service.AccountAggregateService
}

func NewChangePasswordHandler(appContext *appCtx.AppContext, baseRepo repos.Repos, aggregateService *service.AccountAggregateService) cqrs.Handler[*in.ChangePasswordRequest, *out.ChangePasswordResponse] {
	return &changePasswordHandler{
		baseRepo:         baseRepo,
		hasher:           appContext.GetHasher(),
		aggregateService: aggregateService,
	}
}

func (u *changePasswordHandler) Handle(ctx context.Context, req *in.ChangePasswordRequest) (*out.ChangePasswordResponse, error) {
	log := logging.FromContext(ctx).Named("ChangePassword")

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

	valid, err := u.hasher.Verify(ctx, req.CurrentPassword, accountEntity.Password.Value())
	if err != nil {
		log.Errorw("Failed to verify current password", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if !valid {
		return nil, stackErr.Error(ErrInvalidCurrentPassword)
	}
	if req.CurrentPassword == req.NewPassword {
		return nil, stackErr.Error(entity.ErrAccountPasswordSameAsOldOne)
	}

	newPassword, err := valueobject.NewPassword(req.NewPassword)
	if err != nil {
		log.Errorw("Failed to validate new password", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	hashedPassword, err := u.hasher.Hash(ctx, newPassword.Value())
	if err != nil {
		log.Errorw("Failed to hash new password", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	hashedPasswordVO, err := valueobject.NewPassword(hashedPassword)
	if err != nil {
		log.Errorw("Failed to create password value object", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	changed, err := accountEntity.ChangePassword(hashedPasswordVO, utils.NowUTC())
	if err != nil {
		log.Errorw("Failed to change password on entity", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if !changed {
		return &out.ChangePasswordResponse{Message: "Password is unchanged"}, nil
	}

	if txErr := u.baseRepo.WithTransaction(ctx, func(txRepos repos.Repos) error {
		if err := txRepos.AccountRepository().UpdateAccount(ctx, accountEntity); err != nil {
			return stackErr.Error(err)
		}
		return u.aggregateService.PublishPasswordChanged(ctx, txRepos.AccountOutboxEventsRepository(), accountEntity)
	}); txErr != nil {
		log.Errorw("Failed to persist changed password", zap.Error(txErr))
		return nil, stackErr.Error(txErr)
	}

	return &out.ChangePasswordResponse{
		Message: "Password changed successfully",
	}, nil
}
