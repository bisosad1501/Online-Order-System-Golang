# 📊 Hệ thống Xử lý Đơn hàng Trực tuyến - Phân tích và Thiết kế

Tài liệu này trình bày quá trình **phân tích** và **thiết kế** cho hệ thống xử lý đơn hàng trực tuyến dựa trên kiến trúc hướng dịch vụ (SOA), tuân thủ các nguyên lý như Standardized Service Contract, Service Loose Coupling, Service Autonomy, và Service Composability. Hệ thống đảm bảo quy trình từ đặt hàng đến giao hàng được thực hiện **tự động** (thông qua Saga và Kafka) và **hiệu quả** (với REST timeout, Kafka partitioning, Redis caching, và tối ưu tài nguyên).

---

## 1. 🎯 Mô tả Nghiệp vụ

Hệ thống xử lý đơn hàng trực tuyến cho phép khách hàng đặt mua sản phẩm từ cửa hàng trực tuyến. Hệ thống tự động xử lý đơn hàng bằng cách xác minh tình trạng hàng trong kho, xác nhận đơn hàng, xử lý thanh toán, và lên lịch giao hàng, đảm bảo hiệu quả và chính xác. Nếu đơn hàng được thanh toán thành công, sản phẩm sẽ được gửi đến khách hàng.

### Quy trình nghiệp vụ (14 bước):

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

---

## 2. 🧩 Các Dịch vụ đã Xác định

| Service Name | Type | Responsibility | Tech Stack |
|--------------|------|----------------|------------|
| order-service | Task Service | Điều phối quy trình đơn hàng từ tạo đến hoàn thành, phát Kafka event để kích hoạt thông báo trạng thái | Go, Gin, PostgreSQL |
| user-service | Entity Service | Quản lý thông tin người dùng (tên, email, địa chỉ), xác minh và cập nhật hồ sơ | Go, Gin, PostgreSQL |
| cart-service | Entity Service | Quản lý giỏ hàng của người dùng, lưu trữ sản phẩm và số lượng, caching giỏ hàng | Go, Gin, PostgreSQL, Redis |
| inventory-service | Entity Service | Quản lý danh mục sản phẩm và tồn kho, kiểm tra và khóa tồn kho, gợi ý sản phẩm tương tự dựa trên category/tags, caching dữ liệu tồn kho và gợi ý | Go, Gin, PostgreSQL, Redis |
| payment-service | Entity Service | Xử lý giao dịch thanh toán, tích hợp với cổng thanh toán (Stripe) | Go, Gin, PostgreSQL |
| shipping-service | Entity Service | Quản lý giao hàng, lập lịch và theo dõi trạng thái với nhà vận chuyển (FedEx, UPS), caching trạng thái giao hàng | Go, Gin, PostgreSQL, Redis |
| notification-service | Utility Service | Lưu và quản lý thông báo trạng thái đơn hàng trong ứng dụng (REST API, database) cho xác nhận đơn hàng, cập nhật giao hàng, giao hàng thành công | Go, Gin, PostgreSQL |
| api-gateway | Infrastructure | Định tuyến yêu cầu từ người dùng, thực thi bảo mật | Traefik |
| kafka | Infrastructure | Xử lý sự kiện bất đồng bộ với topic phân vùng | Apache Kafka |
| redis | Infrastructure | Caching dữ liệu giỏ hàng, tồn kho, gợi ý sản phẩm, và trạng thái giao hàng | Redis |

### Phân bổ các bước nghiệp vụ:
- **cart-service**: Bước 1
- **user-service**: Bước 2, 8, 13
- **order-service**: Bước 2, 5, 8, 12, 14
- **inventory-service**: Bước 3, 4, 5
- **payment-service**: Bước 6, 7, 8
- **shipping-service**: Bước 9, 10, 11, 12
- **notification-service**: Bước 8, 11, 13 (lưu thông báo trạng thái vào database)

**Ghi chú**:
- Thông báo lỗi (Bước 4, 7) được xử lý qua REST API trả lỗi trực tiếp (409, 400) và log vào AuditLogs/Logrus, không lưu trong database.
- Tính năng gửi email (qua SMTP) không phải cốt lõi, hiện chưa triển khai. Có thể thêm sau bằng AWS SES hoặc SendGrid nếu cần.

---

## 3. 🔄 Giao tiếp giữa các Dịch vụ

Hệ thống sử dụng giao tiếp đồng bộ (REST API) cho các bước cần phản hồi tức thì và bất đồng bộ (Kafka) cho các tác vụ nền, với Saga Orchestrator để đảm bảo tuần tự và hiệu quả.

