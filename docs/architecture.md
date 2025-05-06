# Kiến trúc Hệ thống

## Tổng quan

Hệ thống **Xử lý Đơn hàng Trực tuyến** là một kiến trúc microservices được thiết kế để quản lý quy trình đặt hàng trực tuyến cho nền tảng thương mại điện tử. Hệ thống cho phép khách hàng duyệt sản phẩm, đặt hàng, kiểm tra tồn kho, xử lý thanh toán, và quản lý giao hàng, đảm bảo quy trình tự động, hiệu quả và chính xác. Hệ thống được xây dựng dựa trên các nguyên lý Kiến trúc Hướng Dịch vụ (SOA), nhấn mạnh tính tách rời, tự chủ, khám phá dịch vụ, và khả năng tổ hợp để hỗ trợ bảo trì và mở rộng.

---

## Các Thành phần Chính và Vai trò

- **order-service (Task Service)**:
  - **Chức năng**: Điều phối quy trình xử lý đơn hàng, phối hợp các dịch vụ khác để tạo, xác nhận, và hoàn tất đơn hàng.
  - **Vai trò**: Trung tâm đảm bảo quy trình 14 bước (từ đặt hàng đến lưu trữ) được thực hiện trơn tru.
- **user-service (Entity Service)**:
  - **Chức năng**: Quản lý thông tin người dùng (tên, email, địa chỉ), xác minh và cập nhật hồ sơ.
  - **Vai trò**: Cung cấp dữ liệu người dùng để xác minh và cá nhân hóa thông báo.
- **cart-service (Entity Service)**:
  - **Chức năng**: Quản lý giỏ hàng, lưu trữ sản phẩm và số lượng cho đến khi xác nhận đơn hàng.
  - **Vai trò**: Khởi tạo quy trình đặt hàng bằng cách cung cấp chi tiết giỏ hàng.
- **inventory-service (Entity Service)**:
  - **Chức năng**: Quản lý danh mục sản phẩm, tồn kho, kiểm tra tính sẵn có, khóa tồn kho, và gợi ý sản phẩm thay thế khi hết hàng.
  - **Vai trò**: Đảm bảo sản phẩm có sẵn trước khi thanh toán và giao hàng.
- **payment-service (Entity Service)**:
  - **Chức năng**: Tích hợp với cổng thanh toán bên ngoài để xử lý giao dịch và lỗi thanh toán.
  - **Vai trò**: Xác minh và xử lý thanh toán an toàn.
- **shipping-service (Entity Service)**:
  - **Chức năng**: Phối hợp với nhà vận chuyển để lập lịch giao hàng, theo dõi trạng thái, và cập nhật.
  - **Vai trò**: Quản lý hậu cần để hoàn tất giao hàng.
- **notification-service (Utility Service)**:
  - **Chức năng**: Lưu và cung cấp thông báo trạng thái đơn hàng (xác nhận, cập nhật giao hàng, giao hàng thành công) qua REST API cho ứng dụng.
  - **Vai trò**: Giữ người dùng được thông báo về trạng thái đơn hàng.
- **api-gateway (Infrastructure)**:
  - **Chức năng**: Định tuyến yêu cầu HTTP, thực thi xác thực (JWT), giới hạn tỷ lệ (100 yêu cầu/giây), và chính sách CORS.
  - **Vai trò**: Điểm vào duy nhất cho client, tăng cường bảo mật và hiệu suất.
- **kafka (Infrastructure)**:
  - **Chức năng**: Quản lý giao tiếp bất đồng bộ qua các topic như `order-events`, `payment-events`, `shipping-events`.
  - **Vai trò**: Đảm bảo tính tách rời và thông lượng cao.
- **redis (Infrastructure)**:
  - **Chức năng**: Caching dữ liệu giỏ hàng, tồn kho, gợi ý sản phẩm, và trạng thái giao hàng.
  - **Vai trò**: Tăng hiệu suất truy vấn dữ liệu.

---

## Giao tiếp

- **Giao tiếp giữa các Dịch vụ**:
  - **REST API** (đồng bộ): Dùng cho các bước cần phản hồi tức thì:
    - **order-service → user-service**: Xác minh người dùng (Bước 2).
    - **order-service → cart-service**: Lấy giỏ hàng (Bước 1-2).
    - **order-service → inventory-service**: Kiểm tra, khóa tồn kho, gợi ý sản phẩm (Bước 3, 4, 5).
    - **order-service → payment-service**: Xử lý thanh toán (Bước 6, 7, 8).
    - **order-service → shipping-service**: Lập lịch giao hàng (Bước 9-12).
    - **client → notification-service**: Truy vấn thông báo trạng thái (Bước 8, 11, 13).
    - Timeout 5 giây, retry 2 lần (backoff 100ms, 200ms), trừ API nhà vận chuyển (retry 3 lần, backoff 1 giây).
  - **Kafka** (bất đồng bộ): Dùng cho các tác vụ nền:
    - `order_confirmed`: Kích hoạt thanh toán (Bước 8).
    - `payment_succeeded`: Kích hoạt giao hàng (Bước 9).
    - `shipping_scheduled`, `shipping_status_updated`, `shipping_completed`: Cập nhật trạng thái và thông báo (Bước 11, 12, 13).
    - Schema Avro, 8 partitions, retention 3 ngày, Dead-Letter Queue (DLQ) xử lý hàng ngày, alert qua Prometheus nếu >10 message.
