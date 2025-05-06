# Hệ thống Xử lý Đơn hàng Trực tuyến

Hệ thống xử lý đơn hàng trực tuyến dựa trên kiến trúc hướng dịch vụ (SOA), cho phép khách hàng đặt mua sản phẩm từ cửa hàng trực tuyến. Hệ thống xử lý toàn bộ quy trình từ khi khách hàng đặt hàng, kiểm tra tồn kho, xử lý thanh toán, lên lịch giao hàng, đến khi đơn hàng được giao thành công, tuân thủ các nguyên lý như Standardized Service Contract, Service Loose Coupling, Service Autonomy, và Service Composability.

## Kiến trúc

Hệ thống bao gồm các thành phần sau:

- **order-service (Task Service)**: Điều phối quy trình xử lý đơn hàng, phối hợp các dịch vụ khác để tạo, xác nhận, và hoàn tất đơn hàng.
- **user-service (Entity Service)**: Quản lý thông tin người dùng (tên, email, địa chỉ), xác minh và cập nhật hồ sơ.
- **cart-service (Entity Service)**: Quản lý giỏ hàng, lưu trữ sản phẩm và số lượng cho đến khi xác nhận đơn hàng.
- **inventory-service (Entity Service)**: Quản lý danh mục sản phẩm, tồn kho, kiểm tra tính sẵn có, khóa tồn kho, và gợi ý sản phẩm thay thế khi hết hàng.
- **payment-service (Entity Service)**: Tích hợp với cổng thanh toán bên ngoài để xử lý giao dịch và lỗi thanh toán.
- **shipping-service (Entity Service)**: Phối hợp với nhà vận chuyển để lập lịch giao hàng, theo dõi trạng thái, và cập nhật.
- **notification-service (Utility Service)**: Lưu và cung cấp thông báo trạng thái đơn hàng (xác nhận, cập nhật giao hàng, giao hàng thành công) qua REST API cho ứng dụng.
- **api-gateway (Infrastructure)**: Định tuyến yêu cầu HTTP, thực thi xác thực (JWT), giới hạn tỷ lệ (100 yêu cầu/giây), và chính sách CORS.
- **kafka (Infrastructure)**: Quản lý giao tiếp bất đồng bộ qua các topic như `order-events`, `payment-events`, `shipping-events`.
- **redis (Infrastructure)**: Caching dữ liệu giỏ hàng, tồn kho, gợi ý sản phẩm, và trạng thái giao hàng.

## Công nghệ

- **Ngôn ngữ**: Go (Golang)
- **Framework**: Gin
- **Cơ sở dữ liệu**: PostgreSQL (schema riêng cho mỗi dịch vụ)
- **Caching**: Redis
- **Message Broker**: Apache Kafka (8 partitions, retention 3 ngày)
- **API Gateway**: Traefik
- **Containerization**: Docker, Docker Compose
- **API Documentation**: Swagger/OpenAPI
- **Logging**: Logrus (structured logging)
- **Monitoring**: Prometheus (CPU, memory, request latency)

## Cài đặt

### Yêu cầu

- Docker
- Docker Compose

### Khởi động hệ thống

```bash
docker-compose up -d
```

Lệnh này sẽ:
1. Build và khởi động tất cả các dịch vụ
2. Tạo các Kafka topics cần thiết

### Dừng hệ thống

```bash
docker-compose down
```

## API Documentation

Mỗi service đều có tài liệu API được tạo bằng Swagger/OpenAPI. Bạn có thể truy cập tài liệu API của từng service tại các URL sau:

- **order-service**: http://localhost:8081/swagger
- **inventory-service**: http://localhost:8082/swagger
- **payment-service**: http://localhost:8083/swagger
- **shipping-service**: http://localhost:8084/swagger
- **notification-service**: http://localhost:8085/swagger
- **user-service**: http://localhost:8086/swagger
- **cart-service**: http://localhost:8087/swagger

## Quy trình nghiệp vụ

Hệ thống xử lý đơn hàng trực tuyến tuân theo quy trình 14 bước sau:

