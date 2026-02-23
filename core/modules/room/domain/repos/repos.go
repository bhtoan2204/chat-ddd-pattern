package repos

type Repos interface {
	RoomRepository() RoomRepository
	MessageRepository() MessageRepository
	RoomMemberRepository() RoomMemberRepository
}