- **Mạng Nội bộ**:
  - Các dịch vụ giao tiếp qua tên dịch vụ Docker Compose (DNS-based discovery), ví dụ: `order-service`, `inventory-service`.

---

## Luồng Dữ liệu

1. **Yêu cầu từ Client**: Client gửi yêu cầu HTTP POST để tạo đơn hàng qua api-gateway, bao gồm ID người dùng, giỏ hàng, địa chỉ giao hàng, và phương thức thanh toán.
2. **Tạo Đơn hàng**: api-gateway định tuyến đến order-service, dịch vụ này:
   - Gọi user-service (REST) để xác minh người dùng.
   - Gọi cart-service (REST) để lấy giỏ hàng.
   - Gọi inventory-service (REST) để kiểm tra và khóa tồn kho.
3. **Xử lý Tồn kho**: inventory-service kiểm tra tồn kho, khóa hàng, trả lỗi 409 nếu không đủ (Bước 4). Nếu hết hàng, gợi ý sản phẩm thay thế hoặc hủy đơn.
4. **Xử lý Thanh toán**: order-service gọi payment-service (REST) để xử lý thanh toán, phát `payment_succeeded` hoặc trả lỗi 400 nếu thất bại (Bước 6, 7).
5. **Giao hàng**: shipping-service đăng ký `payment_succeeded`, lập lịch giao hàng, phát `shipping_scheduled`, `shipping_status_updated`, `shipping_completed` (Bước 9-12).
6. **Thông báo**: notification-service đăng ký `order_confirmed`, `shipping_status_updated`, `shipping_completed`, lưu thông báo trạng thái vào database, cung cấp qua REST API (Bước 8, 11, 13).
7. **Hoàn tất Đơn hàng**: order-service đăng ký `shipping_completed`, cập nhật trạng thái "Đã giao" và lưu vào database (Bước 12, 14).
8. **Cập nhật Người dùng**: user-service đăng ký `order_confirmed`, `order_completed`, cập nhật hồ sơ người dùng (Bước 8, 13).

**Phụ thuộc Bên ngoài**:
- **Cơ sở Dữ liệu**: PostgreSQL (schema riêng cho mỗi dịch vụ), Redis (caching).
- **API Bên thứ ba**:
  - payment-service tích hợp với Stripe.
  - shipping-service tích hợp với FedEx/UPS.
- **AWS S3**: Lưu trữ snapshot database và Kafka (nén gzip, retention 1 năm).

---

## Cơ sở Dữ liệu

- **order-service (PostgreSQL)**:
  - **Orders**: `id` (UUID), `customer_id`, `status` (enum: CREATED, CONFIRMED, DELIVERED), `total_amount` (numeric), `shipping_address`, `created_at`, `updated_at` (index: `id`, `customer_id`, retention: 1 năm)
  - **OrderItems**: `id` (UUID), `order_id`, `product_id`, `quantity` (integer), `price` (numeric) (index: `order_id`, retention: 1 năm)
  - **AuditLogs**: `id` (UUID), `service_name`, `action` (enum: create_order, process_payment, inventory_error, shipment_updated), `user_id`, `timestamp`, `details` (index: `user_id`, `timestamp`, `action`, retention: 6 tháng, max 1GB)
- **user-service (PostgreSQL)**:
  - **Customers**: `id` (UUID), `name`, `email`, `address`, `created_at`, `updated_at` (index: `id`, `email`, retention: 1 năm)
- **cart-service (PostgreSQL)**:
  - **Carts**: `id` (UUID), `customer_id`, `created_at`, `updated_at` (index: `customer_id`, retention: 90 ngày)
  - **CartItems**: `id` (UUID), `cart_id`, `product_id`, `quantity` (integer) (index: `cart_id`, retention: 90 ngày)
  - **Cache**: Redis key `cart:{user_id}` (TTL: 3600 giây, memory limit: 100MB, invalidation khi cập nhật hoặc hết hạn 90 ngày)
