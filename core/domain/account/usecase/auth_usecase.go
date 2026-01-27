package usecase

import (
	"context"
	"go-socket/core/delivery/http/data/in"
	"go-socket/core/delivery/http/data/out"
)

type AuthUsecase interface {
	Login(ctx context.Context, in *in.LoginRequest) (*out.LoginResponse, error)
	Register(ctx context.Context, in *in.RegisterRequest) (*out.RegisterResponse, error)
	Logout(ctx context.Context, in *in.LogoutRequest) (*out.LogoutResponse, error)
	GetProfile(ctx context.Context, in *in.GetProfileRequest) (*out.GetProfileResponse, error)
}
