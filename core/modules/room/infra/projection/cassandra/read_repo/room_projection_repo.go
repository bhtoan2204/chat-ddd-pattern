package read_repo

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-socket/core/modules/room/infra/projection/cassandra/views"
	"go-socket/core/shared/pkg/stackErr"
	"go-socket/core/shared/utils"

	"github.com/gocql/gocql"
)

type RoomProjectionRow struct {
	RoomID              string
	Name                string
	Description         string
	RoomType            string
	OwnerID             string
	PinnedMessageID     string
	MemberCount         int
	LastMessageID       string
	LastMessageAt       *time.Time
	LastMessageContent  string
	LastMessageSenderID string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type RoomMemberProjectionRow struct {
	RoomID          string
	MemberID        string
	AccountID       string
	DisplayName     string
	Username        string
	AvatarObjectKey string
	Role            string
	LastDeliveredAt *time.Time
	LastReadAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type RoomProjectionRepo struct {
	session             *gocql.Session
	roomTable           string
	roomsByAccountTable string
	roomMembersTable    string
}

func NewRoomProjectionRepo(session *gocql.Session, tables views.ProjectionTableNames) *RoomProjectionRepo {
	return &RoomProjectionRepo{
		session:             session,
		roomTable:           tables.RoomByID,
		roomsByAccountTable: tables.RoomByAccount,
		roomMembersTable:    tables.RoomMemberByRoom,
	}
}

func (r *RoomProjectionRepo) GetRoomRow(ctx context.Context, roomID string) (*RoomProjectionRow, error) {
	statement := fmt.Sprintf(`
		SELECT
			room_id,
			name,
			description,
			room_type,
			owner_id,
			pinned_message_id,
			member_count,
			last_message_id,
			last_message_at,
			last_message_content,
			last_message_sender_id,
			created_at,
			updated_at
		FROM %s
		WHERE room_id = ?
	`, r.roomTable)

	row := &RoomProjectionRow{}
	if err := r.session.Query(statement, strings.TrimSpace(roomID)).WithContext(ctx).Scan(
		&row.RoomID,
		&row.Name,
		&row.Description,
		&row.RoomType,
		&row.OwnerID,
		&row.PinnedMessageID,
		&row.MemberCount,
		&row.LastMessageID,
		&row.LastMessageAt,
		&row.LastMessageContent,
		&row.LastMessageSenderID,
		&row.CreatedAt,
		&row.UpdatedAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, nil
		}
		return nil, stackErr.Error(fmt.Errorf("get cassandra room projection failed: %v", err))
	}

	row.CreatedAt = row.CreatedAt.UTC()
	row.UpdatedAt = row.UpdatedAt.UTC()
	row.LastMessageAt = cloneTime(row.LastMessageAt)
	return row, nil
}

func (r *RoomProjectionRepo) UpsertRoomRow(ctx context.Context, row *RoomProjectionRow) error {
	statement := fmt.Sprintf(`
		INSERT INTO %s (
			room_id,
			name,
			description,
			room_type,
			owner_id,
			pinned_message_id,
			member_count,
			last_message_id,
			last_message_at,
			last_message_content,
			last_message_sender_id,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.roomTable)

	if err := r.session.Query(
		statement,
		row.RoomID,
		row.Name,
		row.Description,
		row.RoomType,
		row.OwnerID,
		nullableProjectionString(row.PinnedMessageID),
		row.MemberCount,
		nullableProjectionString(row.LastMessageID),
		row.LastMessageAt,
		nullableProjectionString(row.LastMessageContent),
		nullableProjectionString(row.LastMessageSenderID),
		row.CreatedAt.UTC(),
		row.UpdatedAt.UTC(),
	).WithContext(ctx).Exec(); err != nil {
		return stackErr.Error(fmt.Errorf("upsert cassandra room projection failed: %v", err))
	}
	return nil
}

func (r *RoomProjectionRepo) DeleteRoomRow(ctx context.Context, roomID string) error {
	if err := r.session.Query(
		fmt.Sprintf(`DELETE FROM %s WHERE room_id = ?`, r.roomTable),
		strings.TrimSpace(roomID),
	).WithContext(ctx).Exec(); err != nil {
		return stackErr.Error(fmt.Errorf("delete cassandra room projection failed: %v", err))
	}
	return nil
}

func (r *RoomProjectionRepo) ListRoomsByAccount(ctx context.Context, accountID string, limit, offset int) ([]*RoomProjectionRow, error) {
	queryLimit := limitWithOffset(limit, offset)
	statement := fmt.Sprintf(`
		SELECT
			room_id,
			name,
			description,
			room_type,
			owner_id,
			pinned_message_id,
			member_count,
			last_message_id,
			last_message_at,
			last_message_content,
			last_message_sender_id,
			created_at,
			room_updated_at
		FROM %s
		WHERE account_id = ?
		LIMIT ?
	`, r.roomsByAccountTable)

	rows := make([]*RoomProjectionRow, 0, queryLimit)
	iter := r.session.Query(statement, strings.TrimSpace(accountID), queryLimit).WithContext(ctx).Iter()
	defer iter.Close()

	var (
		roomID              string
		name                string
		description         string
		roomType            string
		ownerID             string
		pinnedMessageID     string
		memberCount         int
		lastMessageID       string
		lastMessageAt       *time.Time
		lastMessageContent  string
		lastMessageSenderID string
		createdAt           time.Time
		updatedAt           time.Time
	)
	scanner := iter.Scanner()
	for scanner.Next() {
		if err := scanner.Scan(
			&roomID,
			&name,
			&description,
			&roomType,
			&ownerID,
			&pinnedMessageID,
			&memberCount,
			&lastMessageID,
			&lastMessageAt,
			&lastMessageContent,
			&lastMessageSenderID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, stackErr.Error(fmt.Errorf("scan cassandra account room projection failed: %v", err))
		}
		rows = append(rows, &RoomProjectionRow{
			RoomID:              roomID,
			Name:                name,
			Description:         description,
			RoomType:            roomType,
			OwnerID:             ownerID,
			PinnedMessageID:     pinnedMessageID,
			MemberCount:         memberCount,
			LastMessageID:       lastMessageID,
			LastMessageAt:       cloneTime(lastMessageAt),
			LastMessageContent:  lastMessageContent,
			LastMessageSenderID: lastMessageSenderID,
			CreatedAt:           createdAt.UTC(),
			UpdatedAt:           updatedAt.UTC(),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, stackErr.Error(fmt.Errorf("iterate cassandra account room projections failed: %v", err))
	}
	if err := iter.Close(); err != nil {
		return nil, stackErr.Error(fmt.Errorf("close cassandra account room projection iterator failed: %v", err))
	}

	return sliceRoomRows(rows, offset, limit), nil
}

func (r *RoomProjectionRepo) ListRoomsFromBaseProjection(ctx context.Context, limit, offset int) ([]*RoomProjectionRow, error) {
	statement := fmt.Sprintf(`
		SELECT
			room_id,
			name,
			description,
			room_type,
			owner_id,
			pinned_message_id,
			member_count,
			last_message_id,
			last_message_at,
			last_message_content,
			last_message_sender_id,
			created_at,
			updated_at
		FROM %s
	`, r.roomTable)

	rows := make([]*RoomProjectionRow, 0)
	iter := r.session.Query(statement).WithContext(ctx).Iter()
	defer iter.Close()

	var (
		roomID              string
		name                string
		description         string
		roomType            string
		ownerID             string
		pinnedMessageID     string
		memberCount         int
		lastMessageID       string
		lastMessageAt       *time.Time
		lastMessageContent  string
		lastMessageSenderID string
		createdAt           time.Time
		updatedAt           time.Time
	)
	scanner := iter.Scanner()
	for scanner.Next() {
		if err := scanner.Scan(
			&roomID,
			&name,
			&description,
			&roomType,
			&ownerID,
			&pinnedMessageID,
			&memberCount,
			&lastMessageID,
			&lastMessageAt,
			&lastMessageContent,
			&lastMessageSenderID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, stackErr.Error(fmt.Errorf("scan cassandra room projection failed: %v", err))
		}
		rows = append(rows, &RoomProjectionRow{
			RoomID:              roomID,
			Name:                name,
			Description:         description,
			RoomType:            roomType,
			OwnerID:             ownerID,
			PinnedMessageID:     pinnedMessageID,
			MemberCount:         memberCount,
			LastMessageID:       lastMessageID,
			LastMessageAt:       cloneTime(lastMessageAt),
			LastMessageContent:  lastMessageContent,
			LastMessageSenderID: lastMessageSenderID,
			CreatedAt:           createdAt.UTC(),
			UpdatedAt:           updatedAt.UTC(),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, stackErr.Error(fmt.Errorf("iterate cassandra room projections failed: %v", err))
	}
	if err := iter.Close(); err != nil {
		return nil, stackErr.Error(fmt.Errorf("close cassandra room projection iterator failed: %v", err))
	}

	sort.Slice(rows, func(i, j int) bool {
		if !rows[i].UpdatedAt.Equal(rows[j].UpdatedAt) {
			return rows[i].UpdatedAt.After(rows[j].UpdatedAt)
		}
		return rows[i].RoomID > rows[j].RoomID
	})

	return sliceRoomRows(rows, offset, limit), nil
}

func (r *RoomProjectionRepo) ListRoomMemberRows(ctx context.Context, roomID string) ([]*RoomMemberProjectionRow, error) {
	statement := fmt.Sprintf(`
		SELECT
			room_id,
			account_id,
			member_id,
			display_name,
			username,
			avatar_object_key,
			role,
			last_delivered_at,
			last_read_at,
			created_at,
			updated_at
		FROM %s
		WHERE room_id = ?
	`, r.roomMembersTable)

	rows := make([]*RoomMemberProjectionRow, 0)
	iter := r.session.Query(statement, strings.TrimSpace(roomID)).WithContext(ctx).Iter()
	defer iter.Close()

	var (
		lastDeliveredAt *time.Time
		lastReadAt      *time.Time
	)
	scanner := iter.Scanner()
	for scanner.Next() {
		row := &RoomMemberProjectionRow{}
		lastDeliveredAt = nil
		lastReadAt = nil
		if err := scanner.Scan(
			&row.RoomID,
			&row.AccountID,
			&row.MemberID,
			&row.DisplayName,
			&row.Username,
			&row.AvatarObjectKey,
			&row.Role,
			&lastDeliveredAt,
			&lastReadAt,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, stackErr.Error(fmt.Errorf("scan cassandra room member projection failed: %v", err))
		}
		row.LastDeliveredAt = cloneTime(lastDeliveredAt)
		row.LastReadAt = cloneTime(lastReadAt)
		row.CreatedAt = row.CreatedAt.UTC()
		row.UpdatedAt = row.UpdatedAt.UTC()
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, stackErr.Error(fmt.Errorf("iterate cassandra room member projections failed: %v", err))
	}
	if err := iter.Close(); err != nil {
		return nil, stackErr.Error(fmt.Errorf("close cassandra room member projection iterator failed: %v", err))
	}
	return rows, nil
}

func (r *RoomProjectionRepo) UpsertRoomMemberRow(ctx context.Context, row *RoomMemberProjectionRow) error {
	statement := fmt.Sprintf(`
		INSERT INTO %s (
			room_id,
			account_id,
			member_id,
			display_name,
			username,
			avatar_object_key,
			role,
			last_delivered_at,
			last_read_at,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.roomMembersTable)

	if err := r.session.Query(
		statement,
		row.RoomID,
		row.AccountID,
		row.MemberID,
		nullableProjectionString(row.DisplayName),
		nullableProjectionString(row.Username),
		nullableProjectionString(row.AvatarObjectKey),
		row.Role,
		row.LastDeliveredAt,
		row.LastReadAt,
		row.CreatedAt.UTC(),
		row.UpdatedAt.UTC(),
	).WithContext(ctx).Exec(); err != nil {
		return stackErr.Error(fmt.Errorf("upsert cassandra room member projection failed: %v", err))
	}
	return nil
}

func (r *RoomProjectionRepo) DeleteRoomMemberRow(ctx context.Context, roomID, accountID string) error {
	if err := r.session.Query(
		fmt.Sprintf(`DELETE FROM %s WHERE room_id = ? AND account_id = ?`, r.roomMembersTable),
		strings.TrimSpace(roomID),
		strings.TrimSpace(accountID),
	).WithContext(ctx).Exec(); err != nil {
		return stackErr.Error(fmt.Errorf("delete cassandra room member projection failed: %v", err))
	}
	return nil
}

func (r *RoomProjectionRepo) GetRoomMemberByAccount(ctx context.Context, roomID, accountID string) (*RoomMemberProjectionRow, error) {
	statement := fmt.Sprintf(`
		SELECT
			room_id,
			account_id,
			member_id,
			display_name,
			username,
			avatar_object_key,
			role,
			last_delivered_at,
			last_read_at,
			created_at,
			updated_at
		FROM %s
		WHERE room_id = ? AND account_id = ?
	`, r.roomMembersTable)

	row := &RoomMemberProjectionRow{}
	if err := r.session.Query(statement, strings.TrimSpace(roomID), strings.TrimSpace(accountID)).WithContext(ctx).Scan(
		&row.RoomID,
		&row.AccountID,
		&row.MemberID,
		&row.DisplayName,
		&row.Username,
		&row.AvatarObjectKey,
		&row.Role,
		&row.LastDeliveredAt,
		&row.LastReadAt,
		&row.CreatedAt,
		&row.UpdatedAt,
	); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, nil
		}
		return nil, stackErr.Error(err)
	}

	row.LastDeliveredAt = cloneTime(row.LastDeliveredAt)
	row.LastReadAt = cloneTime(row.LastReadAt)
	row.CreatedAt = row.CreatedAt.UTC()
	row.UpdatedAt = row.UpdatedAt.UTC()
	return row, nil
}