### 3.1. Giao tiếp đồng bộ (REST API)

**Khi nào dùng**:
- **Bước 1**: cart-service lấy giỏ hàng (phản hồi tức thì).
- **Bước 2**: user-service xác minh người dùng, order-service tạo đơn hàng (phản hồi tức thì).
- **Bước 3, 4, 5**: inventory-service kiểm tra và khóa tồn kho, gợi ý sản phẩm, trả lỗi nếu hết hàng (phản hồi tức thì).
- **Bước 6, 7, 8**: payment-service xử lý thanh toán, trả lỗi nếu thất bại (phản hồi tức thì).
- **Bước 9, 10, 12**: shipping-service lập lịch và xác nhận giao hàng (phản hồi để cập nhật trạng thái).
- **Bước 8, 11, 13**: notification-service cung cấp thông báo trạng thái qua REST API cho client.

**Tương tác**:
- **api-gateway ⇄ Tất cả dịch vụ**: Định tuyến yêu cầu từ client.
- **order-service ⇄ user-service**: Xác minh người dùng (Bước 2).
- **order-service ⇄ cart-service**: Lấy giỏ hàng (Bước 1-2).
- **order-service ⇄ inventory-service**: Kiểm tra, khóa tồn kho, gợi ý sản phẩm (Bước 3, 4, 5).
- **order-service ⇄ payment-service**: Xử lý thanh toán (Bước 6, 7, 8).
- **order-service ⇄ shipping-service**: Lập lịch giao hàng (Bước 9-12).
- **client ⇄ notification-service**: Truy vấn thông báo (Bước 8, 11, 13).

**Hợp đồng dịch vụ**:
- Các dịch vụ cung cấp REST API với hợp đồng chuẩn hóa, bao gồm các endpoint để tạo đơn hàng, xác minh người dùng, quản lý giỏ hàng, kiểm tra tồn kho, xử lý thanh toán, lập lịch giao hàng, và truy vấn thông báo trạng thái.
- Mã lỗi: 400 (Invalid request), 401 (Unauthorized), 409 (Conflict, e.g., inventory not available), 404 (Not found).
- Timeout: 5 giây cho tất cả API, trừ `/shipments` (retry 3 lần, backoff 1 giây).

### 3.2. Giao tiếp bất đồng bộ (Kafka)

**Khi nào dùng**:
- **Bước 8**: notification-service lưu thông báo xác nhận đơn hàng (nền).
- **Bước 11**: notification-service lưu thông báo cập nhật trạng thái giao hàng (nền).
- **Bước 13**: notification-service lưu thông báo giao hàng thành công (nền).
- **Bước 8, 12**: order-service cập nhật trạng thái sau khi nhận sự kiện từ payment-service (`payment_succeeded`) hoặc shipping-service (`shipping_completed`).

**Kafka Topics**:
- **order-events**: `order_confirmed`, `order_completed` (8 partitions, retention: 3 ngày)
- **payment-events**: `payment_succeeded` (8 partitions, retention: 3 ngày)
- **shipping-events**: `shipping_scheduled`, `shipping_status_updated`, `shipping_completed` (8 partitions, retention: 3 ngày)

**Service Subscriptions**:
- **payment-service**: Lắng nghe `order_confirmed` để xử lý thanh toán (Bước 6).
- **shipping-service**: Lắng nghe `payment_succeeded` để lập lịch giao hàng (Bước 9).
- **order-service**: Lắng nghe `payment_succeeded`, `shipping_completed` để cập nhật trạng thái (Bước 8, 12).
- **user-service**: Lắng nghe `order_confirmed`, `order_completed` để cập nhật thông tin người dùng (Bước 8, 13).
- **notification-service**: Lắng nghe `order_confirmed`, `shipping_status_updated`, `shipping_completed` để lưu thông báo (Bước 8, 11, 13).

**Error Handling**:
- **Dead-Letter Queue (DLQ)**: Sự kiện thất bại gửi đến DLQ, xử lý hàng ngày, alert qua Prometheus nếu >10 message trong DLQ.
- **Retry Mechanism**: Retry 2 lần với exponential backoff cho Kafka consumer và lưu thông báo.

### 3.3. Saga Orchestrator & Compensation Flow

