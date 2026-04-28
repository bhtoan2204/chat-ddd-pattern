package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	ledgeraggregate "wechat-clone/core/modules/ledger/domain/aggregate"
	"wechat-clone/core/modules/ledger/domain/entity"
	ledgerrepos "wechat-clone/core/modules/ledger/domain/repos"
	valueobject "wechat-clone/core/modules/ledger/domain/value_object"
	sharedevents "wechat-clone/core/shared/contracts/events"
	eventpkg "wechat-clone/core/shared/pkg/event"
	"wechat-clone/core/shared/pkg/stackErr"
)

type TransferToAccountCommand struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Currency      string
	Amount        int64
	FeeAmount     int64
	FeeAccountID  string
}

type RecordPaymentSucceededCommand struct {
	PaymentID          string
	TransactionID      string
	ClearingAccountKey string
	CreditAccountID    string
	Currency           string
	Amount             int64
}

type RecordPaymentReversedCommand struct {
	PaymentID          string
	TransactionID      string
	ClearingAccountKey string
	CreditAccountID    string
	Currency           string
	Amount             int64
	ReversalType       string
}

type RecordLedgerEventsCommand struct {
	Events []eventpkg.Event
}

type RecordPaymentReconciliationFailedCommand struct {
	PaymentID          string
	TransactionID      string
	Provider           string
	ClearingAccountKey string
	CreditAccountID    string
	Currency           string
	Amount             int64
	FeeAmount          int64
	ProviderAmount     int64
	Reason             string
	FailedAt           time.Time
}

const LedgerAccountLockKeyPrefix = "ledger-account"

//go:generate mockgen -package=service -destination=ledger_service_mock.go -source=ledger_service.go
type LedgerService interface {
	TransferToAccount(ctx context.Context, command TransferToAccountCommand) (*entity.LedgerTransaction, error)
	RecordLedgerEvents(ctx context.Context, command RecordLedgerEventsCommand) error
	RecordPaymentSucceeded(ctx context.Context, command RecordPaymentSucceededCommand) error
	RecordPaymentReversed(ctx context.Context, command RecordPaymentReversedCommand) error
	RecordPaymentReconciliationFailed(ctx context.Context, command RecordPaymentReconciliationFailedCommand) error
}

type ledgerService struct {
	baseRepo ledgerrepos.Repos
}

type expectedLedgerPosting struct {
	accountID string
	posting   entity.LedgerAccountPosting
}

type loadedLedgerAccounts map[string]*ledgeraggregate.LedgerAccountAggregate

func NewLedgerService(baseRepo ledgerrepos.Repos) *ledgerService {
	return &ledgerService{baseRepo: baseRepo}
}

