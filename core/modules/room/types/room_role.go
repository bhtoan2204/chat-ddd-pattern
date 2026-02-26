package types

import "database/sql/driver"

type RoomRole string

const (
	RoomRoleOwner  RoomRole = "owner"
	RoomRoleAdmin  RoomRole = "admin"
	RoomRoleMember RoomRole = "member"
)

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
