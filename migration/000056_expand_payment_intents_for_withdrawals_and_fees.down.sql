DROP INDEX IF EXISTS idx_payment_intents_workflow_status_created_at;

UPDATE payment_intents
SET credit_account_id = COALESCE(credit_account_id, debit_account_id, 'legacy-payment-account');

ALTER TABLE payment_intents
    ALTER COLUMN credit_account_id SET NOT NULL,
    DROP COLUMN IF EXISTS destination_account_id,
    DROP COLUMN IF EXISTS debit_account_id,
    DROP COLUMN IF EXISTS provider_amount,
    DROP COLUMN IF EXISTS fee_amount,
    DROP COLUMN IF EXISTS workflow;
