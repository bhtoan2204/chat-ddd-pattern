package usecase

import (
	"context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
)

type AccountUsecase interface {
	GetProfile(ctx context.Context, in *in.GetProfileRequest) (*out.GetProfileResponse, error)
}
