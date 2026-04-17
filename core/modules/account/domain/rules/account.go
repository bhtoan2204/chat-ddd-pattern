package rules

import (
	"errors"
	"strings"

	accounttypes "go-socket/core/modules/account/types"
)

var (
	ErrAccountIDRequired           = errors.New("account id is required")
	ErrAccountDisplayNameRequired  = errors.New("display_name is required")
	ErrAccountStatusInvalid        = errors.New("status is invalid")
	ErrAccountAlreadyVerified      = errors.New("email already verified")
	ErrAccountEmailMismatch        = errors.New("account email does not match")
	ErrAccountPasswordSameAsOldOne = errors.New("new password must be different from current password")
	ErrAccountAlreadyRegistered    = errors.New("account already registered")
	ErrAccountNotRegistered        = errors.New("account is not registered")
	ErrAccountNotFound             = errors.New("account not found")
)

func NormalizeAccountID(id string) (string, error) {
	normalized := strings.TrimSpace(id)
	if normalized == "" {
		return "", ErrAccountIDRequired
	}
	return normalized, nil
}

func NormalizeDisplayName(displayName string) (string, error) {
	normalized := strings.TrimSpace(displayName)
	if normalized == "" {
		return "", ErrAccountDisplayNameRequired
	}
	return normalized, nil
}

func NormalizeStatus(status accounttypes.AccountStatus) (accounttypes.AccountStatus, error) {
	parsed, err := accounttypes.ParseAccountStatus(status.String())
	if err != nil {
		return "", ErrAccountStatusInvalid
	}
	return parsed, nil
}

func NormalizeOptionalString(value string) *string {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func EqualOptionalString(left, right *string) bool {
	switch {
	case left == nil && right == nil:
		return true
	case left == nil || right == nil:
		return false
	default:
		return *left == *right
	}
}
