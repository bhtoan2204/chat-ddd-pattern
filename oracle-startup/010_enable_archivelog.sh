#!/bin/bash

set -euo pipefail

archive_dest="${ORACLE_ARCHIVE_LOG_DEST:-/opt/oracle/oradata/FREE/archivelog}"
mkdir -p "$archive_dest"

current_log_mode="$(
  sqlplus -s / as sysdba <<'SQL'
whenever sqlerror exit failure
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
else
  echo "Enabling ARCHIVELOG mode"

  sqlplus -s / as sysdba <<'SQL'
whenever sqlerror exit failure
shutdown immediate;
startup mount;
alter database archivelog;
alter database open;
exit;
SQL
fi

sqlplus -s / as sysdba <<SQL
whenever sqlerror exit failure
alter system set log_archive_dest_1='LOCATION=${archive_dest}' scope=both;
archive log list;
exit;
SQL
