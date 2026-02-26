package valueobject

import (
	"encoding/json"
	"errors"
	"strings"
)

type Password struct {
	value string
}

func NewPassword(value string) (Password, error) {
	if strings.TrimSpace(value) == "" {
		return Password{}, errors.New("password is required")
	}
	return Password{value: value}, nil
}

func (p Password) Value() string {
	return p.value
}

func (p Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.value)
}

func (p *Password) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	password, err := NewPassword(value)
	if err != nil {
		return err
	}

	*p = password
	return nil
}
