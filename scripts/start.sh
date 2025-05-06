#!/bin/bash

# Script to start the entire system
echo "Khởi động hệ thống..."

# Không tạo thư mục logs và backups

# Đảm bảo các script có quyền thực thi
chmod +x scripts/*.sh

# Dừng các container đang chạy (nếu có)
echo "Dừng các container đang chạy (nếu có)..."
docker-compose down

# Xóa các volume cũ (tùy chọn, bỏ comment nếu muốn xóa dữ liệu cũ)
# echo "Xóa các volume cũ..."
# docker volume rm $(docker volume ls -q | grep mid-project) 2>/dev/null || true

# Build và khởi động các dịch vụ
echo "Build và khởi động các dịch vụ..."
docker-compose build
docker-compose up -d

# Đợi PostgreSQL khởi động
echo "Đợi PostgreSQL khởi động..."
sleep 10

# Đợi Kafka khởi động
echo "Đợi Kafka khởi động..."
sleep 10

# Tạo các Kafka topics cần thiết
echo "Tạo các Kafka topics cần thiết..."
docker-compose exec -T kafka kafka-topics --create --if-not-exists \
  --bootstrap-server kafka:9092 \
  --replication-factor 1 \
  --partitions 8 \
  --topic orders

docker-compose exec -T kafka kafka-topics --create --if-not-exists \
  --bootstrap-server kafka:9092 \
  --replication-factor 1 \
  --partitions 8 \
  --topic payments

docker-compose exec -T kafka kafka-topics --create --if-not-exists \
  --bootstrap-server kafka:9092 \
  --replication-factor 1 \
  --partitions 8 \
  --topic shipments

# Tạo các DLQ topics
echo "Tạo các DLQ topics..."
./scripts/create-dlq-topics.sh

# Kiểm tra trạng thái các dịch vụ
echo "Kiểm tra trạng thái các dịch vụ..."
docker-compose ps

echo "Hệ thống đã được khởi động thành công!"
echo "API Gateway: http://localhost:9090"
echo "Kafka UI: http://localhost:8090"
