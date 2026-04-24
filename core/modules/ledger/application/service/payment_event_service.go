package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	ledgeraggregate "wechat-clone/core/modules/ledger/domain/aggregate"
	ledgerentity "wechat-clone/core/modules/ledger/domain/entity"
	ledgerrepos "wechat-clone/core/modules/ledger/domain/repos"
	valueobject "wechat-clone/core/modules/ledger/domain/value_object"
	paymententity "wechat-clone/core/modules/payment/domain/entity"
	sharedevents "wechat-clone/core/shared/contracts/events"
	eventpkg "wechat-clone/core/shared/pkg/event"
	"wechat-clone/core/shared/pkg/stackErr"
	"wechat-clone/core/shared/utils"
)

type PaymentEventService interface {
	HandleWithdrawalRequested(ctx context.Context, payload sharedevents.PaymentWithdrawalRequestedEvent) error
	HandleSucceeded(ctx context.Context, payload sharedevents.PaymentSucceededEvent) error
	HandleFailed(ctx context.Context, payload sharedevents.PaymentFailedEvent) error
	HandleRefunded(ctx context.Context, payload sharedevents.PaymentRefundedEvent) error
	HandleChargeback(ctx context.Context, payload sharedevents.PaymentChargebackEvent) error
}

type paymentEventService struct {
	ledgerService LedgerService
	feeAccountID  string
}

func NewPaymentEventService(baseRepo ledgerrepos.Repos, feeAccountID string) PaymentEventService {
	if baseRepo == nil {
		return nil
	}

	return &paymentEventService{
		ledgerService: NewLedgerService(baseRepo),
		feeAccountID:  strings.TrimSpace(feeAccountID),
	}
}

func (s *paymentEventService) HandleWithdrawalRequested(ctx context.Context, payload sharedevents.PaymentWithdrawalRequestedEvent) error {
	events, err := s.paymentWithdrawalRequestedLedgerEvents(payload)
	if err != nil {
		return stackErr.Error(err)
	}
	return stackErr.Error(s.ledgerService.RecordLedgerEvents(ctx, RecordLedgerEventsCommand{Events: events}))
}

func (s *paymentEventService) HandleSucceeded(ctx context.Context, payload sharedevents.PaymentSucceededEvent) error {
	events, err := s.paymentSucceededLedgerEvents(payload)
	if err != nil {
		return stackErr.Error(err)
	}
	if len(events) == 0 {
		return nil
	}
	return stackErr.Error(s.ledgerService.RecordLedgerEvents(ctx, RecordLedgerEventsCommand{Events: events}))
}

func (s *paymentEventService) HandleFailed(ctx context.Context, payload sharedevents.PaymentFailedEvent) error {
	events, err := s.paymentFailedLedgerEvents(payload)
	if err != nil {
		return stackErr.Error(err)
	}
	if len(events) == 0 {
		return nil
	}
	return stackErr.Error(s.ledgerService.RecordLedgerEvents(ctx, RecordLedgerEventsCommand{Events: events}))
}

func (s *paymentEventService) HandleRefunded(ctx context.Context, payload sharedevents.PaymentRefundedEvent) error {
	events, err := s.paymentReversedLedgerEvents(
		resolvePaymentEventID(payload.PaymentID, payload.TransactionID),
		payload.TransactionID,
		payload.ClearingAccountKey,
		payload.CreditAccountID,
		payload.Currency,
		payload.Amount,
		payload.FeeAmount,
		sharedevents.EventPaymentRefunded,
		payload.RefundedAt,
	)
	if err != nil {
		return stackErr.Error(err)
	}
	return stackErr.Error(s.ledgerService.RecordLedgerEvents(ctx, RecordLedgerEventsCommand{Events: events}))
}

func (s *paymentEventService) HandleChargeback(ctx context.Context, payload sharedevents.PaymentChargebackEvent) error {
	events, err := s.paymentReversedLedgerEvents(
		resolvePaymentEventID(payload.PaymentID, payload.TransactionID),
		payload.TransactionID,
		payload.ClearingAccountKey,
		payload.CreditAccountID,
		payload.Currency,
		payload.Amount,
		payload.FeeAmount,
		sharedevents.EventPaymentChargeback,
		payload.ChargedBackAt,
	)
	if err != nil {
		return stackErr.Error(err)
	}
	return stackErr.Error(s.ledgerService.RecordLedgerEvents(ctx, RecordLedgerEventsCommand{Events: events}))
}

