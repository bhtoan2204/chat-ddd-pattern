package eventstore

var insertSnapshotSQL = `INSERT INTO payment_balance_snapshots (id, aggregate_id, version, state, created_at) VALUES (?, ?, ?, ?, ?)`

var readSnapshotSQL = `SELECT * FROM payment_events WHERE aggregate_id = ? AND version >= ? ORDER BY version DESC`

var updateVersionSQL = `UPDATE payment_aggregates SET version = ? WHERE aggregate_id = ? AND version = ?`
