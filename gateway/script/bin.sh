#!/bin/bash

SCRIPTPATH="$(
    cd "$(dirname "$0")"
    pwd -P
)"
CURRENT_DIR=$SCRIPTPATH
ROOT_DIR="$(dirname $CURRENT_DIR)"

function applyEnv() {
  set -a;
  export $(grep -v '^#' $ROOT_DIR/secret/.env | xargs -0) >/dev/null 2>&1; . $ROOT_DIR/secret/.env;
  set +a;
}

function run() {
  applyEnv
  echo "Running server..."
  go run $ROOT_DIR/cmd/main.go
}

case "$1" in
  run)
    run
    ;;
  *)
    echo "Usage: $0 {run}"
    exit 1
    ;;
esac