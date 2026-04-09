#!/bin/sh

set -eu

CONNECT_URL="${CONNECT_URL:-http://connect:8083}"
CONNECTOR_CONFIG_FILE="${CONNECTOR_CONFIG_FILE:-/connector/connector_config.json}"
CONNECTOR_NAME="${CONNECTOR_NAME:-}"
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

split_connectors() {
  config_file="$1"
  output_dir="$2"

  awk -v output_dir="$output_dir" '
    BEGIN {
      depth = 0
      item_index = 0
      object = ""
    }
    {
      line = $0 ORS
      for (i = 1; i <= length(line); i++) {
        ch = substr(line, i, 1)

        if (ch == "{") {
          if (depth == 0) {
            object = ""
          }
          depth++
        }

        if (depth > 0) {
          object = object ch
        }

        if (ch == "}") {
          depth--

          if (depth == 0) {
            item_index++
            file = sprintf("%s/connector_%03d.json", output_dir, item_index)
            print object > file
            close(file)
            object = ""
          }
        }
      }
    }
    END {
      if (depth != 0) {
        exit 1
      }
    }
  ' "$config_file"
}

extract_connector_name() {
  config_file="$1"
  sed -n 's/^[[:space:]]*"name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$config_file" | head -n 1
}

extract_connector_config() {
  config_file="$1"
  sed -n '/"config"[[:space:]]*:/,$p' "$config_file" \
    | sed '1s/.*"config"[[:space:]]*:[[:space:]]*//' \
    | sed '$d'
}

upsert_connector() {
  config_file="$1"
  connector_name="$(extract_connector_name "$config_file")"

  if [ -z "$connector_name" ]; then
    echo "Skipping $(basename "$config_file"): missing connector name" >&2
    return 1
  fi

  attempt=1

  while [ "$attempt" -le "$MAX_RETRIES" ]; do
    body_file="$(mktemp)"
    payload_file="$(mktemp)"
    extract_connector_config "$config_file" >"$payload_file"

    status_code="$(
      curl -sS -o "$body_file" -w "%{http_code}" \
        -X PUT \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        --data "@$payload_file" \
        "$CONNECT_URL/connectors/$connector_name/config"
    )"

    case "$status_code" in
      200)
        echo "Updated connector $connector_name from $(basename "$config_file")"
        rm -f "$body_file" "$payload_file"
        return 0
        ;;
      201)
        echo "Created connector $connector_name from $(basename "$config_file")"
        rm -f "$body_file" "$payload_file"
        return 0
        ;;
    esac

    echo "[$attempt/$MAX_RETRIES] connector upsert failed for $connector_name with HTTP $status_code"
    cat "$body_file"
    echo
    rm -f "$body_file" "$payload_file"
    sleep "$SLEEP_SECONDS"
    attempt=$((attempt + 1))
  done

  echo "Failed to upsert connector: $connector_name" >&2
  return 1
}

main() {
  if [ ! -f "$CONNECTOR_CONFIG_FILE" ]; then
    echo "Connector config file not found: $CONNECTOR_CONFIG_FILE" >&2
    exit 1
  fi

  wait_for_connect

  temp_dir="$(mktemp -d)"
  trap 'rm -rf "$temp_dir"' EXIT INT TERM HUP

  split_connectors "$CONNECTOR_CONFIG_FILE" "$temp_dir"

  matched=0
  for config_file in "$temp_dir"/connector_*.json; do
    [ -f "$config_file" ] || continue

    current_name="$(extract_connector_name "$config_file")"
    if [ -n "$CONNECTOR_NAME" ] && [ "$current_name" != "$CONNECTOR_NAME" ]; then
      continue
    fi

    matched=$((matched + 1))
    upsert_connector "$config_file"
  done

  if [ "$matched" -eq 0 ]; then
    if [ -n "$CONNECTOR_NAME" ]; then
      echo "Connector not found in $CONNECTOR_CONFIG_FILE: $CONNECTOR_NAME" >&2
    else
      echo "No connector definitions found in $CONNECTOR_CONFIG_FILE" >&2
    fi
    exit 1
  fi
}

main "$@"
