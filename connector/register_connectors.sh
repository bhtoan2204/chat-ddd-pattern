#!/bin/sh

set -eu

CONNECT_URL="${CONNECT_URL:-http://connect:8083}"
CONNECTOR_DIR="${CONNECTOR_DIR:-/connector}"
MAX_RETRIES="${MAX_RETRIES:-60}"
SLEEP_SECONDS="${SLEEP_SECONDS:-5}"

wait_for_connect() {
  attempt=1
  while [ "$attempt" -le "$MAX_RETRIES" ]; do
    if curl -fsS "$CONNECT_URL/connectors" >/dev/null 2>&1; then
      return 0
    fi
    echo "[$attempt/$MAX_RETRIES] waiting for Kafka Connect at $CONNECT_URL"
    sleep "$SLEEP_SECONDS"
    attempt=$((attempt + 1))
  done

  echo "Kafka Connect did not become ready in time" >&2
  return 1
}

register_connector() {
  config_file="$1"
  attempt=1

  while [ "$attempt" -le "$MAX_RETRIES" ]; do
    body_file="$(mktemp)"
    status_code="$(
      curl -sS -o "$body_file" -w "%{http_code}" \
        -X POST \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        --data "@$config_file" \
        "$CONNECT_URL/connectors"
    )"

    case "$status_code" in
      200|201|202|204|409)
        echo "Registered connector from $(basename "$config_file")"
        rm -f "$body_file"
        return 0
        ;;
    esac

    echo "[$attempt/$MAX_RETRIES] connector bootstrap failed for $(basename "$config_file") with HTTP $status_code"
    cat "$body_file"
    echo
    rm -f "$body_file"
    sleep "$SLEEP_SECONDS"
    attempt=$((attempt + 1))
  done

  echo "Failed to register connector from $(basename "$config_file")" >&2
  return 1
}

main() {
  wait_for_connect

  for config_file in "$CONNECTOR_DIR"/*_config.json; do
    [ -f "$config_file" ] || continue
    register_connector "$config_file"
  done
}

main "$@"