func (s *ledgerService) TransferToAccount(ctx context.Context, command TransferToAccountCommand) (*entity.LedgerTransaction, error) {
	booking, err := entity.NewTransferBooking(entity.TransferBookingInput{
		FromAccountID: command.FromAccountID,
		ToAccountID:   command.ToAccountID,
		Currency:      command.Currency,
		Amount:        command.Amount,
	})
	if err != nil {
		return nil, stackErr.Error(fmt.Errorf("%w: %w", ErrValidation, err))
	}
	feeAccountID := strings.TrimSpace(command.FeeAccountID)
	if command.FeeAmount < 0 {
		return nil, stackErr.Error(fmt.Errorf("%w: fee_amount must be greater than or equal to 0", ErrValidation))
	}
	if command.FeeAmount > 0 && feeAccountID == "" {
		return nil, stackErr.Error(fmt.Errorf("%w: fee_account_id is required when fee_amount > 0", ErrValidation))
	}

	transaction, err := entity.NewLedgerTransaction(strings.TrimSpace(command.TransactionID), booking.LedgerEntries())
	if err != nil {
		return nil, stackErr.Error(fmt.Errorf("%w: %w", ErrValidation, err))
	}
	var feeTransaction *entity.LedgerTransaction
	if command.FeeAmount > 0 {
		feeBooking, feeErr := entity.NewTransferBooking(entity.TransferBookingInput{
			FromAccountID: command.FromAccountID,
			ToAccountID:   feeAccountID,
			Currency:      command.Currency,
			Amount:        command.FeeAmount,
		})
		if feeErr != nil {
			return nil, stackErr.Error(fmt.Errorf("%w: %w", ErrValidation, feeErr))
		}
		feeTransaction, err = entity.NewLedgerTransaction(strings.TrimSpace(command.TransactionID)+":fee", feeBooking.LedgerEntries())
		if err != nil {
			return nil, stackErr.Error(fmt.Errorf("%w: %w", ErrValidation, err))
		}
	}

	fromPosting, err := ledgeraggregate.NewLedgerAccountTransferOutPosting(
		valueobject.LedgerAccountTransferPostingInput{
			AccountID:             booking.FromAccountID,
			TransactionID:         transaction.TransactionID,
			CounterpartyAccountID: booking.ToAccountID,
			Currency:              transaction.Currency,
			Amount:                booking.Amount,
			BookedAt:              transaction.CreatedAt,
		},
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	toPosting, err := ledgeraggregate.NewLedgerAccountTransferInPosting(
		valueobject.LedgerAccountTransferPostingInput{
			AccountID:             booking.ToAccountID,
			TransactionID:         transaction.TransactionID,
			CounterpartyAccountID: booking.FromAccountID,
			Currency:              transaction.Currency,
			Amount:                booking.Amount,
			BookedAt:              transaction.CreatedAt,
		},
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	var (
		feeFromPosting entity.LedgerAccountPosting
		feeToPosting   entity.LedgerAccountPosting
	)
	if feeTransaction != nil {
		feeFromPosting, err = ledgeraggregate.NewLedgerAccountTransferOutPosting(
			valueobject.LedgerAccountTransferPostingInput{
				AccountID:             command.FromAccountID,
				TransactionID:         feeTransaction.TransactionID,
				CounterpartyAccountID: feeAccountID,
				Currency:              feeTransaction.Currency,
				Amount:                command.FeeAmount,
				BookedAt:              feeTransaction.CreatedAt,
			},
		)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		feeToPosting, err = ledgeraggregate.NewLedgerAccountTransferInPosting(
			valueobject.LedgerAccountTransferPostingInput{
				AccountID:             feeAccountID,
				TransactionID:         feeTransaction.TransactionID,
				CounterpartyAccountID: command.FromAccountID,
				Currency:              feeTransaction.Currency,
				Amount:                command.FeeAmount,
				BookedAt:              feeTransaction.CreatedAt,
			},
		)
		if err != nil {
			return nil, stackErr.Error(err)
		}
	}

	if err := s.baseRepo.WithTransaction(ctx, func(txRepos ledgerrepos.Repos) error {
		accountIDs := []string{booking.FromAccountID, booking.ToAccountID}
		if feeTransaction != nil {
			accountIDs = append(accountIDs, feeAccountID)
		}
		accounts, err := s.loadLedgerAccounts(ctx, txRepos, accountIDs...)
		if err != nil {
			return stackErr.Error(err)
		}

		expectedPostings := []expectedLedgerPosting{
			{accountID: booking.FromAccountID, posting: fromPosting},
			{accountID: booking.ToAccountID, posting: toPosting},
		}
		if feeTransaction != nil {
			expectedPostings = append(expectedPostings,
				expectedLedgerPosting{accountID: command.FromAccountID, posting: feeFromPosting},
				expectedLedgerPosting{accountID: feeAccountID, posting: feeToPosting},
			)
		}
		alreadyApplied, err := s.ensureLedgerPostingsState(accounts, expectedPostings)
		if err != nil {
			return stackErr.Error(err)
		}
		if alreadyApplied {
			return nil
		}

		fromAgg := accounts.account(booking.FromAccountID)
		toAgg := accounts.account(booking.ToAccountID)
		var feeAgg *ledgeraggregate.LedgerAccountAggregate
		if feeTransaction != nil {
			feeAgg = accounts.account(feeAccountID)
		}

		fromApplied, err := fromAgg.TransferToAccount(
			transaction.TransactionID,
			booking.ToAccountID,
			transaction.Currency,
			booking.Amount,
			transaction.CreatedAt,
		)
		if err != nil {
			if errors.Is(err, ledgeraggregate.ErrLedgerAccountInsufficientFunds) {
				return stackErr.Error(fmt.Errorf("%w: %w", ErrInsufficientFunds, err))
			}
			return stackErr.Error(err)
		}
		toApplied, err := toAgg.ReceiveTransfer(
			transaction.TransactionID,
			booking.FromAccountID,
			transaction.Currency,
			booking.Amount,
			transaction.CreatedAt,
		)
		if err != nil {
			return stackErr.Error(err)
		}
		feeApplied := false
		if feeTransaction != nil {
			feeApplied, err = fromAgg.TransferToAccount(
				feeTransaction.TransactionID,
				feeAccountID,
				feeTransaction.Currency,
				command.FeeAmount,
				feeTransaction.CreatedAt,
			)
			if err != nil {
				if errors.Is(err, ledgeraggregate.ErrLedgerAccountInsufficientFunds) {
					return stackErr.Error(fmt.Errorf("%w: %w", ErrInsufficientFunds, err))
				}
				return stackErr.Error(err)
			}
			feeReceived, receiveErr := feeAgg.ReceiveTransfer(
				feeTransaction.TransactionID,
				command.FromAccountID,
				feeTransaction.Currency,
				command.FeeAmount,
				feeTransaction.CreatedAt,
			)
			if receiveErr != nil {
				return stackErr.Error(receiveErr)
			}
			if feeApplied != feeReceived {
				return stackErr.Error(fmt.Errorf("ledger fee transfer became inconsistent for transaction_id=%s", feeTransaction.TransactionID))
			}
		}
		if fromApplied != toApplied {
			return stackErr.Error(fmt.Errorf("ledger transfer posting became inconsistent for transaction_id=%s", transaction.TransactionID))
		}
		if !fromApplied && !feeApplied {
			return nil
		}

		if alreadyApplied, err := s.saveLedgerAccount(ctx, txRepos, fromAgg); err != nil {
			return stackErr.Error(err)
		} else if alreadyApplied {
			return nil
		}
		if alreadyApplied, err := s.saveLedgerAccount(ctx, txRepos, toAgg); err != nil {
			return stackErr.Error(err)
		} else if alreadyApplied {
			return stackErr.Error(fmt.Errorf("ledger transfer duplicate state became inconsistent for transaction_id=%s", transaction.TransactionID))
		}
		if feeTransaction != nil {
			if alreadyApplied, err := s.saveLedgerAccount(ctx, txRepos, feeAgg); err != nil {
				return stackErr.Error(err)
			} else if alreadyApplied {
				return stackErr.Error(fmt.Errorf("ledger fee duplicate state became inconsistent for transaction_id=%s", feeTransaction.TransactionID))
			}
		}
		return nil
	}); err != nil {
		return nil, stackErr.Error(err)
	}

	if feeTransaction != nil {
		transaction.Entries = append(transaction.Entries, feeTransaction.Entries...)
	}
	return transaction, nil
}

func (s *ledgerService) RecordPaymentSucceeded(ctx context.Context, command RecordPaymentSucceededCommand) error {
	booking, err := entity.NewPaymentSucceededBooking(entity.PaymentSucceededBookingInput{
		PaymentID:          command.PaymentID,
		TransactionID:      command.TransactionID,
		ClearingAccountKey: command.ClearingAccountKey,
		CreditAccountID:    command.CreditAccountID,
		Currency:           command.Currency,
		Amount:             command.Amount,
	})
	if err != nil {
		return stackErr.Error(fmt.Errorf("%w: %w", ErrValidation, err))
	}

	transaction, err := entity.NewLedgerTransaction(booking.LedgerTransactionID(), booking.LedgerEntries())
	if err != nil {
		return stackErr.Error(fmt.Errorf("%w: %w", ErrValidation, err))
	}

	debitPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		valueobject.LedgerAccountPostingInput{
			AccountID:             booking.DebitAccountID,
			TransactionID:         transaction.TransactionID,
			ReferenceType:         ledgeraggregate.EventNameLedgerAccountWithdrawFromIntent,
			ReferenceID:           booking.PaymentID,
			CounterpartyAccountID: booking.CreditAccountID,
			Currency:              transaction.Currency,
			AmountDelta:           -booking.Amount,
			BookedAt:              transaction.CreatedAt,
		},
	)
	if err != nil {
		return stackErr.Error(err)
	}
	creditPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		valueobject.LedgerAccountPostingInput{
			AccountID:             booking.CreditAccountID,
			TransactionID:         transaction.TransactionID,
			ReferenceType:         ledgeraggregate.EventNameLedgerAccountDepositFromIntent,
			ReferenceID:           booking.PaymentID,
			CounterpartyAccountID: booking.DebitAccountID,
			Currency:              transaction.Currency,
			AmountDelta:           booking.Amount,
			BookedAt:              transaction.CreatedAt,
		},
	)
	if err != nil {
		return stackErr.Error(err)
	}

	events, err := ledgerPaymentEventsFromPostings([]ledgerPostingEventInput{
		{accountID: booking.DebitAccountID, posting: debitPosting},
		{accountID: booking.CreditAccountID, posting: creditPosting},
	})
	if err != nil {
		return stackErr.Error(err)
	}

	return stackErr.Error(s.RecordLedgerEvents(ctx, RecordLedgerEventsCommand{Events: events}))
}

