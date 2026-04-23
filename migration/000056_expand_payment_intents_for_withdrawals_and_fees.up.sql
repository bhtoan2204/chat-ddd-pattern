ALTER TABLE payment_intents
    ADD COLUMN IF NOT EXISTS workflow VARCHAR(64),
    ADD COLUMN IF NOT EXISTS fee_amount BIGINT DEFAULT 0 NOT NULL,
    ADD COLUMN IF NOT EXISTS provider_amount BIGINT DEFAULT 0 NOT NULL,
    ADD COLUMN IF NOT EXISTS debit_account_id VARCHAR(1024),
    ADD COLUMN IF NOT EXISTS destination_account_id VARCHAR(1024);

UPDATE payment_intents
SET workflow = COALESCE(NULLIF(workflow, ''), 'TOP_UP'),
    fee_amount = COALESCE(fee_amount, 0),
    provider_amount = CASE
        WHEN provider_amount IS NULL OR provider_amount = 0 THEN amount + COALESCE(fee_amount, 0)
        ELSE provider_amount
    END
WHERE workflow IS NULL
   OR workflow = ''
   OR provider_amount IS NULL
   OR provider_amount = 0;

ALTER TABLE payment_intents
    ALTER COLUMN workflow SET NOT NULL,
    ALTER COLUMN credit_account_id DROP NOT NULL;

CREATE INDEX IF NOT EXISTS idx_payment_intents_workflow_status_created_at
    ON payment_intents(workflow, status, created_at);
