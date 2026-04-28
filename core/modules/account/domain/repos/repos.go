package repos

import (
	"context"
)

//go:generate mockgen -package=repos -destination=repos_mock.go -source=repos.go
type Repos interface {
	AccountAggregateRepository() AccountAggregateRepository
	DeviceAggregateRepository() DeviceAggregateRepository
	SessionAggregateRepository() SessionAggregateRepository
	DeviceRepository() DeviceRepository
	SessionRepository() SessionRepository

	WithTransaction(ctx context.Context, fn func(Repos) error) error
}