func (s *ledgerService) RecordPaymentReversed(ctx context.Context, command RecordPaymentReversedCommand) error {
	booking, err := entity.NewPaymentReversalBooking(entity.PaymentReversalBookingInput{
		PaymentID:          command.PaymentID,
		TransactionID:      command.TransactionID,
		ClearingAccountKey: command.ClearingAccountKey,
		CreditAccountID:    command.CreditAccountID,
		Currency:           command.Currency,
		Amount:             command.Amount,
		ReversalType:       command.ReversalType,
	})
	if err != nil {
		return stackErr.Error(fmt.Errorf("%w: %w", ErrValidation, err))
	}

	transaction, err := entity.NewLedgerTransaction(booking.LedgerTransactionID(), booking.LedgerEntries())
	if err != nil {
		return stackErr.Error(fmt.Errorf("%w: %w", ErrValidation, err))
	}

	debitPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		valueobject.LedgerAccountPostingInput{
			AccountID:             booking.DebitAccountID,
			TransactionID:         transaction.TransactionID,
			ReferenceType:         debitLedgerEventNameForPaymentReversal(booking.ReversalType),
			ReferenceID:           booking.PaymentID,
			CounterpartyAccountID: booking.CreditAccountID,
			Currency:              transaction.Currency,
			AmountDelta:           -booking.Amount,
			BookedAt:              transaction.CreatedAt,
		},
	)
	if err != nil {
		return stackErr.Error(err)
	}
	creditPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		valueobject.LedgerAccountPostingInput{
			AccountID:             booking.CreditAccountID,
			TransactionID:         transaction.TransactionID,
			ReferenceType:         creditLedgerEventNameForPaymentReversal(booking.ReversalType),
			ReferenceID:           booking.PaymentID,
			CounterpartyAccountID: booking.DebitAccountID,
			Currency:              transaction.Currency,
			AmountDelta:           booking.Amount,
			BookedAt:              transaction.CreatedAt,
		},
	)
	if err != nil {
		return stackErr.Error(err)
	}

	events, err := ledgerPaymentEventsFromPostings([]ledgerPostingEventInput{
		{accountID: booking.DebitAccountID, posting: debitPosting},
		{accountID: booking.CreditAccountID, posting: creditPosting},
	})
	if err != nil {
		return stackErr.Error(err)
	}

	return stackErr.Error(s.RecordLedgerEvents(ctx, RecordLedgerEventsCommand{Events: events}))
}

