package service

import (
	"context"

	"wechat-clone/core/modules/account/domain/rules"
	valueobject "wechat-clone/core/modules/account/domain/value_object"
	"wechat-clone/core/shared/pkg/stackErr"
)

//go:generate mockgen -package=service -destination=password_reuse_mock.go -source=password_reuse.go
type PasswordReuseChecker interface {
	Verify(ctx context.Context, val string, hash string) (bool, error)
}

func EnsurePasswordIsNew(
	ctx context.Context,
	checker PasswordReuseChecker,
	newPassword valueobject.PlainPassword,
	currentHash valueobject.HashedPassword,
) error {
	isSamePassword, err := checker.Verify(ctx, newPassword.Value(), currentHash.Value())
	if err != nil {
		return stackErr.Error(err)
	}
	if isSamePassword {
		return stackErr.Error(rules.ErrAccountPasswordSameAsOldOne)
	}
	return nil
}
