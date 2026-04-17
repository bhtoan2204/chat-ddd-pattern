package aggregate

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"wechat-clone/core/modules/ledger/domain/entity"
	"wechat-clone/core/shared/pkg/event"
	"wechat-clone/core/shared/pkg/stackErr"
)

var (
	ErrLedgerAccountAggregateRequired    = errors.New("ledger account aggregate is required")
	ErrLedgerAccountIDRequired           = errors.New("ledger account id is required")
	ErrLedgerAccountIDMismatch           = errors.New("ledger account id mismatch")
	ErrLedgerAccountTransactionRequired  = errors.New("ledger transaction id is required")
	ErrLedgerAccountReferenceTypeInvalid = errors.New("ledger reference type is invalid")
	ErrLedgerAccountReferenceIDRequired  = errors.New("ledger reference id is required")
	ErrLedgerAccountCounterpartyRequired = errors.New("ledger counterparty_account_id is required")
	ErrLedgerAccountAccountsMustDiffer   = errors.New("ledger account_id and counterparty_account_id must be different")
	ErrLedgerAccountCurrencyRequired     = errors.New("ledger currency is required")
	ErrLedgerAccountAmountInvalid        = errors.New("ledger amount must be positive")
	ErrLedgerAccountBookedAtRequired     = errors.New("ledger booked_at is required")
	ErrLedgerAccountInsufficientFunds    = errors.New("ledger account has insufficient funds")
)

const ledgerReferenceTransferToAccount = "ledger.transfer_to_account"

type LedgerAccountPosting struct {
	TransactionID         string    `json:"transaction_id"`
	ReferenceType         string    `json:"reference_type"`
	ReferenceID           string    `json:"reference_id"`
	CounterpartyAccountID string    `json:"counterparty_account_id"`
	Currency              string    `json:"currency"`
	AmountDelta           int64     `json:"amount_delta"`
	BookedAt              time.Time `json:"booked_at"`
}

type LedgerAccountAggregate struct {
	event.AggregateRoot

	AccountID string           `json:"account_id"`
	Balances  map[string]int64 `json:"balances"`
}

func NewLedgerAccountAggregate(accountID string) (*LedgerAccountAggregate, error) {
	agg := &LedgerAccountAggregate{}
	if err := event.InitAggregate(&agg.AggregateRoot, agg, accountID); err != nil {
		return nil, stackErr.Error(err)
	}
	agg.ensureState()
	return agg, nil
}

func (a *LedgerAccountAggregate) RegisterEvents(register event.RegisterEventsFunc) error {
	return register(
		&EventLedgerAccountPaymentBooked{},
		&EventLedgerAccountTransferredToAccount{},
		&EventLedgerAccountReceivedTransfer{},
	)
}

func (a *LedgerAccountAggregate) Transition(evt event.Event) error {
	switch data := evt.EventData.(type) {
	case *EventLedgerAccountPaymentBooked:
		return a.applyPaymentBooked(evt.AggregateID, data)
	case *EventLedgerAccountTransferredToAccount:
		return a.applyTransferredToAccount(evt.AggregateID, data)
	case *EventLedgerAccountReceivedTransfer:
		return a.applyReceivedTransfer(evt.AggregateID, data)
	default:
		return event.ErrUnsupportedEventType
	}
}

func (a *LedgerAccountAggregate) Balance(currency string) int64 {
	if a == nil {
		return 0
	}
	a.ensureState()
	return a.Balances[strings.ToUpper(strings.TrimSpace(currency))]
}

func (a *LedgerAccountAggregate) BookPayment(
	transactionID string,
	paymentID string,
	counterpartyAccountID string,
	currency string,
	amountDelta int64,
	bookedAt time.Time,
) (bool, error) {
	return a.bookPaymentPosting(
		transactionID,
		entity.PaymentReferenceSucceeded,
		paymentID,
		counterpartyAccountID,
		currency,
		amountDelta,
		bookedAt,
	)
}

