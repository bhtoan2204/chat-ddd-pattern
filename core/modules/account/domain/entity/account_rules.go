package entity

import (
	"errors"
	"strings"
	"time"

	valueobject "go-socket/core/modules/account/domain/value_object"
)

var (
	ErrAccountIDRequired           = errors.New("account id is required")
	ErrAccountDisplayNameRequired  = errors.New("display_name is required")
	ErrAccountStatusInvalid        = errors.New("status is invalid")
	ErrAccountAlreadyVerified      = errors.New("email already verified")
	ErrAccountPasswordSameAsOldOne = errors.New("new password must be different from current password")
)

const (
	AccountStatusActive   = "active"
	AccountStatusInactive = "inactive"
)

func NewAccount(id string, email valueobject.Email, password valueobject.Password, displayName string, now time.Time) (*Account, error) {
	id = strings.TrimSpace(id)
	displayName = strings.TrimSpace(displayName)

	if id == "" {
		return nil, ErrAccountIDRequired
	}
	if displayName == "" {
		return nil, ErrAccountDisplayNameRequired
	}

	now = normalizeAccountTime(now)
	return &Account{
		ID:          id,
		Email:       email,
		Password:    password,
		DisplayName: displayName,
		Status:      AccountStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (a *Account) UpdateProfile(displayName string, username, avatarObjectKey *string, now time.Time) (bool, error) {
	if a == nil {
		return false, ErrAccountIDRequired
	}

	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		return false, ErrAccountDisplayNameRequired
	}

	updated := false
	if a.DisplayName != displayName {
		a.DisplayName = displayName
		updated = true
	}

	if username != nil {
		normalizedUsername := normalizeOptionalString(*username)
		if !equalOptionalString(a.Username, normalizedUsername) {
			a.Username = normalizedUsername
			updated = true
		}
	}

	if avatarObjectKey != nil {
		normalizedAvatarObjectKey := normalizeOptionalString(*avatarObjectKey)
		if !equalOptionalString(a.AvatarObjectKey, normalizedAvatarObjectKey) {
			a.AvatarObjectKey = normalizedAvatarObjectKey
			updated = true
		}
	}

	if updated {
		a.UpdatedAt = normalizeAccountTime(now)
	}
	return updated, nil
}

func (a *Account) MarkEmailVerified(now time.Time) (bool, error) {
	if a == nil {
		return false, ErrAccountIDRequired
	}
	if a.EmailVerifiedAt != nil {
		return false, ErrAccountAlreadyVerified
	}

	verifiedAt := normalizeAccountTime(now)
	a.EmailVerifiedAt = &verifiedAt
	a.UpdatedAt = verifiedAt
	return true, nil
}

func (a *Account) ChangePassword(password valueobject.Password, now time.Time) (bool, error) {
	if a == nil {
		return false, ErrAccountIDRequired
	}
	if a.Password.Value() == password.Value() {
		return false, ErrAccountPasswordSameAsOldOne
	}

	changedAt := normalizeAccountTime(now)
	a.Password = password
	a.PasswordChangedAt = &changedAt
	a.UpdatedAt = changedAt
	return true, nil
}

func normalizeOptionalString(value string) *string {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func equalOptionalString(left, right *string) bool {
	switch {
	case left == nil && right == nil:
		return true
	case left == nil || right == nil:
		return false
	default:
		return *left == *right
	}
}

func normalizeAccountTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now().UTC()
	}
	return value.UTC()
}
