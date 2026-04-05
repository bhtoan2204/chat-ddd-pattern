package aggregate

import "time"

type EventAccountCreated struct {
	AccountID   string
	Email       string
	DisplayName string
	Status      string
	CreatedAt   time.Time
}

type EventAccountUpdated struct {
	AccountID string
	Email     string
	UpdatedAt time.Time
}

type EventAccountProfileUpdated struct {
	AccountID       string
	DisplayName     string
	Username        *string
	AvatarObjectKey *string
	UpdatedAt       time.Time
}

type EventAccountEmailVerificationRequested struct {
	AccountID         string
	Email             string
	VerificationToken string
	RequestedAt       time.Time
}

type EventAccountEmailVerified struct {
	AccountID       string
	EmailVerifiedAt time.Time
}

type EventAccountPasswordChanged struct {
	AccountID         string
	PasswordChangedAt time.Time
}

type EventAccountBanned struct {
	AccountID string
	BanReason string
	BanUntil  *time.Time
}
