# Inventory Service

## Overview
Inventory Service quản lý thông tin sản phẩm và kiểm tra tình trạng tồn kho. Service này xử lý việc hiển thị thông tin sản phẩm, kiểm tra tình trạng tồn kho, cập nhật số lượng tồn kho khi đơn hàng được xác nhận, và cung cấp các API gợi ý sản phẩm. Service sử dụng Redis để cache dữ liệu, giúp tăng hiệu suất và giảm tải cho database.

## Cấu trúc thư mục
```
inventory-service/
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
    ├── cache/                # Tương tác với Redis cache
    │   └── redis.go          # Redis cache client
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
- Redis (cho caching)

### Biến môi trường
- `PORT`: Port để chạy service (mặc định: 8082)
- `DB_HOST`: Host của PostgreSQL (mặc định: localhost)
- `DB_PORT`: Port của PostgreSQL (mặc định: 5432)
- `DB_USER`: Username của PostgreSQL (mặc định: postgres)
- `DB_PASSWORD`: Password của PostgreSQL (mặc định: postgres)
- `DB_NAME`: Tên database (mặc định: inventorydb)
- `KAFKA_BOOTSTRAP_SERVERS`: Kafka bootstrap servers (mặc định: localhost:9092)
- `KAFKA_TOPIC`: Kafka topic (mặc định: inventory)
- `REDIS_HOST`: Host của Redis (mặc định: localhost)
- `REDIS_PORT`: Port của Redis (mặc định: 6379)
- `REDIS_PASSWORD`: Password của Redis (mặc định: "")
- `REDIS_DB`: Redis database index (mặc định: 0)
- `REDIS_CACHE_TTL`: Thời gian cache hết hạn (mặc định: 3600 giây)

### Chạy với Docker
```bash
docker build -t inventory-service .
docker run -p 8082:8082 inventory-service
```

### Chạy với Docker Compose
```bash
docker-compose up -d
```

## API Endpoints

### Health Check
- `GET /health`: Kiểm tra trạng thái của service

### Products
- `POST /products`: Tạo sản phẩm mới
- `GET /products`: Lấy danh sách sản phẩm
- `GET /products/{id}`: Lấy thông tin sản phẩm theo ID
- `PUT /products/{id}`: Cập nhật thông tin sản phẩm
- `DELETE /products/{id}`: Xóa sản phẩm

### Inventory
- `PUT /inventory/update`: Cập nhật số lượng tồn kho
- `GET /inventory/{id}`: Lấy thông tin tồn kho theo ID sản phẩm
- `POST /inventory/check`: Kiểm tra tình trạng tồn kho

### Recommendations
- `GET /recommendations/product/{id}`: Lấy gợi ý sản phẩm dựa trên ID sản phẩm
- `GET /recommendations/category/{id}`: Lấy gợi ý sản phẩm dựa trên ID danh mục
- `POST /recommendations/similar`: Lấy các sản phẩm tương tự dựa trên nhiều tiêu chí

## Database Schema

### Products Table
```sql
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
```

### Inventory Table
```sql
CREATE TABLE IF NOT EXISTS inventory (
    product_id VARCHAR(36) PRIMARY KEY REFERENCES products(id),
    quantity INTEGER NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
```

### Product Tags Table
```sql
CREATE TABLE IF NOT EXISTS product_tags (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL REFERENCES products(id),
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL
)
```

### Product Similarities Table
```sql
CREATE TABLE IF NOT EXISTS product_similarities (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL REFERENCES products(id),
    similar_product_id VARCHAR(36) NOT NULL REFERENCES products(id),
    similarity_score DECIMAL(5, 4) NOT NULL,
    created_at TIMESTAMP NOT NULL
)
```

## Kafka Events

### Produces
- `inventory_updated`: Khi số lượng tồn kho được cập nhật

### Consumes
- `order_created`: Để cập nhật số lượng tồn kho khi đơn hàng được tạo

## Luồng xử lý tồn kho

1. **Tạo sản phẩm**:
   - Client gửi request `POST /products` với thông tin sản phẩm
   - Service tạo sản phẩm mới và lưu vào database
   - Service tạo bản ghi tồn kho mới với số lượng ban đầu là 0

2. **Cập nhật tồn kho**:
   - Client gửi request `PUT /inventory/update` với thông tin sản phẩm và số lượng
   - Service cập nhật số lượng tồn kho trong database
   - Service gửi event `inventory_updated` đến Kafka

3. **Kiểm tra tồn kho**:
   - Client gửi request `POST /inventory/check` với danh sách sản phẩm và số lượng cần kiểm tra
   - Service kiểm tra số lượng tồn kho cho từng sản phẩm
   - Service trả về kết quả kiểm tra

4. **Xử lý đơn hàng**:
   - Order Service gửi event `order_created` đến Kafka
   - Inventory Service nhận event `order_created` và cập nhật số lượng tồn kho cho từng sản phẩm trong đơn hàng

## Xử lý lỗi

- **Sản phẩm không tồn tại**: Trả về lỗi 404 Not Found
- **Số lượng tồn kho không đủ**: Trả về danh sách sản phẩm không đủ số lượng
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
