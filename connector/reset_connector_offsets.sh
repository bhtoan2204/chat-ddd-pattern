#!/bin/sh

set -eu

CONNECT_URL="${CONNECT_URL:-http://localhost:8083}"
CONNECTOR_NAME="${CONNECTOR_NAME:-}"
MAX_RETRIES="${MAX_RETRIES:-60}"
SLEEP_SECONDS="${SLEEP_SECONDS:-2}"

if [ -z "$CONNECTOR_NAME" ]; then
  echo "CONNECTOR_NAME is required" >&2
  exit 1
fi

wait_for_connect() {
  attempt=1
  while [ "$attempt" -le "$MAX_RETRIES" ]; do
    if curl -fsS "$CONNECT_URL/" >/dev/null 2>&1; then
      return 0
    fi
    echo "[$attempt/$MAX_RETRIES] waiting for Kafka Connect at $CONNECT_URL"
    sleep "$SLEEP_SECONDS"
    attempt=$((attempt + 1))
  done

  echo "Kafka Connect did not become ready in time" >&2
  return 1
}

get_status() {
  curl -fsS "$CONNECT_URL/connectors/$CONNECTOR_NAME/status"
}

get_offsets() {
  curl -fsS "$CONNECT_URL/connectors/$CONNECTOR_NAME/offsets"
}

request() {
  method="$1"
  url="$2"
  body_file="$(mktemp)"

  status_code="$(
    curl -sS -o "$body_file" -w "%{http_code}" -X "$method" "$url"
  )"

  cat "$body_file"
  rm -f "$body_file"

  printf '\n__HTTP_STATUS__=%s\n' "$status_code"
}

status_is_stopped() {
  status_json="$1"

  case "$status_json" in
    *'"state":"STOPPED"'*)
      return 0
      ;;
  esac

  return 1
}

wait_for_stopped() {
  attempt=1
  while [ "$attempt" -le "$MAX_RETRIES" ]; do
    status_json="$(get_status)"
    if status_is_stopped "$status_json"; then
      printf '%s\n' "$status_json"
      return 0
    fi

    echo "[$attempt/$MAX_RETRIES] waiting for connector $CONNECTOR_NAME to stop"
    sleep "$SLEEP_SECONDS"
    attempt=$((attempt + 1))
  done

  echo "Connector $CONNECTOR_NAME did not reach STOPPED state in time" >&2
  return 1
}

extract_http_status() {
  printf '%s\n' "$1" | sed -n 's/^__HTTP_STATUS__=\([0-9][0-9][0-9]\)$/\1/p' | tail -n 1
}

print_response_body() {
  printf '%s\n' "$1" | sed '/^__HTTP_STATUS__=/d'
}

wait_for_running() {
  attempt=1
  while [ "$attempt" -le "$MAX_RETRIES" ]; do
    status_json="$(get_status)"

    case "$status_json" in
      *'"state":"RUNNING"'*)
        printf '%s\n' "$status_json"
        return 0
        ;;
    esac

    echo "[$attempt/$MAX_RETRIES] waiting for connector $CONNECTOR_NAME to resume"
    sleep "$SLEEP_SECONDS"
    attempt=$((attempt + 1))
  done

  echo "Connector $CONNECTOR_NAME did not return to RUNNING state in time" >&2
  return 1
}

wait_for_connect

echo "Current status:"
get_status
echo
echo "Current offsets:"
get_offsets
echo
echo "Resetting offsets for $CONNECTOR_NAME will cause the connector to rebuild state from scratch."
echo "If snapshot.mode=initial, downstream consumers may receive duplicate records."

echo "Stopping connector $CONNECTOR_NAME"
stop_response="$(request PUT "$CONNECT_URL/connectors/$CONNECTOR_NAME/stop")"
stop_status="$(extract_http_status "$stop_response")"
case "$stop_status" in
  200|202|204)
    ;;
  *)
    print_response_body "$stop_response" >&2
    echo "Stop request failed for $CONNECTOR_NAME (HTTP $stop_status)" >&2
    exit 1
    ;;
esac

echo "Waiting for connector to stop"
wait_for_stopped
echo

echo "Deleting stored offsets for $CONNECTOR_NAME"
reset_response="$(request DELETE "$CONNECT_URL/connectors/$CONNECTOR_NAME/offsets")"
reset_status="$(extract_http_status "$reset_response")"
case "$reset_status" in
  200|202|204)
    ;;
  *)
    print_response_body "$reset_response" >&2
    echo "Offset reset failed for $CONNECTOR_NAME (HTTP $reset_status)" >&2
    exit 1
    ;;
esac

echo "Resuming connector $CONNECTOR_NAME"
resume_response="$(request PUT "$CONNECT_URL/connectors/$CONNECTOR_NAME/resume")"
resume_status="$(extract_http_status "$resume_response")"
case "$resume_status" in
  200|202|204)
    ;;
  *)
    print_response_body "$resume_response" >&2
    echo "Resume request failed for $CONNECTOR_NAME (HTTP $resume_status)" >&2
    exit 1
    ;;
esac

echo "Waiting for connector to resume"
wait_for_running
echo

echo "Offsets after reset:"
get_offsets