1. **Khách hàng đặt hàng**: Khách hàng duyệt qua các sản phẩm trên trang web, chọn sản phẩm cần mua và thêm vào giỏ hàng.
2. **Xác nhận đơn hàng**: Khách hàng xác nhận đơn hàng và nhập thông tin giao hàng, bao gồm địa chỉ giao hàng và phương thức thanh toán.
3. **Kiểm tra tình trạng hàng tồn kho**: Hệ thống kiểm tra tình trạng hàng tồn kho để đảm bảo sản phẩm còn hàng và có thể giao.
4. **Nếu hàng không có sẵn, thông báo khách hàng**: Nếu sản phẩm đã hết hàng, hệ thống thông báo cho khách hàng qua giao diện ứng dụng và đề xuất các sản phẩm tương tự hoặc hủy đơn hàng.
5. **Nếu hàng có sẵn, tiếp tục xử lý đơn hàng**: Nếu sản phẩm còn hàng, hệ thống tiếp tục xử lý đơn hàng.
6. **Xử lý thanh toán**: Hệ thống chuyển yêu cầu thanh toán đến nhà cung cấp dịch vụ thanh toán, như thẻ tín dụng hoặc ví điện tử.
7. **Nếu thanh toán thất bại, thông báo lỗi**: Nếu quá trình thanh toán gặp lỗi, hệ thống thông báo lỗi qua giao diện ứng dụng và yêu cầu khách hàng nhập lại thông tin thanh toán hoặc chọn phương thức khác.
8. **Nếu thanh toán thành công, xác nhận đơn hàng**: Nếu thanh toán thành công, hệ thống xác nhận đơn hàng và gửi thông báo cho khách hàng qua ứng dụng.
9. **Lên lịch giao hàng**: Hệ thống lên lịch giao hàng dựa trên địa chỉ giao hàng và thông tin từ đơn vị vận chuyển.
10. **Gửi thông tin đơn hàng đến đơn vị vận chuyển**: Hệ thống gửi thông tin về đơn hàng và địa chỉ giao hàng đến đơn vị vận chuyển.
11. **Theo dõi trạng thái giao hàng**: Hệ thống theo dõi trạng thái giao hàng và cập nhật cho khách hàng qua thông báo trong ứng dụng.
12. **Xác nhận đơn hàng đã giao**: Sau khi đơn hàng được giao thành công, hệ thống cập nhật trạng thái đơn hàng thành "Đã giao".
13. **Gửi thông báo giao hàng thành công**: Hệ thống gửi thông báo đến khách hàng qua ứng dụng rằng đơn hàng đã được giao thành công và yêu cầu khách hàng đánh giá sản phẩm.
14. **Lưu thông tin đơn hàng vào hệ thống quản lý**: Hệ thống lưu trữ thông tin đơn hàng vào cơ sở dữ liệu để theo dõi và quản lý các đơn hàng đã hoàn thành.

## Giao tiếp giữa các Service

Hệ thống sử dụng hai phương thức giao tiếp chính giữa các service:

1. **Giao tiếp đồng bộ (REST API)**: Được sử dụng cho các tác vụ yêu cầu phản hồi ngay lập tức, như:
   - Kiểm tra tồn kho (order-service → inventory-service)
   - Xử lý thanh toán (order-service → payment-service)
   - Lên lịch giao hàng (order-service → shipping-service)
   - Xác minh thông tin người dùng (order-service → user-service)
   - Lấy thông tin giỏ hàng (order-service → cart-service)

2. **Giao tiếp bất đồng bộ (Kafka)**: Được sử dụng cho các tác vụ không yêu cầu phản hồi ngay lập tức, như:
   - Thông báo đơn hàng đã tạo (order-service → notification-service)
   - Cập nhật trạng thái đơn hàng (order-service → các service khác)
   - Thông báo thanh toán thành công (payment-service → order-service)
   - Thông báo cập nhật trạng thái giao hàng (shipping-service → order-service)
   - Thông báo giỏ hàng đã xóa (cart-service → các service khác)

Các Kafka topics chính:
- `order-events`: Sự kiện liên quan đến đơn hàng (tạo, xác nhận, hoàn thành, hủy)
- `payment-events`: Sự kiện liên quan đến thanh toán (thành công, thất bại)
- `shipping-events`: Sự kiện liên quan đến giao hàng (lên lịch, cập nhật trạng thái)
- `inventory-events`: Sự kiện liên quan đến tồn kho (cập nhật, khóa, giải phóng)
- `cart-events`: Sự kiện liên quan đến giỏ hàng (tạo, cập nhật, xóa)
- `notification-events`: Sự kiện liên quan đến thông báo (gửi, đã đọc)

## Tài liệu

- [Phân tích và Thiết kế](docs/analysis-and-design.md)
- [Kiến trúc Hệ thống](docs/architecture.md)
- [API Specifications](docs/api-specs/)

## Truy cập UI

