CREATE UNIQUE INDEX idx_acc_outbox_aggregate_version
    ON account_outbox_events(aggregate_id, version);
