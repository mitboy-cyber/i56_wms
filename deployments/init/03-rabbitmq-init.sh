#!/bin/sh
# I56 Framework — RabbitMQ Init Script
# Creates required queues and exchanges

echo "[rabbitmq-init] Waiting for RabbitMQ to be ready..."

# Wait for RabbitMQ to be ready
until rabbitmq-diagnostics check_port_connectivity 2>/dev/null; do
  sleep 2
done

echo "[rabbitmq-init] RabbitMQ is ready. Setting up queues..."

# Declare exchanges
rabbitmqadmin declare exchange name=i56.events type=topic durable=true

# Declare queues
rabbitmqadmin declare queue name=i56.parcel.events durable=true
rabbitmqadmin declare queue name=i56.order.events durable=true
rabbitmqadmin declare queue name=i56.notification.events durable=true
rabbitmqadmin declare queue name=i56.webhook.events durable=true
rabbitmqadmin declare queue name=i56.audit.events durable=true
rabbitmqadmin declare queue name=i56.dead_letter durable=true

# Bind queues to exchange with routing keys
rabbitmqadmin declare binding source=i56.events destination=i56.parcel.events routing_key="parcel.*"
rabbitmqadmin declare binding source=i56.events destination=i56.order.events routing_key="order.*"
rabbitmqadmin declare binding source=i56.events destination=i56.notification.events routing_key="notification.*"
rabbitmqadmin declare binding source=i56.events destination=i56.webhook.events routing_key="webhook.*"
rabbitmqadmin declare binding source=i56.events destination=i56.audit.events routing_key="audit.*"

# Set dead letter policy
rabbitmqadmin declare policy name=DLX pattern="^i56\." '{"dead-letter-exchange":"i56.events","dead-letter-routing-key":"dead_letter"}' --priority 1

echo "[rabbitmq-init] Queues created: parcel, order, notification, webhook, audit events"
echo "[rabbitmq-init] RabbitMQ initialization complete."
