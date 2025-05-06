#!/bin/bash

# Script to test the system
echo "Kiểm tra hệ thống..."

# Kiểm tra API Gateway
echo "Kiểm tra API Gateway..."
curl -s http://localhost:9090/health | grep -q "OK" && echo "API Gateway: OK" || echo "API Gateway: FAILED"

# Kiểm tra Order Service
echo "Kiểm tra Order Service..."
curl -s http://localhost:8081/health | grep -q "UP" && echo "Order Service: OK" || echo "Order Service: FAILED"

# Kiểm tra Inventory Service
echo "Kiểm tra Inventory Service..."
curl -s http://localhost:8082/health | grep -q "UP" && echo "Inventory Service: OK" || echo "Inventory Service: FAILED"

# Kiểm tra Payment Service
echo "Kiểm tra Payment Service..."
curl -s http://localhost:8083/health | grep -q "UP" && echo "Payment Service: OK" || echo "Payment Service: FAILED"

# Kiểm tra Shipping Service
echo "Kiểm tra Shipping Service..."
curl -s http://localhost:8084/health | grep -q "UP" && echo "Shipping Service: OK" || echo "Shipping Service: FAILED"

# Kiểm tra Notification Service
echo "Kiểm tra Notification Service..."
curl -s http://localhost:8085/health | grep -q "UP" && echo "Notification Service: OK" || echo "Notification Service: FAILED"

# Kiểm tra User Service
echo "Kiểm tra User Service..."
curl -s http://localhost:8086/health | grep -q "UP" && echo "User Service: OK" || echo "User Service: FAILED"

# Kiểm tra Cart Service
echo "Kiểm tra Cart Service..."
curl -s http://localhost:8087/health | grep -q "UP" && echo "Cart Service: OK" || echo "Cart Service: FAILED"

# Kiểm tra Kafka
echo "Kiểm tra Kafka..."
docker-compose exec -T kafka kafka-topics --list --bootstrap-server kafka:9092 | grep -q "orders" && echo "Kafka: OK" || echo "Kafka: FAILED"

# Kiểm tra PostgreSQL
echo "Kiểm tra PostgreSQL..."
docker-compose exec -T postgres pg_isready -U postgres | grep -q "accepting connections" && echo "PostgreSQL: OK" || echo "PostgreSQL: FAILED"

# Kiểm tra Redis
echo "Kiểm tra Redis..."
docker-compose exec -T redis redis-cli ping | grep -q "PONG" && echo "Redis: OK" || echo "Redis: FAILED"

echo "Kiểm tra hoàn tất!"