func (a *LedgerAccountAggregate) ReversePayment(
	transactionID string,
	referenceType string,
	paymentID string,
	counterpartyAccountID string,
	currency string,
	amountDelta int64,
	bookedAt time.Time,
) (bool, error) {
	referenceType = strings.TrimSpace(referenceType)
	if referenceType != entity.PaymentReferenceRefunded && referenceType != entity.PaymentReferenceChargeback {
		return false, stackErr.Error(fmt.Errorf("%w: %s", ErrLedgerAccountReferenceTypeInvalid, referenceType))
	}

	return a.bookPaymentPosting(
		transactionID,
		referenceType,
		paymentID,
		counterpartyAccountID,
		currency,
		amountDelta,
		bookedAt,
	)
}

func (a *LedgerAccountAggregate) bookPaymentPosting(
	transactionID string,
	referenceType string,
	referenceID string,
	counterpartyAccountID string,
	currency string,
	amountDelta int64,
	bookedAt time.Time,
) (bool, error) {
	if a == nil {
		return false, stackErr.Error(ErrLedgerAccountAggregateRequired)
	}
	posting, err := NewLedgerAccountPaymentPosting(
		a.AggregateID(),
		transactionID,
		referenceType,
		referenceID,
		counterpartyAccountID,
		currency,
		amountDelta,
		bookedAt,
	)
	if err != nil {
		return false, stackErr.Error(err)
	}

	existing, exists := a.lookupPendingPosting(posting.TransactionID)
	if exists {
		if SameLedgerAccountPosting(existing, posting) {
			return false, nil
		}
		return false, stackErr.Error(fmt.Errorf("ledger payment booking mismatch for transaction_id=%s", posting.TransactionID))
	}
	if err := a.ensurePostingAllowed(posting); err != nil {
		return false, stackErr.Error(err)
	}

	if err := a.ApplyChange(a, &EventLedgerAccountPaymentBooked{
		TransactionID:         posting.TransactionID,
		ReferenceType:         posting.ReferenceType,
		PaymentID:             posting.ReferenceID,
		CounterpartyAccountID: posting.CounterpartyAccountID,
		Currency:              posting.Currency,
		AmountDelta:           posting.AmountDelta,
		BookedAt:              posting.BookedAt,
	}); err != nil {
		return false, stackErr.Error(err)
	}

	return true, nil
}

func (a *LedgerAccountAggregate) TransferToAccount(
	transactionID string,
	toAccountID string,
	currency string,
	amount int64,
	bookedAt time.Time,
) (bool, error) {
	if a == nil {
		return false, stackErr.Error(ErrLedgerAccountAggregateRequired)
	}
	posting, err := NewLedgerAccountTransferOutPosting(
		a.AggregateID(),
		transactionID,
		toAccountID,
		currency,
		amount,
		bookedAt,
	)
	if err != nil {
		return false, stackErr.Error(err)
	}

	existing, exists := a.lookupPendingPosting(posting.TransactionID)
	if exists {
		if SameLedgerAccountPosting(existing, posting) {
			return false, nil
		}
		return false, stackErr.Error(fmt.Errorf("ledger transfer mismatch for transaction_id=%s", posting.TransactionID))
	}
	if err := a.ensurePostingAllowed(posting); err != nil {
		return false, stackErr.Error(err)
	}

	if err := a.ApplyChange(a, &EventLedgerAccountTransferredToAccount{
		TransactionID: posting.TransactionID,
		ToAccountID:   posting.CounterpartyAccountID,
		Currency:      posting.Currency,
		Amount:        amount,
		BookedAt:      posting.BookedAt,
	}); err != nil {
		return false, stackErr.Error(err)
	}

	return true, nil
}

