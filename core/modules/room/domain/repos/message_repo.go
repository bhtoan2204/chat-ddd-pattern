package repos

import (
	"context"
	"go-socket/core/modules/room/domain/entity"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *entity.MessageEntity) error
}
