package repos

import "context"

type Repos interface {
	RoomRepository() RoomRepository
	MessageRepository() MessageRepository
	RoomMemberRepository() RoomMemberRepository
	RoomOutboxEventsRepository() RoomOutboxEventsRepository

	WithTransaction(ctx context.Context, fn func(Repos) error) error
}
