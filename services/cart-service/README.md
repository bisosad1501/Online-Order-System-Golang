# Cart Service

## Tổng quan
Cart Service quản lý giỏ hàng của khách hàng trong hệ thống xử lý đơn hàng trực tuyến. Service này cung cấp các API để tạo, đọc, cập nhật và xóa giỏ hàng, cũng như quản lý các sản phẩm trong giỏ hàng.

## Chức năng chính
- Quản lý giỏ hàng của khách hàng (CRUD)
- Thêm sản phẩm vào giỏ hàng
- Xóa sản phẩm khỏi giỏ hàng
- Cập nhật số lượng sản phẩm trong giỏ hàng
- Lấy thông tin giỏ hàng của khách hàng
- Xóa giỏ hàng sau khi đơn hàng được tạo

## Cấu trúc thư mục
```
cart-service/
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

## API Endpoints

### Health Check
- `GET /health`: Kiểm tra trạng thái của service

### Carts
- `POST /carts`: Tạo giỏ hàng mới
- `GET /carts/{id}`: Lấy thông tin giỏ hàng theo ID
- `GET /carts/customer/{customer_id}`: Lấy thông tin giỏ hàng theo customer ID
- `POST /carts/{id}/items`: Thêm sản phẩm vào giỏ hàng
- `PUT /carts/{id}/items/{item_id}`: Cập nhật số lượng sản phẩm trong giỏ hàng
- `DELETE /carts/{id}/items/{item_id}`: Xóa sản phẩm khỏi giỏ hàng
- `DELETE /carts/{id}`: Xóa giỏ hàng

## Database Schema

### Carts Table
```sql
CREATE TABLE IF NOT EXISTS carts (
    id VARCHAR(36) PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
```

### Cart Items Table
```sql
CREATE TABLE IF NOT EXISTS cart_items (
    id VARCHAR(36) PRIMARY KEY,
    cart_id VARCHAR(36) NOT NULL REFERENCES carts(id),
    product_id VARCHAR(36) NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
)
```

## Kafka Events

### Produces
- `cart_updated`: Khi giỏ hàng được cập nhật
- `cart_cleared`: Khi giỏ hàng được xóa sau khi đơn hàng được tạo

### Consumes
- `order_created`: Để xóa giỏ hàng sau khi đơn hàng được tạo

## Cài đặt và Chạy

### Yêu cầu
- Go 1.20 trở lên
- PostgreSQL
- Kafka

### Biến môi trường
- `PORT`: Port để chạy service (mặc định: 8087)
- `DB_HOST`: Host của PostgreSQL (mặc định: localhost)
- `DB_PORT`: Port của PostgreSQL (mặc định: 5432)
- `DB_USER`: Username của PostgreSQL (mặc định: postgres)
- `DB_PASSWORD`: Password của PostgreSQL (mặc định: postgres)
- `DB_NAME`: Tên database (mặc định: cartdb)
- `KAFKA_BOOTSTRAP_SERVERS`: Kafka bootstrap servers (mặc định: kafka:9092)
- `KAFKA_TOPIC`: Kafka topic (mặc định: carts)

### Chạy với Docker
```bash
docker build -t cart-service .
docker run -p 8087:8087 cart-service
```

### Chạy trực tiếp
```bash
cd src
go run main.go
```

## Luồng xử lý giỏ hàng

1. **Tạo giỏ hàng**:
   - Client gửi request `POST /carts` với customer_id
   - Service tạo giỏ hàng mới
   - Service lưu giỏ hàng vào database
   - Service trả về thông tin giỏ hàng

2. **Thêm sản phẩm vào giỏ hàng**:
   - Client gửi request `POST /carts/{id}/items` với thông tin sản phẩm
   - Service kiểm tra xem sản phẩm đã có trong giỏ hàng chưa
   - Nếu sản phẩm đã có, service cập nhật số lượng
   - Nếu sản phẩm chưa có, service thêm sản phẩm mới vào giỏ hàng
   - Service lưu thông tin vào database
   - Service gửi event `cart_updated` đến Kafka
   - Service trả về thông tin giỏ hàng đã cập nhật

3. **Cập nhật số lượng sản phẩm**:
   - Client gửi request `PUT /carts/{id}/items/{item_id}` với số lượng mới
   - Service cập nhật số lượng sản phẩm trong giỏ hàng
   - Service lưu thông tin vào database
   - Service gửi event `cart_updated` đến Kafka
   - Service trả về thông tin giỏ hàng đã cập nhật

4. **Xóa sản phẩm khỏi giỏ hàng**:
   - Client gửi request `DELETE /carts/{id}/items/{item_id}`
   - Service xóa sản phẩm khỏi giỏ hàng
   - Service lưu thông tin vào database
   - Service gửi event `cart_updated` đến Kafka
   - Service trả về thông tin giỏ hàng đã cập nhật

5. **Xóa giỏ hàng**:
   - Client gửi request `DELETE /carts/{id}`
   - Service xóa giỏ hàng và tất cả sản phẩm trong giỏ hàng
   - Service lưu thông tin vào database
   - Service gửi event `cart_cleared` đến Kafka
   - Service trả về thông báo thành công

6. **Xóa giỏ hàng sau khi đơn hàng được tạo**:
   - Service nhận event `order_created` từ Kafka
   - Service xóa giỏ hàng tương ứng với customer_id trong event
   - Service lưu thông tin vào database
   - Service gửi event `cart_cleared` đến Kafka

## Xử lý lỗi

- **Giỏ hàng không tồn tại**: Trả về lỗi 404 Not Found
- **Sản phẩm không tồn tại trong giỏ hàng**: Trả về lỗi 404 Not Found
- **Số lượng sản phẩm không hợp lệ**: Trả về lỗi 400 Bad Request
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
