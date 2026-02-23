package domainevent

import "time"

type UserCreatedEvent struct {
	UserID    string
	Email     string
	CreatedAt time.Time
}