- **inventory-service (PostgreSQL)**:
  - **Products**: `id` (UUID), `name`, `description`, `price` (numeric), `category`, `created_at` (index: `id`, `category`, retention: 1 năm)
  - **ProductTags**: `product_id`, `tag` (index: `product_id`, `tag`, retention: 1 năm)
  - **Inventory**: `product_id`, `available_quantity` (integer), `locked_quantity` (integer), `updated_at` (index: `product_id`, retention: 1 năm)
  - **ProductRelations**: `product_id`, `related_product_id`, `relation_type` (enum: category_match, tag_match), `score` (float, 0-1), `updated_at` (index: `product_id`, retention: 1 năm)
  - **Cache**: Redis key `inventory:{product_id}` (TTL: 60 giây), `recommendation:{product_id}` (TTL: 3600 giây, memory limit: 100MB)
- **payment-service (PostgreSQL)**:
  - **Payments**: `id` (UUID), `order_id`, `amount` (numeric), `status` (enum: SUCCEEDED, FAILED), `payment_method`, `transaction_id`, `created_at` (index: `order_id`, retention: 1 năm)
- **shipping-service (PostgreSQL)**:
  - **Shipments**: `id` (UUID), `order_id`, `carrier`, `tracking_number`, `status` (enum: SCHEDULED, IN_TRANSIT, DELIVERED), `shipping_address`, `created_at`, `updated_at` (index: `order_id`, retention: 1 năm)
  - **ShipmentUpdates**: `id` (UUID), `shipment_id`, `status`, `description`, `created_at` (index: `shipment_id`, retention: 6 tháng)
  - **Cache**: Redis key `shipment:{order_id}` (TTL: 300 giây, memory limit: 100MB)
- **notification-service (PostgreSQL)**:
  - **Notifications**: `id` (UUID), `user_id`, `type` (enum: ORDER_CONFIRMED, SHIPPING_UPDATED, SHIPPING_COMPLETED), `content`, `status` (enum: UNREAD, READ), `created_at` (index: `user_id`, `created_at`, retention: 6 tháng, max 100 notifications/user)

**Backup**: Snapshot hàng ngày, nén bằng gzip, lưu trong AWS S3 (dung lượng tối đa 100GB, retention 1 năm). Khôi phục bằng `pg_restore`.

---

## Bảo mật

- **Xác thực và Ủy quyền**: JWT kiểm tra tại api-gateway, RBAC với vai trò `customer` (đặt đơn) và `admin` (quản lý).
- **Bảo mật Giao tiếp**: HTTPS với JWT, Kafka mã hóa SSL/TLS, Redis mã hóa SSL/TLS.
- **Bảo mật Dữ liệu**: Mã hóa thanh toán bằng AES-256.
- **Xác thực Đầu vào**: Ngăn SQL Injection, XSS bằng input validation.
- **Rate Limiting & CORS**: Giới hạn 100 yêu cầu/giây, trả lỗi 429, CORS whitelist `*.example.com`.
- **Quản lý Secrets**: Lưu trong Docker Compose secrets, xoay vòng hàng tháng.
- **Audit Logging**: Ghi hành động nhạy cảm (create_order, process_payment, inventory_error, shipment_updated) vào **AuditLogs** (max 1GB).
- **Security Testing**: Penetration testing hàng năm, vulnerability scanning hàng quý với OWASP ZAP, vá lỗ hổng trong 7 ngày.

---

## Triển khai

- **Containerization**: Mỗi dịch vụ chạy trong Docker container, quản lý bằng Docker Compose.
- **Phần cứng**: Tối thiểu 2 CPU, 4GB RAM, 50GB SSD.
- **Auto-Scaling**: Scaling thủ công bằng `docker-compose up --scale`.
- **CI/CD**: GitHub Actions, rolling update với `docker-compose --no-cache` để tránh downtime.
- **Resilience & Observability**:
  - Readiness/liveness probes qua endpoint `/health`.
  - Retry 2 lần cho API (trừ API nhà vận chuyển: 3 lần), Kafka consumer, lưu thông báo.
  - Structured logging với Logrus, metrics với Prometheus (CPU, memory, request latency).
- **Service Discovery**: DNS-based discovery (Docker Compose DNS).
- **Caching**: Redis 1 node, memory limit 100MB, invalidation khi Kafka event hoặc dữ liệu cập nhật.
- **Disaster Recovery**: Backup database và Kafka snapshot hàng ngày (S3, gzip, retention 1 năm), khôi phục bằng `pg_restore`.

---

## Sơ đồ Kiến trúc

