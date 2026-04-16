DROP INDEX idx_ledger_aggregates_aggregate_id;
CREATE UNIQUE INDEX idx_ledger_aggregates_agg_type_id
    ON ledger_aggregates(aggregate_id, aggregate_type);

DROP INDEX idx_ledger_events_agg_ver;
CREATE UNIQUE INDEX idx_ledger_events_agg_type_ver
    ON ledger_events(aggregate_id, aggregate_type, version);
