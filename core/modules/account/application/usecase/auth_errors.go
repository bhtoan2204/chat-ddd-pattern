package usecase

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountExists      = errors.New("account already exists")
	ErrAccountNotFound    = errors.New("account not found")
)
