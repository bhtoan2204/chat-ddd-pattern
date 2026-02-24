package usecase

import (
	"context"
	"errors"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/logging"
	"time"

	"go.uber.org/zap"
)

type accountUsecaseImpl struct {
	baseRepo repos.Repos
}

func NewAccountUsecase(appCtx *appCtx.AppContext, baseRepo repos.Repos) AccountUsecase {
	_ = appCtx
	return &accountUsecaseImpl{
		baseRepo: baseRepo,
	}
}

func (u *accountUsecaseImpl) GetProfile(ctx context.Context, req *in.GetProfileRequest) (*out.GetProfileResponse, error) {
	_ = req
	log := logging.FromContext(ctx).Named("GetProfile")
	account := ctx.Value("account")
	if account == nil {
		log.Errorw("Account not found", zap.Error(errors.New("account not found")))
		return nil, errors.New("account not found")
	}

	payload, ok := account.(*xpaseto.PasetoPayload)
	if !ok {
		return nil, errors.New("invalid account payload")
	}

	accountEntity, err := u.baseRepo.AccountRepository().GetAccountByID(ctx, payload.AccountID)
	if err != nil {
		log.Errorw("Failed to get account by ID", zap.Error(err))
		return nil, err
	}

	return &out.GetProfileResponse{
		Email:     accountEntity.Email,
		CreatedAt: accountEntity.CreatedAt.Format(time.RFC3339),
		UpdatedAt: accountEntity.UpdatedAt.Format(time.RFC3339),
	}, nil
}