Sơ đồ thể hiện:
- Client gửi yêu cầu qua **api-gateway**.
- **order-service** gọi REST API đến **user-service**, **cart-service**, **inventory-service**, **payment-service**, **shipping-service**.
- **notification-service** cung cấp thông báo trạng thái qua REST API.
- **Kafka** quản lý sự kiện bất đồng bộ (`order_confirmed`, `payment_succeeded`, `shipping_scheduled`, `shipping_status_updated`, `shipping_completed`).
- Mỗi dịch vụ có schema **PostgreSQL** riêng, **Redis** caching cho giỏ hàng, tồn kho, giao hàng.
- Compensation flow (rollback) từ **order-service** đến **inventory-service** (unlock tồn kho), **payment-service** (hoàn tiền), **cart-service** (xóa giỏ hàng).

**Mô tả văn bản**:
```
                                 ┌─────────────┐
                                 │    Client   │
                                 └──────┬──────┘
                                        │
                                        ▼
                                 ┌─────────────┐
                                 │ api-gateway │
                                 └──────┬──────┘
                                        │
         ┌───────────────────────┬─────┴─────┬───────────────────────┬──────┐
         │                       │           │                       │      │
         ▼                       ▼           ▼                       ▼      ▼
┌─────────────────┐     ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  ┌─────────────┐
│  order-service  │<--->│  inventory- │    │   payment-  │    │  shipping-  │  │notification-│
│ (Task Service)  │     │   service   │    │   service   │    │   service   │  │   service   │
└────────┬────────┘     └──────┬──────┘    └──────┬──────┘    └──────┬──────┘  └──────┬──────┘
         │                     │                  │                  │                │
         ▼                     ▼                  ▼                  ▼                ▼
┌─────────────────┐     ┌─────────────┐    ┌─────────────┐    ┌─────────────┐   ┌─────────────┐
│  Order/AuditLogs│     │  Inventory  │    │   Payment   │    │   Shipping  │   │ Notification│
│       DB        │     │     DB      │    │     DB      │    │     DB      │   │     DB      │
└─────────────────┘     └──────┬──────┘    └─────────────┘    └──────┬──────┘   └─────────────┘
                               │                                   │
                               ▼                                   ▼
                         ┌─────────────┐                     ┌─────────────┐
                         │    Redis    │                     │    Redis    │
                         └─────────────┘                     └─────────────┘
         │
         │                            ┌─────────────┐
         └───────────────────────────>│  user-      │
                                      │  service    │
                                      └──────┬──────┘
                                             │
                                             ▼
                                      ┌─────────────┐
                                      │ User        │
                                      │     DB      │
                                      └─────────────┘
         │
         │                            ┌─────────────┐
         └───────────────────────────>│   cart-     │
                                      │  service    │
                                      └──────┬──────┘
                                             │
                                             ▼
                                      ┌─────────────┐
                                      │ Cart        │
                                      │     DB      │
                                      └──────┬──────┘
                                             │
                                             ▼
                                      ┌─────────────┐
                                      │    Redis    │
                                      └─────────────┘
         │
         ▼
┌─────────────────┐
│      kafka      │
│  Event Broker   │
└──────┬──────────┘
       │
       └─────────────────────────┬─────────────┬──────────────────────┬──────┐
                                 │             │                      │      │
                                 ▼             ▼                      ▼      ▼
                          ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  ┌─────────────┐
                          │  inventory- │    │   payment-  │    │  shipping-  │  │notification-│
                          │   service   │    │   service   │    │   service   │  │   service   │
                          └──────┬──────┘    └──────┬──────┘    └──────┬──────┘  └─────────────┘
                                 │                  │                  │
                                 ▼                  ▼                  ▼
                          ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
                          │  Inventory  │    │   Payment   │    │   Shipping  │
                          │     DB      │    │     DB      │    │     DB      │
                          └─────────────┘    └─────────────┘    └─────────────┘
```

---

## Khả năng Mở rộng & Chịu lỗi

- **Khả năng Mở rộng**:
  - **Mở rộng Ngang**: Scaling thủ công bằng Docker Compose (`docker-compose up --scale`), tối thiểu 1 instance/dịch vụ, giám sát CPU/memory qua Prometheus.
  - **Kafka**: 8 phân vùng/topic đảm bảo thông lượng cao, xử lý song song.
  - **DNS-based Discovery**: Tự động ánh xạ tên dịch vụ sang IP/cổng trong Docker Compose.
- **Chịu lỗi**:
  - **Saga Orchestrator**: order-service quản lý giao dịch phân tán, rollback (30 giây) khi lỗi (hết hàng, thanh toán thất bại).
  - **Retry**: Retry 2 lần (backoff 100ms, 200ms) cho REST API, 3 lần (backoff 1 giây) cho API nhà vận chuyển, 2 lần cho Kafka consumer.
  - **DLQ**: Sự kiện Kafka thất bại gửi đến DLQ, xử lý hàng ngày, alert nếu >10 message.
  - **Health Check**: Endpoint `/health` cho readiness/liveness probes.
  - **Observability**: Logrus (structured logging), Prometheus (metrics), phát hiện sự cố nhanh.