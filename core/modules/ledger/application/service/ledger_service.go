package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	ledgerprojection "wechat-clone/core/modules/ledger/application/projection"
	ledgeraggregate "wechat-clone/core/modules/ledger/domain/aggregate"
	"wechat-clone/core/modules/ledger/domain/entity"
	ledgerrepos "wechat-clone/core/modules/ledger/domain/repos"
	eventpkg "wechat-clone/core/shared/pkg/event"
	"wechat-clone/core/shared/pkg/stackErr"
)

type TransferToAccountCommand struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Currency      string
	Amount        int64
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

const LedgerAccountLockKeyPrefix = "ledger-account"

//go:generate mockgen -package=service -destination=ledger_service_mock.go -source=ledger_service.go
type LedgerService interface {
	TransferToAccount(ctx context.Context, command TransferToAccountCommand) (*entity.LedgerTransaction, error)
	RecordPaymentSucceeded(ctx context.Context, command RecordPaymentSucceededCommand) error
	RecordPaymentReversed(ctx context.Context, command RecordPaymentReversedCommand) error
}

type ledgerService struct {
	baseRepo ledgerrepos.Repos
}

type expectedLedgerPosting struct {
	accountID string
	posting   ledgeraggregate.LedgerAccountPosting
}

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

	transaction, err := entity.NewLedgerTransaction(strings.TrimSpace(command.TransactionID), booking.LedgerEntries())
	if err != nil {
		return nil, stackErr.Error(fmt.Errorf("%w: %w", ErrValidation, err))
	}

	fromPosting, err := ledgeraggregate.NewLedgerAccountTransferOutPosting(
		booking.FromAccountID,
		transaction.TransactionID,
		booking.ToAccountID,
		transaction.Currency,
		booking.Amount,
		transaction.CreatedAt,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	toPosting, err := ledgeraggregate.NewLedgerAccountTransferInPosting(
		booking.ToAccountID,
		transaction.TransactionID,
		booking.FromAccountID,
		transaction.Currency,
		booking.Amount,
		transaction.CreatedAt,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	if err := s.baseRepo.WithTransaction(ctx, func(txRepos ledgerrepos.Repos) error {
		alreadyApplied, err := s.ensureLedgerPostingsState(ctx, txRepos, []expectedLedgerPosting{
			{accountID: booking.FromAccountID, posting: fromPosting},
			{accountID: booking.ToAccountID, posting: toPosting},
		})
		if err != nil {
			return stackErr.Error(err)
		}
		if alreadyApplied {
			return nil
		}

		fromAgg, err := s.loadLedgerAccount(ctx, txRepos, booking.FromAccountID)
		if err != nil {
			return stackErr.Error(err)
		}
		toAgg, err := s.loadLedgerAccount(ctx, txRepos, booking.ToAccountID)
		if err != nil {
			return stackErr.Error(err)
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
		if fromApplied != toApplied {
			return stackErr.Error(fmt.Errorf("ledger transfer posting became inconsistent for transaction_id=%s", transaction.TransactionID))
		}
		if !fromApplied {
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
		if err := txRepos.LedgerOutboxEventsRepository().Append(ctx, newLedgerTransactionProjectedEvent(
			transaction,
			ledgeraggregate.EventNameLedgerAccountTransferredToAccount,
			"ledger.transfer_to_account",
			transaction.TransactionID,
		)); err != nil {
			return stackErr.Error(err)
		}

		return nil
	}); err != nil {
		return nil, stackErr.Error(err)
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
		booking.DebitAccountID,
		transaction.TransactionID,
		entity.PaymentReferenceSucceeded,
		booking.PaymentID,
		booking.CreditAccountID,
		transaction.Currency,
		-booking.Amount,
		transaction.CreatedAt,
	)
	if err != nil {
		return stackErr.Error(err)
	}
	creditPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		booking.CreditAccountID,
		transaction.TransactionID,
		entity.PaymentReferenceSucceeded,
		booking.PaymentID,
		booking.DebitAccountID,
		transaction.Currency,
		booking.Amount,
		transaction.CreatedAt,
	)
	if err != nil {
		return stackErr.Error(err)
	}

	if err := s.baseRepo.WithTransaction(ctx, func(txRepos ledgerrepos.Repos) error {
		alreadyApplied, err := s.ensureLedgerPostingsState(ctx, txRepos, []expectedLedgerPosting{
			{accountID: booking.DebitAccountID, posting: debitPosting},
			{accountID: booking.CreditAccountID, posting: creditPosting},
		})
		if err != nil {
			return stackErr.Error(err)
		}
		if alreadyApplied {
			return nil
		}

		debitAgg, err := s.loadLedgerAccount(ctx, txRepos, booking.DebitAccountID)
		if err != nil {
			return stackErr.Error(err)
		}
		creditAgg, err := s.loadLedgerAccount(ctx, txRepos, booking.CreditAccountID)
		if err != nil {
			return stackErr.Error(err)
		}

		debitApplied, err := debitAgg.BookPayment(
			transaction.TransactionID,
			booking.PaymentID,
			booking.CreditAccountID,
			transaction.Currency,
			-booking.Amount,
			transaction.CreatedAt,
		)
		if err != nil {
			return stackErr.Error(err)
		}
		creditApplied, err := creditAgg.BookPayment(
			transaction.TransactionID,
			booking.PaymentID,
			booking.DebitAccountID,
			transaction.Currency,
			booking.Amount,
			transaction.CreatedAt,
		)
		if err != nil {
			return stackErr.Error(err)
		}
		if debitApplied != creditApplied {
			return stackErr.Error(fmt.Errorf("ledger payment booking became inconsistent for transaction_id=%s", transaction.TransactionID))
		}
		if !debitApplied {
			return nil
		}

		if alreadyApplied, err := s.saveLedgerAccount(ctx, txRepos, debitAgg); err != nil {
			return stackErr.Error(err)
		} else if alreadyApplied {
			return nil
		}
		if alreadyApplied, err := s.saveLedgerAccount(ctx, txRepos, creditAgg); err != nil {
			return stackErr.Error(err)
		} else if alreadyApplied {
			return stackErr.Error(fmt.Errorf("ledger payment duplicate state became inconsistent for transaction_id=%s", transaction.TransactionID))
		}
		if err := txRepos.LedgerOutboxEventsRepository().Append(ctx, newLedgerTransactionProjectedEvent(
			transaction,
			ledgeraggregate.EventNameLedgerAccountPaymentBooked,
			entity.PaymentReferenceSucceeded,
			booking.PaymentID,
		)); err != nil {
			return stackErr.Error(err)
		}

		return nil
	}); err != nil {
		return stackErr.Error(err)
	}

	return nil
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
		booking.DebitAccountID,
		transaction.TransactionID,
		booking.ReversalType,
		booking.PaymentID,
		booking.CreditAccountID,
		transaction.Currency,
		-booking.Amount,
		transaction.CreatedAt,
	)
	if err != nil {
		return stackErr.Error(err)
	}
	creditPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		booking.CreditAccountID,
		transaction.TransactionID,
		booking.ReversalType,
		booking.PaymentID,
		booking.DebitAccountID,
		transaction.Currency,
		booking.Amount,
		transaction.CreatedAt,
	)
	if err != nil {
		return stackErr.Error(err)
	}

	if err := s.baseRepo.WithTransaction(ctx, func(txRepos ledgerrepos.Repos) error {
		alreadyApplied, err := s.ensureLedgerPostingsState(ctx, txRepos, []expectedLedgerPosting{
			{accountID: booking.DebitAccountID, posting: debitPosting},
			{accountID: booking.CreditAccountID, posting: creditPosting},
		})
		if err != nil {
			return stackErr.Error(err)
		}
		if alreadyApplied {
			return nil
		}

		debitAgg, err := s.loadLedgerAccount(ctx, txRepos, booking.DebitAccountID)
		if err != nil {
			return stackErr.Error(err)
		}
		creditAgg, err := s.loadLedgerAccount(ctx, txRepos, booking.CreditAccountID)
		if err != nil {
			return stackErr.Error(err)
		}

		debitApplied, err := debitAgg.ReversePayment(
			transaction.TransactionID,
			booking.ReversalType,
			booking.PaymentID,
			booking.CreditAccountID,
			transaction.Currency,
			-booking.Amount,
			transaction.CreatedAt,
		)
		if err != nil {
			return stackErr.Error(err)
		}
		creditApplied, err := creditAgg.ReversePayment(
			transaction.TransactionID,
			booking.ReversalType,
			booking.PaymentID,
			booking.DebitAccountID,
			transaction.Currency,
			booking.Amount,
			transaction.CreatedAt,
		)
		if err != nil {
			return stackErr.Error(err)
		}
		if debitApplied != creditApplied {
			return stackErr.Error(fmt.Errorf("ledger payment reversal became inconsistent for transaction_id=%s", transaction.TransactionID))
		}
		if !debitApplied {
			return nil
		}

		if alreadyApplied, err := s.saveLedgerAccount(ctx, txRepos, debitAgg); err != nil {
			return stackErr.Error(err)
		} else if alreadyApplied {
			return nil
		}
		if alreadyApplied, err := s.saveLedgerAccount(ctx, txRepos, creditAgg); err != nil {
			return stackErr.Error(err)
		} else if alreadyApplied {
			return stackErr.Error(fmt.Errorf("ledger payment reversal duplicate state became inconsistent for transaction_id=%s", transaction.TransactionID))
		}
		if err := txRepos.LedgerOutboxEventsRepository().Append(ctx, newLedgerTransactionProjectedEvent(
			transaction,
			ledgeraggregate.EventNameLedgerAccountPaymentBooked,
			booking.ReversalType,
			booking.PaymentID,
		)); err != nil {
			return stackErr.Error(err)
		}

		return nil
	}); err != nil {
		return stackErr.Error(err)
	}

	return nil
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
	ctx context.Context,
	repos ledgerrepos.Repos,
	expected []expectedLedgerPosting,
) (bool, error) {
	matchedCount := 0
	for _, item := range expected {
		existing, err := repos.LedgerAccountAggregateRepository().FindPostedTransaction(ctx, item.accountID, item.posting.TransactionID)
		if err != nil {
			return false, stackErr.Error(err)
		}
		if existing == nil {
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

func newLedgerTransactionProjectedEvent(
	transaction *entity.LedgerTransaction,
	eventName string,
	referenceType string,
	referenceID string,
) eventpkg.Event {
	entries := make([]ledgerprojection.LedgerTransactionEntry, 0, len(transaction.Entries))
	for _, entry := range transaction.Entries {
		if entry == nil {
			continue
		}
		entries = append(entries, ledgerprojection.LedgerTransactionEntry{
			AccountID: entry.AccountID,
			Currency:  entry.Currency,
			Amount:    entry.Amount,
			CreatedAt: entry.CreatedAt.UTC(),
		})
	}

	return eventpkg.Event{
		AggregateID:   transaction.TransactionID,
		AggregateType: ledgerprojection.AggregateTypeLedgerTransactionProjection,
		Version:       1,
		EventName:     strings.TrimSpace(eventName),
		EventData: &ledgerprojection.LedgerTransactionProjected{
			TransactionID: transaction.TransactionID,
			ReferenceType: strings.TrimSpace(referenceType),
			ReferenceID:   strings.TrimSpace(referenceID),
			Currency:      transaction.Currency,
			CreatedAt:     transaction.CreatedAt.UTC(),
			Entries:       entries,
		},
		CreatedAt: transaction.CreatedAt.Unix(),
	}
}