**order-service** điều phối Saga:
1. Tạo đơn hàng (`status=CREATED`, Bước 2, validate địa chỉ và phương thức thanh toán).
2. Gọi user-service xác minh người dùng (REST, Bước 2, retry 2 lần).
3. Gọi cart-service lấy giỏ hàng (REST, Bước 1-2, retry 2 lần).
4. Gọi inventory-service kiểm tra và khóa tồn kho (REST, Bước 3, retry 2 lần), cập nhật `inventory_locked=true`.
5. Nếu tồn kho không đủ, gọi inventory-service gợi ý sản phẩm (REST, Bước 4), trả lỗi 409 nếu khách hàng hủy (rollback tối đa 30 giây).
6. Gọi payment-service xử lý thanh toán (REST, Bước 6, retry 2 lần), cập nhật `payment_processed=true`.
7. Nếu thanh toán thành công, phát `order_confirmed` (Bước 8), cập nhật `status=CONFIRMED`.
8. Lên lịch giao hàng (Bước 9), cập nhật `shipping_scheduled=true`.
9. Nếu thất bại (hết hàng hoặc thanh toán lỗi):
   - Cập nhật `status=FAILED` và lưu `failure_reason`.
   - Compensation (rollback tối đa 30 giây):
     - Unlock tồn kho (inventory-service, Bước 4, xóa cache `inventory:{product_id}`).
     - Refund nếu đã charge (payment-service, Bước 7).
     - Xóa giỏ hàng (cart-service, Bước 4, 7).
     - Trả lỗi qua REST API (Bước 4, 7, mã lỗi 409, 400).
     - Ghi log vào AuditLogs và Logrus.

**Compensation Flow**:
- **inventory-service**: Unlock `locked_quantity`, xóa cache `inventory:{product_id}`, cập nhật `inventory_locked=false`.
- **payment-service**: Refund giao dịch, cập nhật `payment_processed=false`.
- **cart-service**: Xóa giỏ hàng.
- **notification-service**: Không xử lý thông báo lỗi (chỉ lưu thông báo trạng thái).
- **order-service**: Cập nhật `status=FAILED`, lưu `failure_reason` để theo dõi lý do thất bại.

---

## 4. 🗂️ Thiết kế Dữ liệu

### order-service Database (PostgreSQL)
- **Orders**: `id` (string, UUID), `customer_id` (string), `status` (string, enum: CREATED, CONFIRMED, DELIVERED, FAILED), `total_amount` (numeric), `shipping_address` (string), `created_at` (timestamp), `updated_at` (timestamp), `inventory_locked` (boolean), `payment_processed` (boolean), `shipping_scheduled` (boolean), `failure_reason` (text) (index: `id`, `customer_id`, retention: 1 năm)
- **OrderItems**: `id` (string, UUID), `order_id` (string), `product_id` (string), `quantity` (integer), `price` (numeric) (index: `order_id`, retention: 1 năm)
- **AuditLogs**: `id` (string, UUID), `service_name` (string), `action` (string, enum: create_order, process_payment, inventory_error, shipment_updated), `user_id` (string), `timestamp` (timestamp), `details` (string) (index: `user_id`, `timestamp`, `action`, retention: 6 tháng, max size: 1GB, xóa bản ghi cũ nhất khi vượt)

### user-service Database (PostgreSQL)
- **Customers**: `id` (string, UUID), `name` (string), `email` (string), `address` (string), `created_at` (timestamp), `updated_at` (timestamp) (index: `id`, `email`, retention: 1 năm)

### cart-service Database (PostgreSQL)
- **Carts**: `id` (string, UUID), `customer_id` (string), `created_at` (timestamp), `updated_at` (timestamp) (index: `customer_id`, retention: 90 ngày)
- **CartItems**: `id` (string, UUID), `cart_id` (string), `product_id` (string), `quantity` (integer) (index: `cart_id`, retention: 90 ngày)
- **Cache**: Redis key `cart:{user_id}` (TTL: 3600 giây, memory limit: 100MB), invalidation khi giỏ hàng cập nhật hoặc hết hạn (90 ngày).

### inventory-service Database (PostgreSQL)
- **Products**: `id` (string, UUID), `name` (string), `description` (string), `price` (numeric), `category` (string), `created_at` (timestamp) (index: `id`, `category`, retention: 1 năm)
- **ProductTags**: `product_id` (string), `tag` (string) (index: `product_id`, `tag`, retention: 1 năm)
- **Inventory**: `product_id` (string), `available_quantity` (integer), `locked_quantity` (integer), `updated_at` (timestamp) (index: `product_id`, retention: 1 năm)
- **ProductRelations**: `product_id` (string), `related_product_id` (string), `relation_type` (string, enum: category_match, tag_match), `score` (float, 0-1), `updated_at` (timestamp) (index: `product_id`, retention: 1 năm)
  - **Gợi ý sản phẩm**: Truy vấn ProductRelations theo `product_id`, ưu tiên `relation_type=category_match` và `score` cao nhất.
