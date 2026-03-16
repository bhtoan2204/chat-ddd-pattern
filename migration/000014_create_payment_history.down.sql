-- =========================
-- DROP: payment_histories
-- =========================
DROP INDEX idx_payment_histories_receiver_created;
DROP INDEX idx_payment_histories_sender_created;
DROP INDEX idx_payment_histories_receiver_id;
DROP INDEX idx_payment_histories_sender_id;
DROP INDEX idx_payment_histories_type;

DROP TABLE payment_histories;