func (r *RoomProjectionRepo) UpsertAccountRoomIndex(ctx context.Context, accountID string, room *RoomProjectionRow) error {
	statement := fmt.Sprintf(`
		INSERT INTO %s (
			account_id,
			room_updated_at,
			room_id,
			name,
			description,
			room_type,
			owner_id,
			pinned_message_id,
			member_count,
			last_message_id,
			last_message_at,
			last_message_content,
			last_message_sender_id,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.roomsByAccountTable)

	if err := r.session.Query(
		statement,
		accountID,
		room.UpdatedAt.UTC(),
		room.RoomID,
		room.Name,
		room.Description,
		room.RoomType,
		room.OwnerID,
		nullableProjectionString(room.PinnedMessageID),
		room.MemberCount,
		nullableProjectionString(room.LastMessageID),
		room.LastMessageAt,
		nullableProjectionString(room.LastMessageContent),
		nullableProjectionString(room.LastMessageSenderID),
		room.CreatedAt.UTC(),
	).WithContext(ctx).Exec(); err != nil {
		return stackErr.Error(fmt.Errorf("upsert cassandra room-by-account projection failed: %v", err))
	}
	return nil
}

func (r *RoomProjectionRepo) DeleteAccountRoomIndex(ctx context.Context, accountID string, room *RoomProjectionRow) error {
	if room == nil {
		return nil
	}
	if err := r.session.Query(
		fmt.Sprintf(`DELETE FROM %s WHERE account_id = ? AND room_updated_at = ? AND room_id = ?`, r.roomsByAccountTable),
		accountID,
		room.UpdatedAt.UTC(),
		room.RoomID,
	).WithContext(ctx).Exec(); err != nil {
		return stackErr.Error(fmt.Errorf("delete cassandra room-by-account projection failed: %v", err))
	}
	return nil
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copy := value.UTC()
	return &copy
}

func nullableProjectionString(value string) interface{} {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}

func sliceRoomRows(rows []*RoomProjectionRow, offset, limit int) []*RoomProjectionRow {
	if offset >= len(rows) {
		return []*RoomProjectionRow{}
	}
	end := offset + limit
	if end > len(rows) {
		end = len(rows)
	}
	return rows[offset:end]
}

func limitWithOffset(limit, offset int) int {
	value := limit + offset
	if value <= 0 {
		return limit
	}
	return value
}

func RoomRowToEntity(row *RoomProjectionRow) *views.RoomView {
	if row == nil {
		return nil
	}

	pinnedMessageID := utils.StringPtr(row.PinnedMessageID)
	lastMessageID := utils.StringPtr(row.LastMessageID)
	lastMessageContent := utils.StringPtr(row.LastMessageContent)
	lastMessageSenderID := utils.StringPtr(row.LastMessageSenderID)

	return &views.RoomView{
		ID:                  row.RoomID,
		Name:                row.Name,
		Description:         row.Description,
		RoomType:            row.RoomType,
		OwnerID:             row.OwnerID,
		PinnedMessageID:     pinnedMessageID,
		MemberCount:         row.MemberCount,
		LastMessageID:       lastMessageID,
		LastMessageAt:       cloneTime(row.LastMessageAt),
		LastMessageContent:  lastMessageContent,
		LastMessageSenderID: lastMessageSenderID,
		CreatedAt:           row.CreatedAt.UTC(),
		UpdatedAt:           row.UpdatedAt.UTC(),
	}
}

func RoomMemberRowToEntity(row *RoomMemberProjectionRow) *views.RoomMemberView {
	if row == nil {
		return nil
	}
	return &views.RoomMemberView{
		ID:              row.MemberID,
		RoomID:          row.RoomID,
		AccountID:       row.AccountID,
		DisplayName:     strings.TrimSpace(row.DisplayName),
		Username:        strings.TrimSpace(row.Username),
		AvatarObjectKey: strings.TrimSpace(row.AvatarObjectKey),
		Role:            row.Role,
		LastDeliveredAt: cloneTime(row.LastDeliveredAt),
		LastReadAt:      cloneTime(row.LastReadAt),
		CreatedAt:       row.CreatedAt.UTC(),
		UpdatedAt:       row.UpdatedAt.UTC(),
	}
}