- **Cache**: Redis key `inventory:{product_id}` (TTL: 60 giây), `recommendation:{product_id}` (TTL: 3600 giây), memory limit: 100MB, invalidation khi Kafka event `order_confirmed` hoặc ProductRelations cập nhật.

### payment-service Database (PostgreSQL)
- **Payments**: `id` (string, UUID), `order_id` (string), `amount` (numeric), `status` (string, enum: SUCCEEDED, FAILED), `payment_method` (string), `transaction_id` (string), `created_at` (timestamp) (index: `order_id`, retention: 1 năm)

### shipping-service Database (PostgreSQL)
- **Shipments**: `id` (string, UUID), `order_id` (string), `carrier` (string), `tracking_number` (string), `status` (string, enum: SCHEDULED, IN_TRANSIT, DELIVERED), `shipping_address` (string), `created_at` (timestamp), `updated_at` (timestamp) (index: `order_id`, retention: 1 năm)
- **ShipmentUpdates**: `id` (string, UUID), `shipment_id` (string), `status` (string), `description` (string), `created_at` (timestamp) (index: `shipment_id`, retention: 6 tháng)
- **Cache**: Redis key `shipment:{order_id}` (TTL: 300 giây, memory limit: 100MB), invalidation khi Kafka event `shipping_status_updated`, `shipping_completed`.

### notification-service Database (PostgreSQL)
- **Notifications**: `id` (string, UUID), `user_id` (string), `type` (string, enum: ORDER_CONFIRMED, SHIPPING_UPDATED, SHIPPING_COMPLETED), `content` (string), `status` (string, enum: UNREAD, READ), `created_at` (timestamp) (index: `user_id`, `created_at`, retention: 6 tháng, max 100 notifications/user, xóa thông báo cũ nhất khi vượt)

**Data Synchronization**:
- **Kafka Events**: Đồng bộ qua Kafka (ví dụ: notification-service lưu thông báo sau `order_confirmed`).
- **Retention Policy**: Xóa dữ liệu cũ bằng cron job hàng tháng (ví dụ: `DELETE FROM Notifications WHERE created_at < NOW() - INTERVAL '6 months'`).
- **Backup**: Snapshot hàng ngày, nén bằng gzip, lưu trong AWS S3 (dung lượng tối đa 100GB, retention: 1 năm).
- **Restore Process**: Khôi phục từ S3 bằng script `pg_restore`, kiểm tra tính toàn vẹn dữ liệu trước khi áp dụng.

---

## 5. 🔐 Các vấn đề Bảo mật

- **Xác thực và Ủy quyền**:
  - JWT cho xác thực người dùng, kiểm tra tại api-gateway.
  - Role-based access control (RBAC): `customer` (đặt đơn), `admin` (quản lý).
- **Bảo mật Giao tiếp**:
  - HTTPS với JWT cho giao tiếp nội bộ và bên ngoài.
- **Bảo mật Dữ liệu**:
  - Mã hóa thông tin thanh toán (AES-256).
- **Bảo mật Kafka**:
  - Mã hóa sự kiện với SSL/TLS.
  - Xác thực producer/consumer với SASL/SCRAM.
- **Bảo mật Redis**:
  - Mã hóa dữ liệu với SSL/TLS.
- **Xác thực Đầu vào**:
  - Ngăn chặn SQL Injection, XSS bằng input validation.
- **Rate Limiting & CORS**:
  - Rate limit: 100 reqs/giây tại api-gateway, trả lỗi 429 (Too Many Requests).
  - CORS whitelist: `*.example.com`.
- **Quản lý Secrets**:
  - Lưu API keys, credentials trong Docker Compose secrets, xoay vòng hàng tháng.
- **Audit Logging**:
  - Ghi hành động nhạy cảm (create_order, process_payment, inventory_error, shipment_updated) vào bảng AuditLogs trong order-service.
- **Security Testing**:
  - Penetration testing hàng năm.
  - Vulnerability scanning hàng quý với OWASP ZAP, vá lỗ hổng trong 7 ngày sau khi phát hiện.

---

## 6. 📦 Kế hoạch Triển khai

- **Containerization**:
  - Mỗi dịch vụ chạy trong Docker container, quản lý bằng Docker Compose.
- **Auto-Scaling**:
  - Scaling thủ công bằng lệnh `docker-compose up --scale order-service=2`.
