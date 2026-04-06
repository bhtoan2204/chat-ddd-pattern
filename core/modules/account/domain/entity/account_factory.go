package entity

import (
	"time"

	"go-socket/core/modules/account/domain/rules"
	valueobject "go-socket/core/modules/account/domain/value_object"
	accounttypes "go-socket/core/modules/account/types"
)

func NewAccount(
	id string,
	email valueobject.Email,
	passwordHash valueobject.HashedPassword,
	displayName string,
	status accounttypes.AccountStatus,
	now time.Time,
) (*Account, error) {
	normalizedID, err := rules.NormalizeAccountID(id)
	if err != nil {
		return nil, err
	}
	normalizedDisplayName, err := rules.NormalizeDisplayName(displayName)
	if err != nil {
		return nil, err
	}
	normalizedStatus, err := rules.NormalizeStatus(status)
	if err != nil {
		return nil, err
	}

	normalizedNow := rules.NormalizeAccountTime(now)
	return &Account{
		ID:           normalizedID,
		Email:        email,
		PasswordHash: passwordHash,
		DisplayName:  normalizedDisplayName,
		Status:       normalizedStatus,
		CreatedAt:    normalizedNow,
		UpdatedAt:    normalizedNow,
	}, nil
}
