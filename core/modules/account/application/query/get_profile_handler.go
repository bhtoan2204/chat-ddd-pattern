package query

import (
	"context"
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/account/application/dto/in"
	"wechat-clone/core/modules/account/application/dto/out"
	"wechat-clone/core/modules/account/application/service"
	"wechat-clone/core/modules/account/application/support"
	repos "wechat-clone/core/modules/account/domain/repos"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

type getProfileHandler struct {
	accountRepo repos.AccountRepository
}

func NewGetProfileHandler(appCtx *appCtx.AppContext, baseRepo repos.Repos, services service.Services) cqrs.Handler[*in.GetProfileRequest, *out.GetProfileResponse] {
	return &getProfileHandler{
		accountRepo: baseRepo.AccountRepository(),
	}
}

func (u *getProfileHandler) Handle(ctx context.Context, req *in.GetProfileRequest) (*out.GetProfileResponse, error) {
	_ = req
	log := logging.FromContext(ctx).Named("GetProfile")
	accountID, err := support.AccountIDFromCtx(ctx)
	if err != nil {
		log.Errorw("Account not found in context", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	accountEntity, err := u.accountRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		log.Errorw("Failed to get account by ID", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	return support.ToGetProfileResponse(accountEntity), nil
}
