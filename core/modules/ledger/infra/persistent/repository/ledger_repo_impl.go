package repository

import (
	"context"
	"strings"
	"time"

	"wechat-clone/core/modules/ledger/domain/entity"
	ledgerrepos "wechat-clone/core/modules/ledger/domain/repos"
	"wechat-clone/core/modules/ledger/infra/persistent/model"

	"gorm.io/gorm"
)

type ledgerRepoImpl struct {
	db *gorm.DB
}

type ledgerTransactionListRow struct {
	TransactionID string
	Currency      string
	CreatedAt     time.Time
}

func NewLedgerRepoImpl(db *gorm.DB) ledgerrepos.LedgerRepository {
	return &ledgerRepoImpl{db: db}
}

func (r *ledgerRepoImpl) GetBalance(ctx context.Context, accountID, currency string) (int64, error) {
	var balance int64
	err := r.db.WithContext(ctx).
		Model(&model.LedgerEntryModel{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("account_id = ? AND currency = ?", accountID, currency).
		Scan(&balance).Error
	return balance, mapError(err)
}

func (r *ledgerRepoImpl) CountTransactions(ctx context.Context, accountID, currency string) (int64, error) {
	accountID = strings.TrimSpace(accountID)
	currency = strings.ToUpper(strings.TrimSpace(currency))

	query := r.listTransactionsBaseQuery(ctx, accountID, currency)

	var total int64
	err := query.Distinct("t.transaction_id").Count(&total).Error
	return total, mapError(err)
}

func (r *ledgerRepoImpl) GetTransaction(ctx context.Context, transactionID string) (*entity.LedgerTransaction, error) {
	var transactionModel model.LedgerTransactionModel
	if err := r.db.WithContext(ctx).
		Where("transaction_id = ?", transactionID).
		First(&transactionModel).Error; err != nil {
		return nil, mapError(err)
	}

	var entryModels []model.LedgerEntryModel
	if err := r.db.WithContext(ctx).
		Where("transaction_id = ?", transactionID).
		Order("id ASC").
		Find(&entryModels).Error; err != nil {
		return nil, mapError(err)
	}

	entries := make([]*entity.LedgerEntry, 0, len(entryModels))
	for _, entryModel := range entryModels {
		entry := entryModel
		entries = append(entries, &entity.LedgerEntry{
			ID:            entry.ID,
			TransactionID: entry.TransactionID,
			AccountID:     entry.AccountID,
			Currency:      entry.Currency,
			Amount:        entry.Amount,
			CreatedAt:     entry.CreatedAt,
		})
	}

	return &entity.LedgerTransaction{
		TransactionID: transactionModel.TransactionID,
		Currency:      transactionModel.Currency,
		CreatedAt:     transactionModel.CreatedAt,
		Entries:       entries,
	}, nil
}

func (r *ledgerRepoImpl) ListTransactions(ctx context.Context, filter ledgerrepos.ListTransactionsFilter) ([]*entity.LedgerTransaction, error) {
	query := r.listTransactionsBaseQuery(ctx, filter.AccountID, filter.Currency)

	if filter.CursorCreatedAt != nil && strings.TrimSpace(filter.CursorTransactionID) != "" {
		query = query.Where(
			"(t.created_at < ? OR (t.created_at = ? AND t.transaction_id < ?))",
			filter.CursorCreatedAt.UTC(),
			filter.CursorCreatedAt.UTC(),
			strings.TrimSpace(filter.CursorTransactionID),
		)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	var transactionRows []ledgerTransactionListRow
	if err := query.
		Select("t.transaction_id, t.currency, t.created_at").
		Group("t.transaction_id, t.currency, t.created_at").
		Order("t.created_at DESC").
		Order("t.transaction_id DESC").
		Find(&transactionRows).Error; err != nil {
		return nil, mapError(err)
	}
	if len(transactionRows) == 0 {
		return []*entity.LedgerTransaction{}, nil
	}

	transactionIDs := make([]string, 0, len(transactionRows))
	transactionsByID := make(map[string]*entity.LedgerTransaction, len(transactionRows))
	for _, row := range transactionRows {
		transactionIDs = append(transactionIDs, row.TransactionID)
		transactionsByID[row.TransactionID] = &entity.LedgerTransaction{
			TransactionID: row.TransactionID,
			Currency:      row.Currency,
			CreatedAt:     row.CreatedAt,
			Entries:       make([]*entity.LedgerEntry, 0),
		}
	}

	var entryModels []model.LedgerEntryModel
	if err := r.db.WithContext(ctx).
		Where("transaction_id IN ?", transactionIDs).
		Order("transaction_id ASC").
		Order("id ASC").
		Find(&entryModels).Error; err != nil {
		return nil, mapError(err)
	}

	for _, entryModel := range entryModels {
		transaction := transactionsByID[entryModel.TransactionID]
		if transaction == nil {
			continue
		}
		entry := entryModel
		transaction.Entries = append(transaction.Entries, &entity.LedgerEntry{
			ID:            entry.ID,
			TransactionID: entry.TransactionID,
			AccountID:     entry.AccountID,
			Currency:      entry.Currency,
			Amount:        entry.Amount,
			CreatedAt:     entry.CreatedAt,
		})
	}

	transactions := make([]*entity.LedgerTransaction, 0, len(transactionRows))
	for _, row := range transactionRows {
		transactions = append(transactions, transactionsByID[row.TransactionID])
	}

	return transactions, nil
}

func (r *ledgerRepoImpl) listTransactionsBaseQuery(ctx context.Context, accountID, currency string) *gorm.DB {
	accountID = strings.TrimSpace(accountID)
	currency = strings.ToUpper(strings.TrimSpace(currency))

	query := r.db.WithContext(ctx).
		Table(model.LedgerTransactionModel{}.TableName()+" t").
		Joins("JOIN "+model.LedgerEntryModel{}.TableName()+" e ON e.transaction_id = t.transaction_id").
		Where("e.account_id = ?", accountID)

	if currency != "" {
		query = query.Where("t.currency = ?", currency)
	}

	return query
}