func (a *LedgerAccountAggregate) ReceiveTransfer(
	transactionID string,
	fromAccountID string,
	currency string,
	amount int64,
	bookedAt time.Time,
) (bool, error) {
	if a == nil {
		return false, stackErr.Error(ErrLedgerAccountAggregateRequired)
	}
	posting, err := NewLedgerAccountTransferInPosting(
		a.AggregateID(),
		transactionID,
		fromAccountID,
		currency,
		amount,
		bookedAt,
	)
	if err != nil {
		return false, stackErr.Error(err)
	}

	existing, exists := a.lookupPendingPosting(posting.TransactionID)
	if exists {
		if SameLedgerAccountPosting(existing, posting) {
			return false, nil
		}
		return false, stackErr.Error(fmt.Errorf("ledger transfer receive mismatch for transaction_id=%s", posting.TransactionID))
	}
	if err := a.ensurePostingAllowed(posting); err != nil {
		return false, stackErr.Error(err)
	}

	if err := a.ApplyChange(a, &EventLedgerAccountReceivedTransfer{
		TransactionID: posting.TransactionID,
		FromAccountID: posting.CounterpartyAccountID,
		Currency:      posting.Currency,
		Amount:        amount,
		BookedAt:      posting.BookedAt,
	}); err != nil {
		return false, stackErr.Error(err)
	}

	return true, nil
}

func (a *LedgerAccountAggregate) applyPaymentBooked(accountID string, data *EventLedgerAccountPaymentBooked) error {
	posting, ok, err := NewLedgerAccountPostingFromEvent(accountID, data)
	if err != nil {
		return stackErr.Error(err)
	}
	if !ok {
		return stackErr.Error(errors.New("ledger payment booked event is unsupported"))
	}
	return a.applyPosting(accountID, posting)
}

func (a *LedgerAccountAggregate) applyTransferredToAccount(accountID string, data *EventLedgerAccountTransferredToAccount) error {
	posting, ok, err := NewLedgerAccountPostingFromEvent(accountID, data)
	if err != nil {
		return stackErr.Error(err)
	}
	if !ok {
		return stackErr.Error(errors.New("ledger transfer to account event is unsupported"))
	}
	return a.applyPosting(accountID, posting)
}

func (a *LedgerAccountAggregate) applyReceivedTransfer(accountID string, data *EventLedgerAccountReceivedTransfer) error {
	posting, ok, err := NewLedgerAccountPostingFromEvent(accountID, data)
	if err != nil {
		return stackErr.Error(err)
	}
	if !ok {
		return stackErr.Error(errors.New("ledger received transfer event is unsupported"))
	}
	return a.applyPosting(accountID, posting)
}

func (a *LedgerAccountAggregate) applyPosting(accountID string, posting LedgerAccountPosting) error {
	normalizedAccountID, normalizedPosting, err := normalizeLedgerAccountPosting(accountID, posting)
	if err != nil {
		return err
	}
	if err := a.ensurePostingAllowed(normalizedPosting); err != nil {
		return err
	}

	a.ensureState()
	if a.AccountID != "" && a.AccountID != normalizedAccountID {
		return fmt.Errorf("%w: aggregate=%s event=%s", ErrLedgerAccountIDMismatch, a.AccountID, normalizedAccountID)
	}

	a.AccountID = normalizedAccountID
	a.Balances[normalizedPosting.Currency] += normalizedPosting.AmountDelta
	return nil
}

func (a *LedgerAccountAggregate) ensurePostingAllowed(posting LedgerAccountPosting) error {
	if posting.AmountDelta >= 0 {
		return nil
	}
	if !requiresNonNegativeBalance(posting.ReferenceType) {
		return nil
	}

	nextBalance := a.Balance(posting.Currency) + posting.AmountDelta
	if nextBalance < 0 {
		return stackErr.Error(fmt.Errorf(
			"%w: account_id=%s currency=%s balance=%d amount=%d",
			ErrLedgerAccountInsufficientFunds,
			a.AggregateID(),
			posting.Currency,
			a.Balance(posting.Currency),
			-posting.AmountDelta,
		))
	}

	return nil
}

func requiresNonNegativeBalance(referenceType string) bool {
	switch strings.TrimSpace(referenceType) {
	case entity.PaymentReferenceSucceeded, entity.PaymentReferenceRefunded, entity.PaymentReferenceChargeback:
		return false
	default:
		return true
	}
}

func (a *LedgerAccountAggregate) ensureState() {
	if a.Balances == nil {
		a.Balances = make(map[string]int64)
	}
}