- **Kafka UI**: http://localhost:8090
- **API Gateway**: http://localhost:9090

## Mô hình dữ liệu

Hệ thống sử dụng mô hình Database-per-Service, trong đó mỗi service có schema riêng trong cùng một cơ sở dữ liệu PostgreSQL. Các service chỉ truy cập vào schema của chính nó và giao tiếp với các service khác thông qua API hoặc message broker.

Các schema chính:
- **order_schema**: Lưu trữ thông tin đơn hàng, trạng thái, và lịch sử đơn hàng
- **user_schema**: Lưu trữ thông tin người dùng, địa chỉ, và thông tin liên hệ
- **cart_schema**: Lưu trữ thông tin giỏ hàng và các mặt hàng trong giỏ
- **inventory_schema**: Lưu trữ thông tin sản phẩm, danh mục, và tồn kho
- **payment_schema**: Lưu trữ thông tin giao dịch thanh toán và trạng thái
- **shipping_schema**: Lưu trữ thông tin vận chuyển, lịch trình, và trạng thái
- **notification_schema**: Lưu trữ thông tin thông báo và trạng thái gửi

## Saga Pattern và Xử lý lỗi

Hệ thống triển khai mẫu thiết kế **Saga Orchestration** với order-service đóng vai trò là Saga Orchestrator, điều phối toàn bộ quy trình giao dịch phân tán và xử lý bù trừ (compensation) khi có lỗi xảy ra.

### Saga Orchestration Flow

1. **order-service** điều phối quy trình 14 bước:
   - Tạo đơn hàng (`status=CREATED`)
   - Gọi user-service xác minh người dùng (REST API)
   - Gọi cart-service lấy giỏ hàng (REST API)
   - Gọi inventory-service kiểm tra và khóa tồn kho (REST API), cập nhật `inventory_locked=true`
   - Gọi payment-service xử lý thanh toán (REST API), cập nhật `payment_processed=true`
   - Phát event `order_confirmed` khi thanh toán thành công, cập nhật `status=CONFIRMED`
   - Lên lịch giao hàng, cập nhật `shipping_scheduled=true`
   - Lắng nghe event `shipping_completed` để cập nhật `status=DELIVERED`

2. **Các service khác** phản hồi lại các lệnh từ Saga Orchestrator:
   - payment-service lắng nghe `order_confirmed` để xử lý thanh toán
   - shipping-service lắng nghe `payment_succeeded` để lập lịch giao hàng
   - notification-service lắng nghe các event để gửi thông báo

### Cơ chế bù trừ (Compensation)

Khi một bước trong quy trình thất bại, hệ thống thực hiện các hành động bù trừ để đưa hệ thống về trạng thái nhất quán trong thời gian tối đa 30 giây:

- **Lỗi kiểm tra tồn kho**:
  - Đơn hàng được đánh dấu là FAILED với `failure_reason="inventory_unavailable"`
  - Thông báo cho khách hàng và đề xuất sản phẩm thay thế
  - Ghi log vào AuditLogs với `action=inventory_error`

- **Lỗi thanh toán**:
  - Giải phóng tồn kho đã khóa (gọi inventory-service)
  - Đơn hàng được đánh dấu là FAILED với `failure_reason="payment_failed"`
  - Thông báo cho khách hàng để thử lại hoặc chọn phương thức thanh toán khác

- **Lỗi lên lịch giao hàng**:
  - Hoàn tiền cho khách hàng (gọi payment-service)
  - Đơn hàng được đánh dấu là FAILED với `failure_reason="shipping_failed"`
  - Thông báo cho khách hàng

- **Lỗi giao hàng**:
  - Cập nhật trạng thái đơn hàng
  - Thông báo cho khách hàng
  - Lên lịch giao hàng lại nếu cần

### Cơ chế Retry và Dead-Letter Queue

- **Retry với Exponential Backoff**:
  - REST API: Retry 2 lần (backoff 100ms, 200ms)
  - API nhà vận chuyển: Retry 3 lần (backoff 1s)
  - Kafka consumer: Retry 2 lần (backoff 1s, 2s, 4s, 8s, 16s)

- **Dead-Letter Queue (DLQ)**:
  - Sự kiện Kafka thất bại gửi đến DLQ
  - Xử lý hàng ngày
  - Alert qua Prometheus nếu >10 message trong DLQ

## Tác giả

- Lê Đức Thắng - B21DCDT205
- Nguyễn Đức Tài - B21DCDT199
- Nguyễn Vũ Duy Anh- B21DCVT004
