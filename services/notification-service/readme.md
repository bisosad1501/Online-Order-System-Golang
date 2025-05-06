# Notification Service

## Overview
Notification Service quản lý việc gửi thông báo cho khách hàng. Service này xử lý việc tạo thông báo mới, gửi thông báo qua email, SMS, hoặc push notification, và lưu trữ thông tin thông báo.

## Cấu trúc thư mục
```
notification-service/
├── Dockerfile                # Cấu hình Docker
├── README.md                 # Tài liệu hướng dẫn
└── src/                      # Mã nguồn
    ├── main.go               # Điểm khởi đầu của ứng dụng
    ├── config/               # Cấu hình ứng dụng
    │   └── config.go         # Đọc cấu hình từ biến môi trường
    ├── api/                  # API handlers
    │   ├── handlers.go       # Các handler xử lý request
    │   ├── middleware.go     # Middleware (logging, CORS)
    │   └── routes.go         # Định nghĩa routes
    ├── models/               # Định nghĩa các model
    │   └── models.go         # Các struct đại diện cho dữ liệu
    ├── interfaces/           # Định nghĩa các interface
    │   └── interfaces.go     # Các interface cho service và producer
    ├── db/                   # Tương tác với database
    │   ├── db.go             # Kết nối database
    │   └── repository.go     # Các hàm truy vấn database
    ├── kafka/                # Tương tác với Kafka
    │   ├── consumer.go       # Kafka consumer
    │   └── producer.go       # Kafka producer
    ├── service/              # Business logic
    │   ├── service.go        # Các hàm xử lý logic nghiệp vụ
    │   ├── email_sender.go   # Gửi email
    │   ├── sms_sender.go     # Gửi SMS
    │   └── push_sender.go    # Gửi push notification
    └── utils/                # Các tiện ích
        └── utils.go          # Các hàm tiện ích
```

## Cài đặt và Chạy

### Yêu cầu
- Go 1.20 trở lên
- PostgreSQL
- Kafka

### Biến môi trường
- `PORT`: Port để chạy service (mặc định: 8085)
- `DB_HOST`: Host của PostgreSQL (mặc định: localhost)
- `DB_PORT`: Port của PostgreSQL (mặc định: 5432)
- `DB_USER`: Username của PostgreSQL (mặc định: postgres)
- `DB_PASSWORD`: Password của PostgreSQL (mặc định: postgres)
- `DB_NAME`: Tên database (mặc định: notificationdb)
- `KAFKA_BOOTSTRAP_SERVERS`: Kafka bootstrap servers (mặc định: kafka:29092)
- `KAFKA_TOPIC`: Kafka topic (mặc định: notifications)
- `SMTP_HOST`: Host của SMTP server (mặc định: smtp.example.com)
- `SMTP_PORT`: Port của SMTP server (mặc định: 587)
- `SMTP_USERNAME`: Username của SMTP server (mặc định: user@example.com)
- `SMTP_PASSWORD`: Password của SMTP server (mặc định: password)
- `SMTP_FROM`: Email gửi từ (mặc định: noreply@example.com)
- `SMS_API_KEY`: API key của SMS provider (mặc định: sms_api_key)
- `SMS_API_SECRET`: API secret của SMS provider (mặc định: sms_api_secret)
- `SMS_FROM_NUMBER`: Số điện thoại gửi từ (mặc định: +1234567890)
- `PUSH_API_KEY`: API key của push notification provider (mặc định: push_api_key)

### Chạy với Docker
```bash
docker build -t notification-service .
docker run -p 8085:8085 notification-service
```

### Chạy với Docker Compose
```bash
docker-compose up -d
```

## API Endpoints

### Health Check
- `GET /health`: Kiểm tra trạng thái của service

### Notifications
- `POST /notifications`: Tạo thông báo mới
- `GET /notifications`: Lấy danh sách thông báo
- `GET /notifications/{id}`: Lấy thông tin thông báo theo ID
- `GET /notifications/customer/{customer_id}`: Lấy thông tin thông báo theo customer ID
- `PUT /notifications/{id}/status`: Cập nhật trạng thái thông báo

## Database Schema

### Notifications Table
```sql
CREATE TABLE IF NOT EXISTS notifications (
    id VARCHAR(36) PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    sent_at TIMESTAMP
)
```

## Kafka Events

### Produces
- `notification_created`: Khi thông báo được tạo
- `notification_sent`: Khi thông báo được gửi thành công
- `notification_failed`: Khi thông báo gửi thất bại

### Consumes
- `order_created`: Để tạo thông báo khi đơn hàng được tạo
- `order_confirmed`: Để tạo thông báo khi đơn hàng được xác nhận
- `order_cancelled`: Để tạo thông báo khi đơn hàng bị hủy
- `order_completed`: Để tạo thông báo khi đơn hàng hoàn thành
- `payment_successful`: Để tạo thông báo khi thanh toán thành công
- `payment_failed`: Để tạo thông báo khi thanh toán thất bại
- `shipment_created`: Để tạo thông báo khi lô hàng được tạo
- `shipping_completed`: Để tạo thông báo khi lô hàng đã được giao

## Luồng xử lý thông báo

1. **Tạo thông báo**:
   - Service nhận event từ Kafka hoặc request từ API
   - Service tạo thông báo mới với trạng thái `PENDING`
   - Service lưu thông báo vào database
   - Service gửi event `notification_created` đến Kafka

2. **Gửi thông báo**:
   - Service gửi thông báo qua email, SMS, hoặc push notification
   - Nếu gửi thành công, service cập nhật trạng thái thông báo thành `SENT`
   - Nếu gửi thất bại, service cập nhật trạng thái thông báo thành `FAILED`
   - Service gửi event `notification_sent` hoặc `notification_failed` đến Kafka

## Xử lý lỗi

- **Thông báo không tồn tại**: Trả về lỗi 404 Not Found
- **Trạng thái không hợp lệ**: Trả về lỗi 400 Bad Request
- **Lỗi database**: Trả về lỗi 500 Internal Server Error
- **Lỗi Kafka**: Log lỗi và tiếp tục xử lý
- **Lỗi gửi thông báo**: Log lỗi và cập nhật trạng thái thông báo thành `FAILED`

## Bảo mật

- CORS middleware cho phép cross-origin requests
- Validation đầu vào để ngăn chặn các cuộc tấn công như SQL Injection

## Logging

- Logger middleware ghi lại tất cả các request với thông tin như method, path, latency, status code, và IP

## Monitoring

- Health check endpoint để kiểm tra trạng thái của service

## Phát triển

### Thêm API mới
1. Định nghĩa model trong `models/models.go`
2. Thêm hàm xử lý trong `service/service.go`
3. Thêm handler trong `api/handlers.go`
4. Thêm route trong `api/routes.go`

### Thêm Kafka event mới
1. Định nghĩa event trong `models/models.go`
2. Thêm hàm publish trong `kafka/producer.go`
3. Thêm xử lý event trong `kafka/consumer.go`
