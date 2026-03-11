-- =========================
-- DROP: payment_transactions
-- =========================
DROP TABLE payment_transactions CASCADE CONSTRAINTS;

-- =========================
-- DROP: payment_outbox_events
-- =========================
DROP TABLE payment_outbox_events CASCADE CONSTRAINTS;

-- =========================
-- DROP: payment_event_offsets
-- =========================
DROP TABLE payment_event_offsets CASCADE CONSTRAINTS;

-- =========================
-- DROP: payment_events
-- =========================
DROP TABLE payment_events CASCADE CONSTRAINTS;

-- =========================
-- DROP: payment_balance_snapshots
-- =========================
DROP TABLE payment_balance_snapshots CASCADE CONSTRAINTS;

-- =========================
-- DROP: payment_balances
-- =========================
DROP TABLE payment_balances CASCADE CONSTRAINTS;
