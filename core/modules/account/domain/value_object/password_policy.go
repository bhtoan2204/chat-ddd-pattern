package valueobject

import (
	"errors"
	"strings"
)

var ErrPasswordRequired = errors.New("password is required")

func normalizePasswordValue(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", ErrPasswordRequired
	}
	return value, nil
}
