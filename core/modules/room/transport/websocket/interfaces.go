package socket

import "context"

type IClient interface {
	GetID() string
	GetUserID() string
	Send(message []byte)
	ReadPump(ctx context.Context, hub IHub)
	WritePump(ctx context.Context)
	Close()
}

type IRoom interface {
	GetID() string
	AddClient(client IClient)
	RemoveClient(client IClient)
	Broadcast(message []byte)
	IsEmpty() bool
	ClientCount() int
}

type IHub interface {
	Register(client IClient)
	Unregister(client IClient)

	JoinRoom(ctx context.Context, client IClient, roomID string) error
	LeaveRoom(ctx context.Context, client IClient, roomID string) error

	HandleMessage(ctx context.Context, client IClient, msg Message) error
	Close()
}
