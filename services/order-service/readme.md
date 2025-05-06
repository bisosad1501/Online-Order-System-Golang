# Payment Service

## Overview
Payment Service quản lý việc xử lý thanh toán cho các đơn hàng. Service này xử lý việc tạo thanh toán mới, cập nhật trạng thái thanh toán, và lưu trữ thông tin thanh toán.

## Cấu trúc thư mục
```
payment-service/
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
    │   └── service.go        # Các hàm xử lý logic nghiệp vụ
    └── utils/                # Các tiện ích
        └── utils.go          # Các hàm tiện ích
```

## Cài đặt và Chạy

### Yêu cầu
- Go 1.20 trở lên
- PostgreSQL
- Kafka

### Biến môi trường
- `PORT`: Port để chạy service (mặc định: 8083)
- `DB_HOST`: Host của PostgreSQL (mặc định: localhost)
- `DB_PORT`: Port của PostgreSQL (mặc định: 5432)
- `DB_USER`: Username của PostgreSQL (mặc định: postgres)
- `DB_PASSWORD`: Password của PostgreSQL (mặc định: postgres)
- `DB_NAME`: Tên database (mặc định: paymentdb)
- `KAFKA_BOOTSTRAP_SERVERS`: Kafka bootstrap servers (mặc định: localhost:9092)
- `KAFKA_TOPIC`: Kafka topic (mặc định: payments)
- `PAYMENT_GATEWAY_URL`: URL của payment gateway (mặc định: https://api.example.com/payments)
- `PAYMENT_GATEWAY_KEY`: API key của payment gateway (mặc định: test_key)

### Chạy với Docker
```bash
docker build -t payment-service .
docker run -p 8083:8083 payment-service
```

### Chạy với Docker Compose
```bash
docker-compose up -d
```

## API Endpoints

### Health Check
- `GET /health`: Kiểm tra trạng thái của service

### Payments
- `POST /payments`: Tạo thanh toán mới
- `GET /payments`: Lấy danh sách thanh toán
- `GET /payments/{id}`: Lấy thông tin thanh toán theo ID
- `GET /payments/order/{order_id}`: Lấy thông tin thanh toán theo order ID
- `PUT /payments/{id}/status`: Cập nhật trạng thái thanh toán

## Database Schema

### Payments Table
```sql
CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    card_number VARCHAR(16),
    expiry_month VARCHAR(2),
    expiry_year VARCHAR(4),
    cvv VARCHAR(4),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
```

## Kafka Events

### Produces
- `payment_created`: Khi thanh toán được tạo
- `payment_successful`: Khi thanh toán thành công
- `payment_failed`: Khi thanh toán thất bại
- `payment_refunded`: Khi thanh toán được hoàn tiền

### Consumes
- `order_cancelled`: Để hoàn tiền khi đơn hàng bị hủy

## Luồng xử lý thanh toán

1. **Tạo thanh toán**:
   - Order Service gửi request `POST /payments` với thông tin thanh toán
   - Service tạo thanh toán mới với trạng thái `PENDING`
   - Service xử lý thanh toán với payment gateway
   - Nếu thanh toán thành công, service cập nhật trạng thái thanh toán thành `SUCCESSFUL`
   - Nếu thanh toán thất bại, service cập nhật trạng thái thanh toán thành `FAILED`
   - Service gửi event `payment_successful` hoặc `payment_failed` đến Kafka

2. **Hoàn tiền**:
   - Payment Service nhận event `order_cancelled` từ Kafka
   - Service kiểm tra trạng thái thanh toán
   - Nếu thanh toán đã thành công, service cập nhật trạng thái thanh toán thành `REFUNDED`
   - Service gửi event `payment_refunded` đến Kafka

## Xử lý lỗi

- **Thanh toán thất bại**: Trả về lỗi và cập nhật trạng thái thanh toán thành `FAILED`
- **Phương thức thanh toán không hỗ trợ**: Trả về lỗi "unsupported payment method"
- **Thông tin thẻ tín dụng không hợp lệ**: Trả về lỗi "invalid credit card details"
- **Lỗi database**: Trả về lỗi 500 Internal Server Error
- **Lỗi Kafka**: Log lỗi và tiếp tục xử lý

## Bảo mật

- CORS middleware cho phép cross-origin requests
- Validation đầu vào để ngăn chặn các cuộc tấn công như SQL Injection
- Masking thông tin thẻ tín dụng trong response

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
