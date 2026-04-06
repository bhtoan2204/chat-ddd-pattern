package service

import (
	"context"
	"errors"

	valueobject "go-socket/core/modules/account/domain/value_object"
)

var ErrAccountEmailAlreadyExists = errors.New("account email already exists")

type EmailUniquenessChecker interface {
	IsEmailExists(ctx context.Context, email string) (bool, error)
}

func EnsureEmailAvailable(ctx context.Context, checker EmailUniquenessChecker, email valueobject.Email) error {
	exists, err := checker.IsEmailExists(ctx, email.Value())
	if err != nil {
		return err
	}
	if exists {
		return ErrAccountEmailAlreadyExists
	}
	return nil
}
