package usecase

import (
	"context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
)

type AuthUsecase interface {
	Login(ctx context.Context, in *in.LoginRequest) (*out.LoginResponse, error)
	Register(ctx context.Context, in *in.RegisterRequest) (*out.RegisterResponse, error)
	Logout(ctx context.Context, in *in.LogoutRequest) (*out.LogoutResponse, error)
}