func (s *paymentEventService) paymentWithdrawalRequestedLedgerEvents(payload sharedevents.PaymentWithdrawalRequestedEvent) ([]eventpkg.Event, error) {
	events := make([]eventpkg.Event, 0, 4)
	clearingAccountID := ledgerClearingAccountID(utils.FirstNonEmpty(strings.TrimSpace(payload.ClearingAccountKey), providerClearingAccountKey(payload.Provider)))

	principalEvents, err := paymentLedgerEventsFromPostings([]paymentEventLedgerPostingInput{
		{
			accountID: strings.TrimSpace(payload.DebitAccountID),
			posting: newPaymentPosting(
				strings.TrimSpace(payload.DebitAccountID),
				fmt.Sprintf("payment:%s:withdrawal:principal", resolvePaymentEventID(payload.PaymentID, payload.TransactionID)),
				ledgeraggregate.EventNameLedgerAccountReserveWithdrawal,
				resolvePaymentEventID(payload.PaymentID, payload.TransactionID),
				clearingAccountID,
				payload.Currency,
				-payload.Amount,
				payload.RequestedAt,
			),
		},
		{
			accountID: clearingAccountID,
			posting: newPaymentPosting(
				clearingAccountID,
				fmt.Sprintf("payment:%s:withdrawal:principal", resolvePaymentEventID(payload.PaymentID, payload.TransactionID)),
				ledgeraggregate.EventNameLedgerAccountReceiveWithdrawalHold,
				resolvePaymentEventID(payload.PaymentID, payload.TransactionID),
				strings.TrimSpace(payload.DebitAccountID),
				payload.Currency,
				payload.Amount,
				payload.RequestedAt,
			),
		},
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}
	events = append(events, principalEvents...)

	if payload.FeeAmount > 0 && s.feeAccountID != "" {
		feeEvents, err := paymentLedgerEventsFromPostings([]paymentEventLedgerPostingInput{
			{
				accountID: strings.TrimSpace(payload.DebitAccountID),
				posting: newPaymentPosting(
					strings.TrimSpace(payload.DebitAccountID),
					fmt.Sprintf("payment:%s:withdrawal:fee", resolvePaymentEventID(payload.PaymentID, payload.TransactionID)),
					ledgeraggregate.EventNameLedgerAccountReserveWithdrawal,
					resolvePaymentEventID(payload.PaymentID, payload.TransactionID),
					s.feeAccountID,
					payload.Currency,
					-payload.FeeAmount,
					payload.RequestedAt,
				),
			},
			{
				accountID: s.feeAccountID,
				posting: newPaymentPosting(
					s.feeAccountID,
					fmt.Sprintf("payment:%s:withdrawal:fee", resolvePaymentEventID(payload.PaymentID, payload.TransactionID)),
					ledgeraggregate.EventNameLedgerAccountReceiveWithdrawalHold,
					resolvePaymentEventID(payload.PaymentID, payload.TransactionID),
					strings.TrimSpace(payload.DebitAccountID),
					payload.Currency,
					payload.FeeAmount,
					payload.RequestedAt,
				),
			},
		})
		if err != nil {
			return nil, stackErr.Error(err)
		}
		events = append(events, feeEvents...)
	}

	return events, nil
}

func (s *paymentEventService) paymentSucceededLedgerEvents(payload sharedevents.PaymentSucceededEvent) ([]eventpkg.Event, error) {
	if paymententity.NormalizePaymentWorkflow(payload.Workflow) == paymententity.PaymentWorkflowWithdrawal {
		return nil, nil
	}

	booking, err := ledgerentity.NewPaymentSucceededBooking(ledgerentity.PaymentSucceededBookingInput{
		PaymentID:          resolvePaymentEventID(payload.PaymentID, payload.TransactionID),
		TransactionID:      payload.TransactionID,
		ClearingAccountKey: payload.ClearingAccountKey,
		CreditAccountID:    payload.CreditAccountID,
		Currency:           payload.Currency,
		Amount:             payload.Amount,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	events, err := paymentLedgerEventsFromBooking(
		booking.LedgerTransactionID(),
		booking.PaymentID,
		booking.Currency,
		booking.Amount,
		booking.DebitAccountID,
		booking.CreditAccountID,
		ledgeraggregate.EventNameLedgerAccountWithdrawFromIntent,
		ledgeraggregate.EventNameLedgerAccountDepositFromIntent,
		payload.SucceededAt,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	if payload.FeeAmount > 0 && s.feeAccountID != "" {
		feeEvents, err := paymentLedgerEventsFromBooking(
			fmt.Sprintf("payment:%s:succeeded:fee", booking.PaymentID),
			booking.PaymentID,
			booking.Currency,
			payload.FeeAmount,
			ledgerClearingAccountID(payload.ClearingAccountKey),
			s.feeAccountID,
			ledgeraggregate.EventNameLedgerAccountWithdrawFromIntent,
			ledgeraggregate.EventNameLedgerAccountDepositFromIntent,
			payload.SucceededAt,
		)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		events = append(events, feeEvents...)
	}

	return events, nil
}

func (s *paymentEventService) paymentFailedLedgerEvents(payload sharedevents.PaymentFailedEvent) ([]eventpkg.Event, error) {
	if paymententity.NormalizePaymentWorkflow(payload.Workflow) != paymententity.PaymentWorkflowWithdrawal {
		return nil, nil
	}

	paymentID := resolvePaymentEventID(payload.PaymentID, payload.TransactionID)
	events := make([]eventpkg.Event, 0, 4)
	principalEvents, err := paymentLedgerEventsFromPostings([]paymentEventLedgerPostingInput{
		{
			accountID: strings.TrimSpace(payload.DebitAccountID),
			posting: newPaymentPosting(
				strings.TrimSpace(payload.DebitAccountID),
				fmt.Sprintf("payment:%s:withdrawal:principal:failed", paymentID),
				ledgeraggregate.EventNameLedgerAccountReleaseWithdrawal,
				paymentID,
				ledgerClearingAccountID(payload.ClearingAccountKey),
				payload.Currency,
				payload.Amount,
				payload.OccurredAt,
			),
		},
		{
			accountID: ledgerClearingAccountID(payload.ClearingAccountKey),
			posting: newPaymentPosting(
				ledgerClearingAccountID(payload.ClearingAccountKey),
				fmt.Sprintf("payment:%s:withdrawal:principal:failed", paymentID),
				ledgeraggregate.EventNameLedgerAccountWithdrawReleasedHold,
				paymentID,
				strings.TrimSpace(payload.DebitAccountID),
				payload.Currency,
				-payload.Amount,
				payload.OccurredAt,
			),
		},
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}
	events = append(events, principalEvents...)

	if payload.FeeAmount > 0 && s.feeAccountID != "" {
		feeEvents, err := paymentLedgerEventsFromPostings([]paymentEventLedgerPostingInput{
			{
				accountID: strings.TrimSpace(payload.DebitAccountID),
				posting: newPaymentPosting(
					strings.TrimSpace(payload.DebitAccountID),
					fmt.Sprintf("payment:%s:withdrawal:fee:failed", paymentID),
					ledgeraggregate.EventNameLedgerAccountReleaseWithdrawal,
					paymentID,
					s.feeAccountID,
					payload.Currency,
					payload.FeeAmount,
					payload.OccurredAt,
				),
			},
			{
				accountID: s.feeAccountID,
				posting: newPaymentPosting(
					s.feeAccountID,
					fmt.Sprintf("payment:%s:withdrawal:fee:failed", paymentID),
					ledgeraggregate.EventNameLedgerAccountWithdrawReleasedHold,
					paymentID,
					strings.TrimSpace(payload.DebitAccountID),
					payload.Currency,
					-payload.FeeAmount,
					payload.OccurredAt,
				),
			},
		})
		if err != nil {
			return nil, stackErr.Error(err)
		}
		events = append(events, feeEvents...)
	}

	return events, nil
}

func (s *paymentEventService) paymentReversedLedgerEvents(
	paymentID,
	transactionID,
	clearingAccountKey,
	creditAccountID,
	currency string,
	amount,
	feeAmount int64,
	reversalType string,
	bookedAt time.Time,
) ([]eventpkg.Event, error) {
	booking, err := ledgerentity.NewPaymentReversalBooking(ledgerentity.PaymentReversalBookingInput{
		PaymentID:          paymentID,
		TransactionID:      transactionID,
		ClearingAccountKey: clearingAccountKey,
		CreditAccountID:    creditAccountID,
		Currency:           currency,
		Amount:             amount,
		ReversalType:       reversalType,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	events, err := paymentLedgerEventsFromBooking(
		booking.LedgerTransactionID(),
		booking.PaymentID,
		booking.Currency,
		booking.Amount,
		booking.DebitAccountID,
		booking.CreditAccountID,
		debitLedgerEventNameForReversal(booking.ReversalType),
		creditLedgerEventNameForReversal(booking.ReversalType),
		bookedAt,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	if feeAmount > 0 && s.feeAccountID != "" {
		feeEvents, err := paymentLedgerEventsFromBooking(
			fmt.Sprintf("payment:%s:%s:fee", booking.PaymentID, reversalSuffix(reversalType)),
			booking.PaymentID,
			booking.Currency,
			feeAmount,
			s.feeAccountID,
			ledgerClearingAccountID(clearingAccountKey),
			debitLedgerEventNameForReversal(reversalType),
			creditLedgerEventNameForReversal(reversalType),
			bookedAt,
		)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		events = append(events, feeEvents...)
	}

	return events, nil
}

type paymentEventLedgerPostingInput struct {
	accountID string
	posting   ledgerentity.LedgerAccountPosting
}

func paymentLedgerEventsFromBooking(
	transactionID,
	paymentID,
	currency string,
	amount int64,
	debitAccountID,
	creditAccountID,
	debitEventName,
	creditEventName string,
	bookedAt time.Time,
) ([]eventpkg.Event, error) {
	if bookedAt.IsZero() {
		bookedAt = time.Now().UTC()
	}

	debitPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		valueobject.LedgerAccountPostingInput{
			AccountID:             debitAccountID,
			TransactionID:         transactionID,
			ReferenceType:         debitEventName,
			ReferenceID:           paymentID,
			CounterpartyAccountID: creditAccountID,
			Currency:              currency,
			AmountDelta:           -amount,
			BookedAt:              bookedAt,
		},
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	creditPosting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(
		valueobject.LedgerAccountPostingInput{
			AccountID:             creditAccountID,
			TransactionID:         transactionID,
			ReferenceType:         creditEventName,
			ReferenceID:           paymentID,
			CounterpartyAccountID: debitAccountID,
			Currency:              currency,
			AmountDelta:           amount,
			BookedAt:              bookedAt,
		},
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	debitEvent, ok, err := ledgeraggregate.NewLedgerAccountEventFromPosting(debitAccountID, debitPosting)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if !ok {
		return nil, stackErr.Error(fmt.Errorf("unsupported debit ledger event reference_type=%s", debitEventName))
	}
	creditEvent, ok, err := ledgeraggregate.NewLedgerAccountEventFromPosting(creditAccountID, creditPosting)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if !ok {
		return nil, stackErr.Error(fmt.Errorf("unsupported credit ledger event reference_type=%s", creditEventName))
	}

	return []eventpkg.Event{debitEvent, creditEvent}, nil
}

func paymentLedgerEventsFromPostings(inputs []paymentEventLedgerPostingInput) ([]eventpkg.Event, error) {
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

func newPaymentPosting(
	accountID string,
	transactionID string,
	referenceType string,
	referenceID string,
	counterpartyAccountID string,
	currency string,
	amountDelta int64,
	bookedAt time.Time,
) ledgerentity.LedgerAccountPosting {
	posting, err := ledgeraggregate.NewLedgerAccountPaymentPosting(valueobject.LedgerAccountPostingInput{
		AccountID:             accountID,
		TransactionID:         transactionID,
		ReferenceType:         referenceType,
		ReferenceID:           referenceID,
		CounterpartyAccountID: counterpartyAccountID,
		Currency:              currency,
		AmountDelta:           amountDelta,
		BookedAt:              bookedAt,
	})
	if err != nil {
		panic(err)
	}
	return posting
}

func debitLedgerEventNameForReversal(paymentEventName string) string {
	switch strings.TrimSpace(paymentEventName) {
	case sharedevents.EventPaymentRefunded:
		return ledgeraggregate.EventNameLedgerAccountWithdrawFromRefund
	case sharedevents.EventPaymentChargeback:
		return ledgeraggregate.EventNameLedgerAccountWithdrawFromChargeback
	default:
		return ""
	}
}

func creditLedgerEventNameForReversal(paymentEventName string) string {
	switch strings.TrimSpace(paymentEventName) {
	case sharedevents.EventPaymentRefunded:
		return ledgeraggregate.EventNameLedgerAccountDepositFromRefund
	case sharedevents.EventPaymentChargeback:
		return ledgeraggregate.EventNameLedgerAccountDepositFromChargeback
	default:
		return ""
	}
}

func ledgerClearingAccountID(clearingAccountKey string) string {
	return fmt.Sprintf("ledger:clearing:%s", strings.ToLower(strings.TrimSpace(clearingAccountKey)))
}

func providerClearingAccountKey(provider string) string {
	provider = strings.ToLower(strings.TrimSpace(provider))
	if provider == "" {
		return ""
	}
	return fmt.Sprintf("provider:%s", provider)
}

func reversalSuffix(reversalType string) string {
	switch strings.TrimSpace(reversalType) {
	case sharedevents.EventPaymentRefunded:
		return "refunded"
	case sharedevents.EventPaymentChargeback:
		return "chargeback"
	default:
		return "reversed"
	}
}

func resolvePaymentEventID(paymentID, transactionID string) string {
	return utils.FirstNonEmpty(strings.TrimSpace(paymentID), strings.TrimSpace(transactionID))
}
