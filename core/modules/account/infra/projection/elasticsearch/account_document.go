package elasticsearch

import "time"

type accountDocument struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	DisplayName       string     `json:"display_name"`
	Username          *string    `json:"username,omitempty"`
	AvatarObjectKey   *string    `json:"avatar_object_key,omitempty"`
	Status            string     `json:"status"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	BannedReason      string     `json:"banned_reason,omitempty"`
	BannedUntil       *time.Time `json:"banned_until,omitempty"`
}
