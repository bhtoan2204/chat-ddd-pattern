package aggregate

import (
	"time"
	accounttypes "wechat-clone/core/modules/account/types"
)

type EventAccountCreated struct {
	AccountID    string
	Email        string
	PasswordHash string
	DisplayName  string
	Status       accounttypes.AccountStatus
	CreatedAt    time.Time
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
	PasswordHash      string
	PasswordChangedAt time.Time
}

type EventAccountBanned struct {
	AccountID string
	BanReason string
	BanUntil  *time.Time
}