func (s *ledgerService) RecordPaymentReconciliationFailed(ctx context.Context, command RecordPaymentReconciliationFailedCommand) error {
	if s == nil || s.baseRepo == nil {
		return stackErr.Error(fmt.Errorf("ledger repository is required"))
	}

	paymentID := strings.TrimSpace(command.PaymentID)
	if paymentID == "" {
		paymentID = strings.TrimSpace(command.TransactionID)
	}
	if paymentID == "" {
		return stackErr.Error(fmt.Errorf("%w: payment_id is required", ErrValidation))
	}

	failedAt := command.FailedAt.UTC()
	if failedAt.IsZero() {
		failedAt = time.Now().UTC()
	}

	reason := strings.TrimSpace(command.Reason)
	if reason == "" {
		reason = "ledger payment reconciliation failed"
	}

	failureEvent := eventpkg.Event{
		AggregateID:   paymentID,
		AggregateType: "LedgerPaymentReconciliation",
		EventName:     sharedevents.EventLedgerPaymentReconciliationFailed,
		EventData: sharedevents.LedgerPaymentReconciliationFailedEvent{
			PaymentID:          paymentID,
			TransactionID:      strings.TrimSpace(command.TransactionID),
			Provider:           strings.TrimSpace(command.Provider),
			ClearingAccountKey: strings.TrimSpace(command.ClearingAccountKey),
			CreditAccountID:    strings.TrimSpace(command.CreditAccountID),
			Currency:           strings.TrimSpace(command.Currency),
			Amount:             command.Amount,
			FeeAmount:          command.FeeAmount,
			ProviderAmount:     command.ProviderAmount,
			Reason:             reason,
			FailedAt:           failedAt,
		},
		CreatedAt: failedAt.Unix(),
	}

	outboxRepo, ok := ledgerOutboxStore(s.baseRepo)
	if !ok || outboxRepo == nil {
		return stackErr.Error(fmt.Errorf("ledger reconciliation failure publisher is not configured"))
	}

	return stackErr.Error(outboxRepo.Append(ctx, failureEvent))
}

