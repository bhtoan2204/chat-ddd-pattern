CREATE TABLE ledger_posted_transactions (
    id                      VARCHAR2(1024) PRIMARY KEY,
    aggregate_id            VARCHAR2(1024) NOT NULL,
    aggregate_type          VARCHAR2(255)  NOT NULL,
    transaction_id          VARCHAR2(1024) NOT NULL,
    reference_type          VARCHAR2(255)  NOT NULL,
    reference_id            VARCHAR2(1024) NOT NULL,
    counterparty_account_id VARCHAR2(1024) NOT NULL,
    currency                VARCHAR2(16)   NOT NULL,
    amount_delta            NUMBER(19)     NOT NULL,
    booked_at               TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL
);

INSERT INTO ledger_posted_transactions (
    id,
    aggregate_id,
    aggregate_type,
    transaction_id,
    reference_type,
    reference_id,
    counterparty_account_id,
    currency,
    amount_delta,
    booked_at,
    created_at
)
SELECT
    LOWER(RAWTOHEX(SYS_GUID())) AS id,
    aggregate_id,
    aggregate_type,
    JSON_VALUE(event_data, '$.transaction_id') AS transaction_id,
    CASE event_name
        WHEN 'EventLedgerAccountPaymentBooked' THEN NVL(JSON_VALUE(event_data, '$.reference_type'), 'payment.succeeded')
        ELSE 'ledger.transfer_to_account'
    END AS reference_type,
    CASE event_name
        WHEN 'EventLedgerAccountPaymentBooked' THEN JSON_VALUE(event_data, '$.payment_id')
        ELSE JSON_VALUE(event_data, '$.transaction_id')
    END AS reference_id,
    CASE event_name
        WHEN 'EventLedgerAccountPaymentBooked' THEN JSON_VALUE(event_data, '$.counterparty_account_id')
        WHEN 'EventLedgerAccountTransferredToAccount' THEN JSON_VALUE(event_data, '$.to_account_id')
        WHEN 'EventLedgerAccountReceivedTransfer' THEN JSON_VALUE(event_data, '$.from_account_id')
    END AS counterparty_account_id,
    UPPER(JSON_VALUE(event_data, '$.currency')) AS currency,
    CASE event_name
        WHEN 'EventLedgerAccountPaymentBooked' THEN JSON_VALUE(event_data, '$.amount_delta' RETURNING NUMBER)
        WHEN 'EventLedgerAccountTransferredToAccount' THEN -JSON_VALUE(event_data, '$.amount' RETURNING NUMBER)
        WHEN 'EventLedgerAccountReceivedTransfer' THEN JSON_VALUE(event_data, '$.amount' RETURNING NUMBER)
    END AS amount_delta,
    created_at AS booked_at,
    created_at
FROM ledger_events
WHERE aggregate_type = 'LedgerAccountAggregate'
  AND event_name IN (
      'EventLedgerAccountPaymentBooked',
      'EventLedgerAccountTransferredToAccount',
      'EventLedgerAccountReceivedTransfer'
  );

CREATE UNIQUE INDEX idx_ledger_posted_tx_agg_type_tx
    ON ledger_posted_transactions(aggregate_id, aggregate_type, transaction_id);
