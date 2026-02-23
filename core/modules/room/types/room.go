package types

type RoomType string

const (
	RoomTypePublic  RoomType = "public"
	RoomTypePrivate RoomType = "private"
)

type RoomRole string

const (
	RoomRoleOwner  RoomRole = "owner"
	RoomRoleAdmin  RoomRole = "admin"
	RoomRoleMember RoomRole = "member"
)