func ledgerOutboxStore(baseRepo ledgerrepos.Repos) (eventpkg.Store, bool) {
	type paymentReconciliationFailureStoreProvider interface {
		PaymentReconciliationFailureEventStore() eventpkg.Store
	}

	provider, ok := baseRepo.(paymentReconciliationFailureStoreProvider)
	if !ok {
		return nil, false
	}
	return provider.PaymentReconciliationFailureEventStore(), true
}

func (s *ledgerService) RecordLedgerEvents(ctx context.Context, command RecordLedgerEventsCommand) error {
	if len(command.Events) == 0 {
		return nil
	}

	return stackErr.Error(s.baseRepo.WithTransaction(ctx, func(txRepos ledgerrepos.Repos) error {
		accountIDs := make([]string, 0, len(command.Events))
		seen := make(map[string]struct{}, len(command.Events))
		for _, evt := range command.Events {
			accountID := strings.TrimSpace(evt.AggregateID)
			if accountID == "" {
				return stackErr.Error(fmt.Errorf("ledger event aggregate_id is required"))
			}
			if _, exists := seen[accountID]; exists {
				continue
			}
			seen[accountID] = struct{}{}
			accountIDs = append(accountIDs, accountID)
		}

		aggregates, err := s.loadLedgerAccounts(ctx, txRepos, accountIDs...)
		if err != nil {
			return stackErr.Error(err)
		}

		saveOrder := make([]string, 0, len(command.Events))
		appliedCount := 0
		for _, evt := range command.Events {
			accountID := strings.TrimSpace(evt.AggregateID)
			agg := aggregates.account(accountID)
			if agg == nil {
				return stackErr.Error(fmt.Errorf("ledger aggregate not loaded: %s", accountID))
			}
			if _, exists := seen[accountID]; exists {
				delete(seen, accountID)
				saveOrder = append(saveOrder, accountID)
			}

			applied, err := agg.ApplyPostingEvent(evt.EventData)
			if err != nil {
				return stackErr.Error(err)
			}
			if applied {
				appliedCount++
			}
		}
		if appliedCount == 0 {
			return nil
		}
		if appliedCount != len(command.Events) {
			return stackErr.Error(fmt.Errorf("ledger event application became inconsistent"))
		}

		for _, accountID := range saveOrder {
			agg := aggregates[accountID]
			if alreadyApplied, err := s.saveLedgerAccount(ctx, txRepos, agg); err != nil {
				return stackErr.Error(err)
			} else if alreadyApplied {
				return stackErr.Error(fmt.Errorf("ledger duplicate state became inconsistent for aggregate_id=%s", accountID))
			}
		}
		return nil
	}))
}

type ledgerPostingEventInput struct {
	accountID string
	posting   entity.LedgerAccountPosting
}

