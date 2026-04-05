package types

import (
	"database/sql/driver"
	"strings"
)

type RoomRole string

const (
	RoomRoleOwner  RoomRole = "owner"
	RoomRoleAdmin  RoomRole = "admin"
	RoomRoleMember RoomRole = "member"
)

func (r RoomRole) Normalize() RoomRole {
	return RoomRole(strings.ToLower(strings.TrimSpace(string(r))))
}

func (r RoomRole) IsValid() bool {
	switch r.Normalize() {
	case RoomRoleOwner, RoomRoleAdmin, RoomRoleMember:
		return true
	default:
		return false
	}
}

func (r RoomRole) Value() (driver.Value, error) {
	return string(r), nil
}

func (r *RoomRole) Scan(value interface{}) error {
	if value == nil {
		*r = ""
		return nil
	}
	*r = RoomRole(value.(string))
	return nil
}
