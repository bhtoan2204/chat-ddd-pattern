package repos

import (
	"context"
	"go-socket/core/shared/pkg/event"
)

type AccountOutboxEventsRepository interface {
	Append(ctx context.Context, event event.Event) error
}
