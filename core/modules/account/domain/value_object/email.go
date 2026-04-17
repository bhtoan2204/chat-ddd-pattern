package valueobject

import (
	"encoding/json"
	"errors"
	"net/mail"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Email{}, errors.New("email is required")
	}

	addr, err := mail.ParseAddress(value)
	if err != nil {
		return Email{}, errors.New("invalid email format")
	}

	normalized := strings.ToLower(addr.Address)

	if len(normalized) > 254 {
		return Email{}, errors.New("email is too long")
	}

	return Email{value: normalized}, nil

}

func (e Email) Value() string {
	return e.value
}

func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.value)
}

func (e *Email) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return stackErr.Error(err)
	}

	email, err := NewEmail(value)
	if err != nil {
		return stackErr.Error(err)
	}

	*e = email
	return nil
}
