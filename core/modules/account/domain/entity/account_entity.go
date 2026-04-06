package entity

import (
	valueobject "go-socket/core/modules/account/domain/value_object"
	accounttypes "go-socket/core/modules/account/types"
	"time"
)

type Account struct {
	ID                string                     `json:"id"`
	Email             valueobject.Email          `json:"email"`
	PasswordHash      valueobject.HashedPassword `json:"password_hash"`
	DisplayName       string                     `json:"display_name"`
	Username          *string                    `json:"username,omitempty"`
	AvatarObjectKey   *string                    `json:"avatar_object_key,omitempty"`
	Status            accounttypes.AccountStatus `json:"status"`
	EmailVerifiedAt   *time.Time                 `json:"email_verified_at,omitempty"`
	LastLoginAt       *time.Time                 `json:"last_login_at,omitempty"`
	PasswordChangedAt *time.Time                 `json:"password_changed_at,omitempty"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
	BannedReason      string                     `json:"banned_reason"`
	BannedUntil       *time.Time                 `json:"banned_until"`
}