func ledgerPaymentEventsFromPostings(inputs []ledgerPostingEventInput) ([]eventpkg.Event, error) {
	events := make([]eventpkg.Event, 0, len(inputs))
	for _, item := range inputs {
		evt, ok, err := ledgeraggregate.NewLedgerAccountEventFromPosting(item.accountID, item.posting)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		if !ok {
			return nil, stackErr.Error(fmt.Errorf("unsupported ledger posting reference_type=%s", item.posting.ReferenceType))
		}
		events = append(events, evt)
	}
	return events, nil
}

func debitLedgerEventNameForPaymentReversal(paymentEventName string) string {
	switch strings.TrimSpace(paymentEventName) {
	case sharedevents.EventPaymentRefunded:
		return ledgeraggregate.EventNameLedgerAccountWithdrawFromRefund
	case sharedevents.EventPaymentChargeback:
		return ledgeraggregate.EventNameLedgerAccountWithdrawFromChargeback
	default:
		return ""
	}
}

func creditLedgerEventNameForPaymentReversal(paymentEventName string) string {
	switch strings.TrimSpace(paymentEventName) {
	case sharedevents.EventPaymentRefunded:
		return ledgeraggregate.EventNameLedgerAccountDepositFromRefund
	case sharedevents.EventPaymentChargeback:
		return ledgeraggregate.EventNameLedgerAccountDepositFromChargeback
	default:
		return ""
	}
}

func (s *ledgerService) loadLedgerAccount(
	ctx context.Context,
	repos ledgerrepos.Repos,
	accountID string,
) (*ledgeraggregate.LedgerAccountAggregate, error) {
	account, err := repos.LedgerAccountAggregateRepository().Load(ctx, accountID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if account != nil {
		return account, nil
	}
	account, err = ledgeraggregate.NewLedgerAccountAggregate(accountID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return account, nil
}

func (s *ledgerService) ensureLedgerPostingsState(
	accounts loadedLedgerAccounts,
	expected []expectedLedgerPosting,
) (bool, error) {
	matchedCount := 0
	for _, item := range expected {
		account := accounts.account(item.accountID)
		if account == nil {
			continue
		}
		existing, ok := account.PostedTransaction(item.posting.TransactionID)
		if existing == nil {
			continue
		}
		if !ok {
			continue
		}
		if !ledgeraggregate.SameLedgerAccountPosting(*existing, item.posting) {
			return false, stackErr.Error(fmt.Errorf(
				"existing ledger posting mismatch for account_id=%s transaction_id=%s",
				item.accountID,
				item.posting.TransactionID,
			))
		}
		matchedCount++
	}

	if matchedCount == 0 {
		return false, nil
	}
	if matchedCount != len(expected) {
		return false, stackErr.Error(fmt.Errorf(
			"ledger posting state became inconsistent for transaction_id=%s",
			expected[0].posting.TransactionID,
		))
	}

	return true, nil
}

func (s *ledgerService) loadLedgerAccounts(
	ctx context.Context,
	repos ledgerrepos.Repos,
	accountIDs ...string,
) (loadedLedgerAccounts, error) {
	accounts := make(loadedLedgerAccounts, len(accountIDs))
	for _, accountID := range accountIDs {
		accountID = strings.TrimSpace(accountID)
		if accountID == "" {
			continue
		}
		if _, exists := accounts[accountID]; exists {
			continue
		}

		account, err := s.loadLedgerAccount(ctx, repos, accountID)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		accounts[accountID] = account
	}
	return accounts, nil
}

func (a loadedLedgerAccounts) account(accountID string) *ledgeraggregate.LedgerAccountAggregate {
	if a == nil {
		return nil
	}
	return a[strings.TrimSpace(accountID)]
}

func (s *ledgerService) saveLedgerAccount(
	ctx context.Context,
	repos ledgerrepos.Repos,
	account *ledgeraggregate.LedgerAccountAggregate,
) (bool, error) {
	err := repos.LedgerAccountAggregateRepository().Save(ctx, account)
	if err == nil {
		return false, nil
	}
	if errors.Is(err, ledgerrepos.ErrAlreadyApplied) {
		return true, nil
	}
	return false, stackErr.Error(err)
}