func (a *LedgerAccountAggregate) lookupPendingPosting(transactionID string) (LedgerAccountPosting, bool) {
	if a == nil {
		return LedgerAccountPosting{}, false
	}
	transactionID = strings.TrimSpace(transactionID)
	if transactionID == "" {
		return LedgerAccountPosting{}, false
	}

	for _, evt := range a.Root().Events() {
		posting, ok, err := NewLedgerAccountPostingFromEvent(evt.AggregateID, evt.EventData)
		if err != nil || !ok {
			continue
		}
		if posting.TransactionID == transactionID {
			return posting, true
		}
	}

	return LedgerAccountPosting{}, false
}

func normalizeLedgerAccountPosting(accountID string, posting LedgerAccountPosting) (string, LedgerAccountPosting, error) {
	accountID = strings.TrimSpace(accountID)
	if accountID == "" {
		return "", LedgerAccountPosting{}, ErrLedgerAccountIDRequired
	}

	normalizedReferenceType := strings.TrimSpace(posting.ReferenceType)
	switch normalizedReferenceType {
	case entity.PaymentReferenceSucceeded, entity.PaymentReferenceRefunded, entity.PaymentReferenceChargeback, ledgerReferenceTransferToAccount:
	default:
		return "", LedgerAccountPosting{}, ErrLedgerAccountReferenceTypeInvalid
	}

	normalizedTransactionID := strings.TrimSpace(posting.TransactionID)
	if normalizedTransactionID == "" {
		return "", LedgerAccountPosting{}, ErrLedgerAccountTransactionRequired
	}

	normalizedReferenceID := strings.TrimSpace(posting.ReferenceID)
	if normalizedReferenceID == "" {
		return "", LedgerAccountPosting{}, ErrLedgerAccountReferenceIDRequired
	}

	normalizedCounterpartyAccountID := strings.TrimSpace(posting.CounterpartyAccountID)
	if normalizedCounterpartyAccountID == "" {
		return "", LedgerAccountPosting{}, ErrLedgerAccountCounterpartyRequired
	}
	if normalizedCounterpartyAccountID == accountID {
		return "", LedgerAccountPosting{}, ErrLedgerAccountAccountsMustDiffer
	}

	normalizedCurrency := strings.ToUpper(strings.TrimSpace(posting.Currency))
	if normalizedCurrency == "" {
		return "", LedgerAccountPosting{}, ErrLedgerAccountCurrencyRequired
	}
	if posting.AmountDelta == 0 {
		return "", LedgerAccountPosting{}, ErrLedgerAccountAmountInvalid
	}
	if posting.BookedAt.IsZero() {
		return "", LedgerAccountPosting{}, ErrLedgerAccountBookedAtRequired
	}

	return accountID, LedgerAccountPosting{
		TransactionID:         normalizedTransactionID,
		ReferenceType:         normalizedReferenceType,
		ReferenceID:           normalizedReferenceID,
		CounterpartyAccountID: normalizedCounterpartyAccountID,
		Currency:              normalizedCurrency,
		AmountDelta:           posting.AmountDelta,
		BookedAt:              posting.BookedAt.UTC(),
	}, nil
}

func NewLedgerAccountPaymentPosting(
	accountID string,
	transactionID string,
	referenceType string,
	referenceID string,
	counterpartyAccountID string,
	currency string,
	amountDelta int64,
	bookedAt time.Time,
) (LedgerAccountPosting, error) {
	posting, err := newLedgerPosting(
		transactionID,
		referenceType,
		referenceID,
		counterpartyAccountID,
		currency,
		amountDelta,
		bookedAt,
	)
	if err != nil {
		return LedgerAccountPosting{}, stackErr.Error(err)
	}

	_, normalizedPosting, err := normalizeLedgerAccountPosting(accountID, posting)
	if err != nil {
		return LedgerAccountPosting{}, stackErr.Error(err)
	}

	return normalizedPosting, nil
}

