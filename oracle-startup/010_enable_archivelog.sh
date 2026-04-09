#!/bin/bash

set -euo pipefail

current_log_mode="$(
  sqlplus -s / as sysdba <<'SQL'
set heading off
set feedback off
set verify off
set pagesize 0
select log_mode from v$database;
exit;
SQL
)"

current_log_mode="$(printf '%s' "$current_log_mode" | tr -d '[:space:]')"

if [ "$current_log_mode" = "ARCHIVELOG" ]; then
  echo "ARCHIVELOG already enabled"
  exit 0
fi

echo "Enabling ARCHIVELOG mode"

sqlplus -s / as sysdba <<'SQL'
shutdown immediate;
startup mount;
alter database archivelog;
alter database open;
archive log list;
exit;
SQL
