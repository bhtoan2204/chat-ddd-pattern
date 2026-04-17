package repos

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicate      = errors.New("duplicate value")
	ErrAlreadyApplied = errors.New("ledger posting already applied")
)