func NewLedgerAccountTransferOutPosting(
	accountID string,
	transactionID string,
	toAccountID string,
	currency string,
	amount int64,
	bookedAt time.Time,
) (LedgerAccountPosting, error) {
	return NewLedgerAccountPaymentPosting(
		accountID,
		transactionID,
		ledgerReferenceTransferToAccount,
		transactionID,
		toAccountID,
		currency,
		-amount,
		bookedAt,
	)
}

func NewLedgerAccountTransferInPosting(
	accountID string,
	transactionID string,
	fromAccountID string,
	currency string,
	amount int64,
	bookedAt time.Time,
) (LedgerAccountPosting, error) {
	return NewLedgerAccountPaymentPosting(
		accountID,
		transactionID,
		ledgerReferenceTransferToAccount,
		transactionID,
		fromAccountID,
		currency,
		amount,
		bookedAt,
	)
}

func NewLedgerAccountPostingFromEvent(accountID string, eventData interface{}) (LedgerAccountPosting, bool, error) {
	switch data := eventData.(type) {
	case *EventLedgerAccountPaymentBooked:
		if data == nil {
			return LedgerAccountPosting{}, false, stackErr.Error(errors.New("ledger payment booked event is nil"))
		}

		referenceType := strings.TrimSpace(data.ReferenceType)
		if referenceType == "" {
			referenceType = entity.PaymentReferenceSucceeded
		}

		posting, err := NewLedgerAccountPaymentPosting(
			accountID,
			data.TransactionID,
			referenceType,
			data.PaymentID,
			data.CounterpartyAccountID,
			data.Currency,
			data.AmountDelta,
			data.BookedAt,
		)
		return posting, true, stackErr.Error(err)
	case *EventLedgerAccountTransferredToAccount:
		if data == nil {
			return LedgerAccountPosting{}, false, stackErr.Error(errors.New("ledger transfer to account event is nil"))
		}

		posting, err := NewLedgerAccountTransferOutPosting(
			accountID,
			data.TransactionID,
			data.ToAccountID,
			data.Currency,
			data.Amount,
			data.BookedAt,
		)
		return posting, true, stackErr.Error(err)
	case *EventLedgerAccountReceivedTransfer:
		if data == nil {
			return LedgerAccountPosting{}, false, stackErr.Error(errors.New("ledger received transfer event is nil"))
		}

		posting, err := NewLedgerAccountTransferInPosting(
			accountID,
			data.TransactionID,
			data.FromAccountID,
			data.Currency,
			data.Amount,
			data.BookedAt,
		)
		return posting, true, stackErr.Error(err)
	default:
		return LedgerAccountPosting{}, false, nil
	}
}

func newLedgerPosting(
	transactionID string,
	referenceType string,
	referenceID string,
	counterpartyAccountID string,
	currency string,
	amountDelta int64,
	bookedAt time.Time,
) (LedgerAccountPosting, error) {
	transactionID = strings.TrimSpace(transactionID)
	if transactionID == "" {
		return LedgerAccountPosting{}, ErrLedgerAccountTransactionRequired
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return LedgerAccountPosting{}, ErrLedgerAccountCurrencyRequired
	}
	if amountDelta == 0 {
		return LedgerAccountPosting{}, ErrLedgerAccountAmountInvalid
	}
	if bookedAt.IsZero() {
		return LedgerAccountPosting{}, ErrLedgerAccountBookedAtRequired
	}

	return LedgerAccountPosting{
		TransactionID:         transactionID,
		ReferenceType:         strings.TrimSpace(referenceType),
		ReferenceID:           strings.TrimSpace(referenceID),
		CounterpartyAccountID: strings.TrimSpace(counterpartyAccountID),
		Currency:              currency,
		AmountDelta:           amountDelta,
		BookedAt:              bookedAt.UTC(),
	}, nil
}

func SameLedgerAccountPosting(left LedgerAccountPosting, right LedgerAccountPosting) bool {
	return left.TransactionID == right.TransactionID &&
		left.ReferenceType == right.ReferenceType &&
		left.ReferenceID == right.ReferenceID &&
		left.CounterpartyAccountID == right.CounterpartyAccountID &&
		left.Currency == right.Currency &&
		left.AmountDelta == right.AmountDelta
}
