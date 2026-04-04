#!/usr/bin/env bash
set -euo pipefail

CONNECT_URL="${CONNECT_URL:-http://localhost:8083}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_PATH="${SCRIPT_DIR}/payment-outbox-connector.json"

curl -sS -X PUT \
  -H "Content-Type: application/json" \
  --data @"${CONFIG_PATH}" \
  "${CONNECT_URL}/connectors/payment-outbox-connector/config"
