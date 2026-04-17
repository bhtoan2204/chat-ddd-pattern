package provider

import (
	"context"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiryUnix   int64
	IDToken      string
}

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

//go:generate mockgen -package=provider -destination=auth_provider_mock.go -source=auth_provider.go
type AuthProvider interface {
	Login() string
	Callback(ctx context.Context, code string) (*AuthResult, error)
	UserInfo(ctx context.Context, accessToken string) (*UserInfo, error)
	Name() string
}
