DROP INDEX idx_ledger_events_agg_type_ver;

CREATE UNIQUE INDEX idx_ledger_events_agg_ver
    ON ledger_events(aggregate_id, version);

DROP INDEX idx_ledger_aggregates_agg_type_id;

CREATE UNIQUE INDEX idx_ledger_aggregates_aggregate_id
    ON ledger_aggregates(aggregate_id);