- **Yêu cầu Phần cứng**:
  - Server tối thiểu: 2 CPU, 4GB RAM, 50GB SSD.
- **Cost Optimization**:
  - Resource limits: CPU 0.5, memory 512Mi.
- **Cấu hình Môi trường**:
  - Biến môi trường trong Docker Compose.
- **CI/CD**:
  - GitHub Actions: build, deploy, rolling update với `docker-compose --no-cache` để tránh downtime.
- **Resilience & Observability**:
  - **Readiness/Liveness Probes**: Mỗi dịch vụ cung cấp endpoint `/health` để kiểm tra trạng thái.
  - **Retry**: Retry 2 lần cho API calls (trừ `/v1/shipments`: 3 lần, backoff 1 giây), Kafka consumer, và lưu thông báo vào Notifications.
  - **Structured Logging**: JSON logs với Logrus.
  - **Metrics**: Prometheus, giám sát CPU, memory, request latency sau khi scale.
- **Service Discovery**:
  - DNS-based discovery (Docker Compose DNS).
- **Caching**:
  - Redis 1 node, memory limit 100MB, invalidation khi Kafka event phát hoặc dữ liệu cập nhật.
- **Disaster Recovery**:
  - Backup database và Kafka snapshot hàng ngày, nén bằng gzip, lưu trong AWS S3 (dung lượng tối đa 100GB, retention: 1 năm).
  - Restore từ S3 bằng script `pg_restore`, kiểm tra tính toàn vẹn dữ liệu.

---

## 7. 🎨 Sơ đồ Kiến trúc

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

**Kafka Data Flow**:
- `order-events` → payment-service (Bước 6), order-service (Bước 8, 12), user-service (Bước 8, 13), notification-service (Bước 8)
- `payment-events` → order-service (Bước 8), shipping-service (Bước 9)
- `shipping-events` → order-service (Bước 12), notification-service (Bước 11, 13)

**Event Structures**:
- **OrderEvent**: `event_type`, `order_id`, `customer_id`, `status`, `total_amount`, `timestamp`, `items`, `failure_reason`
- **PaymentEvent**: `event_type`, `payment_id`, `order_id`, `customer_id`, `amount`, `status`, `payment_method`, `timestamp`
- **ShipmentEvent**: `event_type`, `shipment_id`, `order_id`, `status`, `timestamp`

**Compensation Flow Summary**:
- Nếu thanh toán thất bại hoặc hết hàng: order-service cập nhật `status=FAILED` và lưu `failure_reason`, trả lỗi qua REST API, inventory-service unlock tồn kho (Bước 4, xóa cache), payment-service hoàn tiền (Bước 7), cart-service xóa giỏ hàng, lỗi được ghi vào AuditLogs và Logrus.

---

## 8. ✅ Tóm tắt

Kiến trúc SOA này phù hợp với UC "Xử lý đơn hàng trực tuyến" vì:
1. **Tính Mô-đun**: 7 dịch vụ tách biệt trách nhiệm, hỗ trợ quy trình từ đặt hàng đến giao hàng.
2. **Tính Tự động**: Saga Orchestrator và Kafka đảm bảo quy trình thực hiện tự động, tuần tự qua 14 bước.
3. **Hiệu suất**: REST timeout (5 giây), Kafka partitioning (8), Redis caching (100MB), và resource limits tối ưu hóa hiệu suất.
4. **Khả năng Phục hồi**: Compensation flow xử lý lỗi (hết hàng, thanh toán thất bại).
5. **Khả năng Theo dõi**: Logrus và Prometheus hỗ trợ giám sát, AuditLogs ghi hành động nhạy cảm.
6. **Tính Linh hoạt**: Dễ tích hợp với Stripe, FedEx, UPS, và mở rộng với email (chưa triển khai).
7. **Khám phá Dịch vụ**: DNS-based discovery đáp ứng quy mô UC.

Hệ thống đáp ứng đầy đủ 14 bước nghiệp vụ, tuân thủ các nguyên lý SOA:
- **Standardized Service Contract**: Hợp đồng REST API chuẩn hóa.
- **Service Loose Coupling**: Dịch vụ độc lập, giao tiếp qua api-gateway/Kafka.
- **Service Abstraction**: DNS-based discovery và mã hóa che giấu chi tiết triển khai.
- **Service Reusability**: payment-service, user-service, notification-service tái sử dụng được.
- **Service Autonomy**: Database per Service (schema riêng).
- **Service Statelessness**: Kafka và REST API.
- **Service Discoverability**: DNS-based discovery.
- **Service Composability**: Saga và Kafka.