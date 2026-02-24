package command

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountExists      = errors.New("account already exists")
	ErrAccountNotFound    = errors.New("account not found")
	ErrCheckAccountFailed = errors.New("check account failed")
)
