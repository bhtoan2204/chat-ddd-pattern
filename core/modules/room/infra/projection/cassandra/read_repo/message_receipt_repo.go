package read_repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	roomprojection "go-socket/core/modules/room/application/projection"
	"go-socket/core/modules/room/infra/projection/cassandra/views"
	"go-socket/core/shared/pkg/stackErr"
	"go-socket/core/shared/utils"

	"github.com/gocql/gocql"
)

type MessageReceiptRepo struct {
	session              *gocql.Session
	messageReceiptsTable string
}

func NewMessageReceiptRepo(session *gocql.Session, tables views.ProjectionTableNames) *MessageReceiptRepo {
	return &MessageReceiptRepo{
		session:              session,
		messageReceiptsTable: tables.MessageReceipts,
	}
}

func (r *MessageReceiptRepo) Upsert(ctx context.Context, projection *roomprojection.MessageReceiptProjection) error {
	if projection == nil {
		return nil
	}
	statement := fmt.Sprintf(`INSERT INTO %s (message_id,account_id,room_id,status,delivered_at,seen_at,created_at,updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, r.messageReceiptsTable)
	return stackErr.Error(r.session.Query(statement, projection.MessageID, projection.AccountID, projection.RoomID, projection.Status, projection.DeliveredAt, projection.SeenAt, projection.CreatedAt.UTC(), projection.UpdatedAt.UTC()).WithContext(ctx).Exec())
}

func (r *MessageReceiptRepo) GetMessageReceipt(ctx context.Context, lookup roomprojection.MessageReceiptLookup) (*roomprojection.MessageReceiptStatus, error) {
	statement := fmt.Sprintf(`SELECT status, delivered_at, seen_at FROM %s WHERE message_id = ? AND account_id = ?`, r.messageReceiptsTable)
	var (
		status      string
		deliveredAt *time.Time
		seenAt      *time.Time
	)
	if err := r.session.Query(statement, strings.TrimSpace(lookup.MessageID), strings.TrimSpace(lookup.AccountID)).WithContext(ctx).Scan(&status, &deliveredAt, &seenAt); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, nil
		}
		return nil, stackErr.Error(err)
	}
	return &roomprojection.MessageReceiptStatus{Status: status, DeliveredAt: utils.ClonePtr(deliveredAt), SeenAt: utils.ClonePtr(seenAt)}, nil
}

func (r *MessageReceiptRepo) CountByStatus(ctx context.Context, messageID, status string) (int64, error) {
	statement := fmt.Sprintf(`SELECT status FROM %s WHERE message_id = ?`, r.messageReceiptsTable)
	iter := r.session.Query(statement, strings.TrimSpace(messageID)).WithContext(ctx).Iter()
	defer iter.Close()
	var (
		value string
		count int64
	)
	scanner := iter.Scanner()
	for scanner.Next() {
		if err := scanner.Scan(&value); err != nil {
			return 0, stackErr.Error(fmt.Errorf("scan cassandra message receipt failed: %w", err))
		}
		if strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(status)) {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, stackErr.Error(fmt.Errorf("iterate cassandra message receipts failed: %w", err))
	}
	if err := iter.Close(); err != nil {
		return 0, stackErr.Error(fmt.Errorf("close cassandra message receipts iterator failed: %w", err))
	}
	return count, nil
}
