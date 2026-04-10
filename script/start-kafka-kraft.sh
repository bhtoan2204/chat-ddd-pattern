#!/usr/bin/env bash

set -euo pipefail

cat >/tmp/kraft.properties <<EOF
process.roles=${KAFKA_PROCESS_ROLES}
node.id=${KAFKA_NODE_ID}
controller.quorum.voters=${KAFKA_CONTROLLER_QUORUM_VOTERS}
listeners=${KAFKA_LISTENERS}
advertised.listeners=${KAFKA_ADVERTISED_LISTENERS}
listener.security.protocol.map=${KAFKA_LISTENER_SECURITY_PROTOCOL_MAP}
inter.broker.listener.name=${KAFKA_INTER_BROKER_LISTENER_NAME}
controller.listener.names=${KAFKA_CONTROLLER_LISTENER_NAMES}
log.dirs=${KAFKA_LOG_DIRS}
num.partitions=1
offsets.topic.replication.factor=${KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR}
transaction.state.log.replication.factor=${KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR}
transaction.state.log.min.isr=${KAFKA_TRANSACTION_STATE_LOG_MIN_ISR}
group.initial.rebalance.delay.ms=${KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS}
confluent.license.topic.replication.factor=1
confluent.metadata.topic.replication.factor=1
confluent.security.event.logger.exporter.kafka.topic.replicas=1
EOF

mkdir -p "${KAFKA_LOG_DIRS}"

if [ ! -f "${KAFKA_LOG_DIRS}/meta.properties" ]; then
  kafka-storage format --ignore-formatted -t "${CLUSTER_ID}" -c /tmp/kraft.properties
fi

exec kafka-server-start /tmp/kraft.properties
