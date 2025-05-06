#!/bin/bash

# Create DLQ topics for each service
echo "Creating DLQ topics..."

# Order Service DLQ
docker-compose exec kafka kafka-topics --create --if-not-exists \
  --bootstrap-server kafka:9092 \
  --replication-factor 1 \
  --partitions 8 \
  --topic order-service-dlq \
  --config retention.ms=259200000

# Inventory Service DLQ
docker-compose exec kafka kafka-topics --create --if-not-exists \
  --bootstrap-server kafka:9092 \
  --replication-factor 1 \
  --partitions 8 \
  --topic inventory-service-dlq \
  --config retention.ms=259200000

# Payment Service DLQ
docker-compose exec kafka kafka-topics --create --if-not-exists \
  --bootstrap-server kafka:9092 \
  --replication-factor 1 \
  --partitions 8 \
  --topic payment-service-dlq \
  --config retention.ms=259200000

# Shipping Service DLQ
docker-compose exec kafka kafka-topics --create --if-not-exists \
  --bootstrap-server kafka:9092 \
  --replication-factor 1 \
  --partitions 8 \
  --topic shipping-service-dlq \
  --config retention.ms=259200000

# Notification Service DLQ
docker-compose exec kafka kafka-topics --create --if-not-exists \
  --bootstrap-server kafka:9092 \
  --replication-factor 1 \
  --partitions 8 \
  --topic notification-service-dlq \
  --config retention.ms=259200000

echo "DLQ topics created successfully!"
