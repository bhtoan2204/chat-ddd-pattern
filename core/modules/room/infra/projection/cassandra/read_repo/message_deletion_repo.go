package read_repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	roomprojection "go-socket/core/modules/room/application/projection"
	"go-socket/core/modules/room/infra/projection/cassandra/views"
	"go-socket/core/shared/pkg/stackErr"

	"github.com/gocql/gocql"
)

type MessageDeletionRepo struct {
	session             *gocql.Session
	messageDeletesTable string
}

func NewMessageDeletionRepo(session *gocql.Session, tables views.ProjectionTableNames) *MessageDeletionRepo {
	return &MessageDeletionRepo{
		session:             session,
		messageDeletesTable: tables.MessageDeletions,
	}
}

func (r *MessageDeletionRepo) Upsert(ctx context.Context, projection *roomprojection.MessageDeletionProjection) error {
	if projection == nil {
		return nil
	}
	statement := fmt.Sprintf(`INSERT INTO %s (account_id,room_id,message_sent_at,message_id,created_at) VALUES (?, ?, ?, ?, ?)`, r.messageDeletesTable)
	return stackErr.Error(r.session.Query(statement, projection.AccountID, projection.RoomID, projection.MessageSentAt.UTC(), projection.MessageID, projection.CreatedAt.UTC()).WithContext(ctx).Exec())
}

func (r *MessageDeletionRepo) DeletePartition(ctx context.Context, accountID, roomID string) error {
	statement := fmt.Sprintf(`DELETE FROM %s WHERE account_id = ? AND room_id = ?`, r.messageDeletesTable)
	return stackErr.Error(r.session.Query(statement, strings.TrimSpace(accountID), strings.TrimSpace(roomID)).WithContext(ctx).Exec())
}

func (r *MessageDeletionRepo) ListDeletedMessageIDs(ctx context.Context, accountID, roomID string, from, to *time.Time) (map[string]struct{}, error) {
	if accountID == "" || roomID == "" || from == nil || to == nil {
		return map[string]struct{}{}, nil
	}
	statement := fmt.Sprintf(`SELECT message_id FROM %s WHERE account_id = ? AND room_id = ? AND message_sent_at >= ? AND message_sent_at <= ?`, r.messageDeletesTable)
	iter := r.session.Query(statement, accountID, roomID, from.UTC(), to.UTC()).WithContext(ctx).Iter()
	defer iter.Close()
	results := make(map[string]struct{})
	scanner := iter.Scanner()
	var messageID string
	for scanner.Next() {
		if err := scanner.Scan(&messageID); err != nil {
			return nil, stackErr.Error(fmt.Errorf("scan cassandra message deletion projection failed: %v", err))
		}
		results[strings.TrimSpace(messageID)] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, stackErr.Error(fmt.Errorf("iterate cassandra message deletions failed: %v", err))
	}
	if err := iter.Close(); err != nil {
		return nil, stackErr.Error(fmt.Errorf("close cassandra message deletion iterator failed: %v", err))
	}
	return results, nil
}
