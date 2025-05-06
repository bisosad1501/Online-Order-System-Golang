# User Service

## Overview
User Service quản lý thông tin người dùng trong hệ thống. Service này xử lý việc tạo, cập nhật, xóa và xác thực người dùng, cũng như quản lý thông tin đơn hàng của người dùng.

## Cấu trúc thư mục
```
user-service/
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
    └── docs/                 # Tài liệu API
        ├── swagger.html      # Swagger UI
        └── swagger.yaml      # Swagger specification
```

## Cài đặt và Chạy

### Yêu cầu
- Go 1.20 trở lên
- PostgreSQL
- Kafka

### Biến môi trường
- `PORT`: Port để chạy service (mặc định: 8086)
- `DB_HOST`: Host của PostgreSQL (mặc định: localhost)
- `DB_PORT`: Port của PostgreSQL (mặc định: 5432)
- `DB_USER`: Username của PostgreSQL (mặc định: useruser)
- `DB_PASSWORD`: Password của PostgreSQL (mặc định: userpass)
- `DB_NAME`: Tên database (mặc định: userdb)
- `KAFKA_BOOTSTRAP_SERVERS`: Kafka bootstrap servers (mặc định: kafka:9092)
- `KAFKA_TOPIC`: Kafka topic (mặc định: orders)

### Chạy với Docker
```bash
docker build -t user-service .
docker run -p 8086:8086 user-service
```

### Chạy với Docker Compose
```bash
docker-compose up -d
```

## API Endpoints

### Health Check
- `GET /health`: Kiểm tra trạng thái của service

### Users
- `POST /users`: Tạo người dùng mới
- `GET /users`: Lấy danh sách người dùng
- `GET /users/{id}`: Lấy thông tin người dùng theo ID
- `PUT /users/{id}`: Cập nhật thông tin người dùng
- `DELETE /users/{id}`: Xóa người dùng
- `GET /users/{id}/orders`: Lấy danh sách đơn hàng của người dùng
- `POST /users/verify`: Xác thực thông tin người dùng

## Database Schema

### Users Table
```sql
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### User Orders Table
```sql
CREATE TABLE IF NOT EXISTS user_orders (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36) NOT NULL,
    order_date TIMESTAMP NOT NULL,
    order_status VARCHAR(20) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## Kafka Events

### Produces
- `user_verified`: Khi người dùng được xác thực
- `user_updated`: Khi thông tin người dùng được cập nhật

### Consumes
- `order_created`: Để thêm đơn hàng mới vào danh sách đơn hàng của người dùng
- `order_updated`: Để cập nhật trạng thái đơn hàng của người dùng

## Xử lý lỗi

- **Người dùng không tồn tại**: Trả về lỗi 404 Not Found
- **Email đã tồn tại**: Trả về lỗi 400 Bad Request
- **Lỗi database**: Trả về lỗi 500 Internal Server Error
- **Lỗi Kafka**: Log lỗi và tiếp tục xử lý

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
