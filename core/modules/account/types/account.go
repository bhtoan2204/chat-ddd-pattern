package types

import (
	"errors"
	"strings"
)

type AccountType string

const (
	AccountTypeUser  AccountType = "user"
	AccountTypeAdmin AccountType = "admin"
)

type AccountRole string

const (
	AccountRoleUser  AccountRole = "user"
	AccountRoleAdmin AccountRole = "admin"
)

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
)

func ParseAccountStatus(value string) (AccountStatus, error) {
	switch normalized := AccountStatus(strings.ToLower(strings.TrimSpace(value))); normalized {
	case AccountStatusActive, AccountStatusInactive:
		return normalized, nil
	default:
		return "", errors.New("status is invalid")
	}
}

func (s AccountStatus) String() string {
	return string(s)
}
