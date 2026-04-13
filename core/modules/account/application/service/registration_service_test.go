package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-socket/core/modules/account/domain/aggregate"
	"go-socket/core/modules/account/domain/entity"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/infra/xpaseto"
)

func TestRegistrationServiceRegister(t *testing.T) {
	t.Parallel()

	expiresAt := time.Date(2026, time.April, 14, 10, 0, 0, 0, time.UTC)
	fakeRepos := &stubAccountRepos{
		accountRepo: &stubAccountRepository{
			isEmailExists: func(context.Context, string) (bool, error) {
				return false, nil
			},
		},
		aggregateRepo: &stubAccountAggregateRepository{
			save: func(_ context.Context, agg *aggregate.AccountAggregate) error {
				if agg == nil || agg.AccountID == "" {
					t.Fatalf("expected persisted aggregate with generated account id")
				}
				return nil
			},
		},
	}

	service := &registrationService{
		baseRepo: fakeRepos,
		hasher: &stubHasher{
			hash: func(context.Context, string) (string, error) {
				return "hashed-password", nil
			},
		},
		paseto: &stubPasetoService{
			generateToken: func(_ context.Context, account *entity.Account) (string, time.Time, error) {
				if account == nil {
					t.Fatalf("expected generated account snapshot")
				}
				if account.Email.Value() != "alice@example.com" {
					t.Fatalf("unexpected email in token payload: %s", account.Email.Value())
				}
				return "signed-token", expiresAt, nil
			},
		},
	}

	result, err := service.Register(context.Background(), RegisterAccountCommand{
		Email:       "alice@example.com",
		Password:    "password123",
		DisplayName: "Alice",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if result.Token != "signed-token" {
		t.Fatalf("expected signed token, got %q", result.Token)
	}
	if !result.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("expected expiresAt %v, got %v", expiresAt, result.ExpiresAt)
	}
	if fakeRepos.transactionCalls != 1 {
		t.Fatalf("expected 1 transaction call, got %d", fakeRepos.transactionCalls)
	}
}

func TestRegistrationServiceRegisterReturnsAccountExists(t *testing.T) {
	t.Parallel()

	service := &registrationService{
		baseRepo: &stubAccountRepos{
			accountRepo: &stubAccountRepository{
				isEmailExists: func(context.Context, string) (bool, error) {
					return true, nil
				},
			},
			aggregateRepo: &stubAccountAggregateRepository{},
		},
		hasher: &stubHasher{},
		paseto: &stubPasetoService{},
	}

	_, err := service.Register(context.Background(), RegisterAccountCommand{
		Email:       "alice@example.com",
		Password:    "password123",
		DisplayName: "Alice",
	})
	if !errors.Is(err, ErrRegistrationAccountExists) {
		t.Fatalf("expected ErrRegistrationAccountExists, got %v", err)
	}
}

type stubAccountRepos struct {
	accountRepo      repos.AccountRepository
	aggregateRepo    repos.AccountAggregateRepository
	transactionCalls int
}

func (s *stubAccountRepos) AccountRepository() repos.AccountRepository {
	return s.accountRepo
}

func (s *stubAccountRepos) AccountAggregateRepository() repos.AccountAggregateRepository {
	return s.aggregateRepo
}

func (s *stubAccountRepos) WithTransaction(_ context.Context, fn func(repos.Repos) error) error {
	s.transactionCalls++
	return fn(s)
}

type stubAccountRepository struct {
	isEmailExists func(ctx context.Context, email string) (bool, error)
}

func (s *stubAccountRepository) GetAccountByID(context.Context, string) (*entity.Account, error) {
	return nil, nil
}

func (s *stubAccountRepository) GetAccountByEmail(context.Context, string) (*entity.Account, error) {
	return nil, nil
}

func (s *stubAccountRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	if s.isEmailExists == nil {
		return false, nil
	}
	return s.isEmailExists(ctx, email)
}

func (s *stubAccountRepository) CreateAccount(context.Context, *entity.Account) error {
	return nil
}

func (s *stubAccountRepository) UpdateAccount(context.Context, *entity.Account) error {
	return nil
}

func (s *stubAccountRepository) DeleteAccount(context.Context, string) error {
	return nil
}

func (s *stubAccountRepository) ListAccountsByRoomID(context.Context, string) ([]*entity.Account, error) {
	return nil, nil
}

func (s *stubAccountRepository) SearchUsers(context.Context, string, int, int) ([]*entity.Account, int64, error) {
	return nil, 0, nil
}

type stubAccountAggregateRepository struct {
	load        func(ctx context.Context, accountID string) (*aggregate.AccountAggregate, error)
	loadByEmail func(ctx context.Context, email string) (*aggregate.AccountAggregate, error)
	save        func(ctx context.Context, agg *aggregate.AccountAggregate) error
}

func (s *stubAccountAggregateRepository) Load(ctx context.Context, accountID string) (*aggregate.AccountAggregate, error) {
	if s.load == nil {
		return nil, nil
	}
	return s.load(ctx, accountID)
}

func (s *stubAccountAggregateRepository) LoadByEmail(ctx context.Context, email string) (*aggregate.AccountAggregate, error) {
	if s.loadByEmail == nil {
		return nil, nil
	}
	return s.loadByEmail(ctx, email)
}

func (s *stubAccountAggregateRepository) Save(ctx context.Context, agg *aggregate.AccountAggregate) error {
	if s.save == nil {
		return nil
	}
	return s.save(ctx, agg)
}

type stubHasher struct {
	hash   func(ctx context.Context, value string) (string, error)
	verify func(ctx context.Context, val string, hash string) (bool, error)
}

func (s *stubHasher) Hash(ctx context.Context, value string) (string, error) {
	if s.hash == nil {
		return "", nil
	}
	return s.hash(ctx, value)
}

func (s *stubHasher) Verify(ctx context.Context, val string, hash string) (bool, error) {
	if s.verify == nil {
		return false, nil
	}
	return s.verify(ctx, val, hash)
}

type stubPasetoService struct {
	generateToken func(ctx context.Context, account *entity.Account) (string, time.Time, error)
}

func (s *stubPasetoService) GenerateToken(ctx context.Context, account *entity.Account) (string, time.Time, error) {
	if s.generateToken == nil {
		return "", time.Time{}, nil
	}
	return s.generateToken(ctx, account)
}

func (s *stubPasetoService) ParseToken(context.Context, string) (*xpaseto.PasetoPayload, error) {
	return nil, nil
}
