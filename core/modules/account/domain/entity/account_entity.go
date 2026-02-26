package entity

import (
	valueobject "go-socket/core/modules/account/domain/value_object"
	"time"
)

type Account struct {
	ID        string               `json:"id"`
	Email     valueobject.Email    `json:"email"`
	Password  valueobject.Password `json:"password"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
