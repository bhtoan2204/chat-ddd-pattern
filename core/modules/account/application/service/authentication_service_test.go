package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-socket/core/modules/account/domain/aggregate"
	"go-socket/core/modules/account/domain/entity"
	valueobject "go-socket/core/modules/account/domain/value_object"

	"gorm.io/gorm"
)

func TestAuthenticationServiceAuthenticate(t *testing.T) {
	t.Parallel()

	email, err := valueobject.NewEmail("alice@example.com")
	if err != nil {
		t.Fatalf("NewEmail() error = %v", err)
	}
	passwordHash, err := valueobject.NewHashedPassword("hashed-password")
	if err != nil {
		t.Fatalf("NewHashedPassword() error = %v", err)
	}

	accountAggregate, err := aggregate.NewAccountAggregate("acc-1")
	if err != nil {
		t.Fatalf("NewAccountAggregate() error = %v", err)
	}
	if err := accountAggregate.Register(email, passwordHash, "Alice", time.Date(2026, time.April, 14, 8, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	expiresAt := time.Date(2026, time.April, 14, 11, 0, 0, 0, time.UTC)
	service := &authenticationService{
		baseRepo: &stubAccountRepos{
			aggregateRepo: &stubAccountAggregateRepository{
				loadByEmail: func(_ context.Context, email string) (*aggregate.AccountAggregate, error) {
					if email != "alice@example.com" {
						t.Fatalf("unexpected email lookup: %s", email)
					}
					return accountAggregate, nil
				},
			},
		},
		hasher: &stubHasher{
			verify: func(_ context.Context, raw string, hash string) (bool, error) {
				if raw != "password123" {
					t.Fatalf("unexpected raw password: %s", raw)
				}
				if hash != "hashed-password" {
					t.Fatalf("unexpected hash: %s", hash)
				}
				return true, nil
			},
		},
		paseto: &stubPasetoService{
			generateToken: func(_ context.Context, account *entity.Account) (string, time.Time, error) {
				if account == nil || account.ID != "acc-1" {
					t.Fatalf("expected account snapshot from aggregate")
				}
				return "signed-token", expiresAt, nil
			},
		},
	}

	result, err := service.Authenticate(context.Background(), AuthenticateAccountCommand{
		Email:    "alice@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
	if result.Token != "signed-token" {
		t.Fatalf("expected signed-token, got %q", result.Token)
	}
	if !result.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("expected expiresAt %v, got %v", expiresAt, result.ExpiresAt)
	}
}

func TestAuthenticationServiceAuthenticateMapsNotFound(t *testing.T) {
	t.Parallel()

	service := &authenticationService{
		baseRepo: &stubAccountRepos{
			aggregateRepo: &stubAccountAggregateRepository{
				loadByEmail: func(context.Context, string) (*aggregate.AccountAggregate, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
		},
		hasher: &stubHasher{},
		paseto: &stubPasetoService{},
	}

	_, err := service.Authenticate(context.Background(), AuthenticateAccountCommand{
		Email:    "missing@example.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrAuthenticationAccountNotFound) {
		t.Fatalf("expected ErrAuthenticationAccountNotFound, got %v", err)
	}
}

func TestAuthenticationServiceAuthenticateMapsInvalidPassword(t *testing.T) {
	t.Parallel()

	email, _ := valueobject.NewEmail("alice@example.com")
	passwordHash, _ := valueobject.NewHashedPassword("hashed-password")
	accountAggregate, _ := aggregate.NewAccountAggregate("acc-1")
	if err := accountAggregate.Register(email, passwordHash, "Alice", time.Date(2026, time.April, 14, 8, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	service := &authenticationService{
		baseRepo: &stubAccountRepos{
			aggregateRepo: &stubAccountAggregateRepository{
				loadByEmail: func(context.Context, string) (*aggregate.AccountAggregate, error) {
					return accountAggregate, nil
				},
			},
		},
		hasher: &stubHasher{
			verify: func(context.Context, string, string) (bool, error) {
				return false, nil
			},
		},
		paseto: &stubPasetoService{},
	}

	_, err := service.Authenticate(context.Background(), AuthenticateAccountCommand{
		Email:    "alice@example.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrAuthenticationInvalidPassword) {
		t.Fatalf("expected ErrAuthenticationInvalidPassword, got %v", err)
	}
}
