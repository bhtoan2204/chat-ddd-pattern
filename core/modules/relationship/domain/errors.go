package domain

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrDuplicate   = errors.New("duplicate value")
	ErrEmpty       = errors.New("empty value")
	ErrInvalidData = errors.New("invalid data")
)
