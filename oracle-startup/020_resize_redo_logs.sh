#!/bin/bash

set -euo pipefail

target_mb="${ORACLE_REDO_LOG_SIZE_MB:-1024}"
target_mb="$(printf '%s' "$target_mb" | tr -d '[:space:]')"
target_group_count="${ORACLE_REDO_LOG_GROUP_COUNT:-4}"
target_group_count="$(printf '%s' "$target_group_count" | tr -d '[:space:]')"

case "$target_mb" in
  ''|*[!0-9]*)
    echo "ORACLE_REDO_LOG_SIZE_MB must be a positive integer in MB" >&2
    exit 1
    ;;
esac

if [ "$target_mb" -lt 512 ]; then
  echo "ORACLE_REDO_LOG_SIZE_MB must be at least 512MB for Debezium LogMiner stability" >&2
  exit 1
fi

case "$target_group_count" in
  ''|*[!0-9]*)
    echo "ORACLE_REDO_LOG_GROUP_COUNT must be a positive integer" >&2
    exit 1
    ;;
esac

if [ "$target_group_count" -lt 4 ]; then
  echo "ORACLE_REDO_LOG_GROUP_COUNT must be at least 4 for Debezium LogMiner stability" >&2
  exit 1
fi

sql_query() {
  local sql="$1"

  sqlplus -s / as sysdba <<SQL
whenever sqlerror exit failure
set heading off
set feedback off
set verify off
set pagesize 0
$sql
exit;
SQL
}

sql_scalar() {
  local sql="$1"

  sql_query "$sql" | tr -d '[:space:]'
}

log_member_dir="$(
  sql_query "select regexp_replace(member, '/[^/]+$', '') from v\$logfile where rownum = 1;"
)"
log_member_dir="$(printf '%s' "$log_member_dir" | tr -d '[:space:]')"
if [ -z "$log_member_dir" ]; then
  echo "Unable to determine redo log directory from v\$logfile" >&2
  exit 1
fi

small_group_count="$(sql_scalar "select count(*) from v\$log where bytes < (${target_mb} * 1024 * 1024);")"
total_group_count="$(sql_scalar "select count(*) from v\$log;")"
large_group_count="$(sql_scalar "select count(*) from v\$log where bytes >= (${target_mb} * 1024 * 1024);")"
desired_group_count="$total_group_count"
if [ "$desired_group_count" -lt "$target_group_count" ]; then
  desired_group_count="$target_group_count"
fi

if [ "${small_group_count:-0}" -eq 0 ] && [ "${large_group_count:-0}" -ge "$desired_group_count" ]; then
  echo "Redo log groups are already at least ${target_mb}MB and count >= ${desired_group_count}"
  exit 0
fi

additional_groups_needed=$((desired_group_count - large_group_count))
if [ "$additional_groups_needed" -lt 0 ]; then
  additional_groups_needed=0
fi

if [ "$additional_groups_needed" -gt 0 ]; then
  echo "Adding ${additional_groups_needed} redo log groups sized ${target_mb}MB to reach ${desired_group_count} durable groups"
  for _ in $(seq 1 "$additional_groups_needed"); do
    next_group="$(sql_scalar "select nvl(max(group#), 0) + 1 from v\$log;")"
    next_member="${log_member_dir}/redo$(printf '%02d' "$next_group").log"
    sql_query "alter database add logfile group ${next_group} ('${next_member}') size ${target_mb}M;"
  done
fi

max_switches=$((small_group_count * 4 + 4))
switch_attempt=1
while [ "$switch_attempt" -le "$max_switches" ]; do
  active_small_groups="$(sql_scalar "select count(*) from v\$log where bytes < (${target_mb} * 1024 * 1024) and status in ('CURRENT', 'ACTIVE');")"
  if [ "${active_small_groups:-0}" -eq 0 ]; then
    break
  fi

  echo "Cycling redo logs to retire small active groups (${switch_attempt}/${max_switches})"
  sql_query "alter system switch logfile;"
  sql_query "alter system checkpoint;"
  switch_attempt=$((switch_attempt + 1))
done

droppable_groups="$(
  sql_query "select group# from v\$log where bytes < (${target_mb} * 1024 * 1024) and status not in ('CURRENT', 'ACTIVE') order by group#;"
)"

while IFS= read -r group_id; do
  group_id="$(printf '%s' "$group_id" | tr -d '[:space:]')"
  if [ -z "$group_id" ]; then
    continue
  fi

  echo "Dropping redo log group ${group_id} because it is smaller than ${target_mb}MB"
  sql_query "alter database drop logfile group ${group_id};"
done <<< "$droppable_groups"

remaining_small_groups="$(
  sql_query "select group# || ':' || status || ':' || round(bytes / 1024 / 1024) || 'MB' from v\$log where bytes < (${target_mb} * 1024 * 1024) order by group#;"
)"

if [ -n "$(printf '%s' "$remaining_small_groups" | tr -d '[:space:]')" ]; then
  echo "Some redo log groups are still smaller than ${target_mb}MB:"
  printf '%s\n' "$remaining_small_groups"
else
  echo "All redo log groups are now at least ${target_mb}MB and the database has at least ${desired_group_count} groups"
fi

echo "Current redo log layout:"
sql_query "select group# || ':' || round(bytes / 1024 / 1024) || 'MB:' || status from v\$log order by group#;"
