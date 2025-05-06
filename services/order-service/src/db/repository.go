package db

import (
"time"

"github.com/online-order-system/order-service/models"
)

// OrderRepository handles database operations for orders
type OrderRepository struct {
db *Database
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *Database) *OrderRepository {
return &OrderRepository{db: db}
}

// CreateAuditLog creates a new audit log entry in the database
func (r *OrderRepository) CreateAuditLog(log models.AuditLog) error {
	_, err := r.db.Exec(
		"INSERT INTO audit_logs (id, service_name, action, user_id, timestamp, details) VALUES ($1, $2, $3, $4, $5, $6)",
		log.ID, log.ServiceName, log.Action, log.CustomerID, log.Timestamp, log.Details,
	)
	return err
}

// CreateOrder creates a new order in the database
func (r *OrderRepository) CreateOrder(order models.Order) error {
// Begin transaction
tx, err := r.db.Begin()
if err != nil {
return err
}
defer tx.Rollback()

// Insert order
_, err = tx.Exec(
"INSERT INTO orders (id, user_id, status, total_amount, shipping_address, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
order.ID, order.CustomerID, order.Status, order.TotalAmount, order.ShippingAddress, order.CreatedAt, order.UpdatedAt,
)
if err != nil {
return err
}

// Insert order items
for _, item := range order.Items {
_, err = tx.Exec(
"INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4, $5)",
item.ID, order.ID, item.ProductID, item.Quantity, item.Price,
)
if err != nil {
return err
}
}

// Commit transaction
return tx.Commit()
}

// GetOrderByID retrieves an order by ID
func (r *OrderRepository) GetOrderByID(id string) (models.Order, error) {
var order models.Order
var createdAt, updatedAt time.Time
var status string


// Get order
err := r.db.QueryRow(
"SELECT id, user_id, status, total_amount, shipping_address, created_at, updated_at, inventory_locked, payment_processed, shipping_scheduled, failure_reason FROM orders WHERE id = $1",
id,
).Scan(&order.ID, &order.CustomerID, &status, &order.TotalAmount, &order.ShippingAddress, &createdAt, &updatedAt, &order.InventoryLocked, &order.PaymentProcessed, &order.ShippingScheduled, &order.FailureReason)
if err != nil {
return order, err
}

order.Status = models.OrderStatus(status)
order.CreatedAt = createdAt
order.UpdatedAt = updatedAt

// Get order items
rows, err := r.db.Query(
"SELECT id, product_id, quantity, price FROM order_items WHERE order_id = $1",
id,
)
if err != nil {
return order, err
}
defer rows.Close()

var items []models.OrderItem
for rows.Next() {
var item models.OrderItem
err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price)
if err != nil {
return order, err
}
item.OrderID = id
items = append(items, item)
}

order.Items = items
return order, nil
}

// GetOrders retrieves all orders
func (r *OrderRepository) GetOrders() ([]models.Order, error) {
// Get all orders
rows, err := r.db.Query(
"SELECT id, user_id, status, total_amount, shipping_address, created_at, updated_at FROM orders",
)
if err != nil {
return nil, err
}
defer rows.Close()

var orders []models.Order
for rows.Next() {
var order models.Order
var status string
var createdAt, updatedAt time.Time

err := rows.Scan(&order.ID, &order.CustomerID, &status, &order.TotalAmount, &order.ShippingAddress, &createdAt, &updatedAt)
if err != nil {
return nil, err
}

order.Status = models.OrderStatus(status)
order.CreatedAt = createdAt
order.UpdatedAt = updatedAt

// Get order items for this order
itemRows, err := r.db.Query(
"SELECT id, product_id, quantity, price FROM order_items WHERE order_id = $1",
order.ID,
)
if err != nil {
return nil, err
}
defer itemRows.Close()

var items []models.OrderItem
for itemRows.Next() {
var item models.OrderItem
err := itemRows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price)
if err != nil {
return nil, err
}
item.OrderID = order.ID
items = append(items, item)
}

order.Items = items
orders = append(orders, order)
}

return orders, nil
}

// UpdateOrderStatus updates the status of an order
func (r *OrderRepository) UpdateOrderStatus(id string, status models.OrderStatus) error {
_, err := r.db.Exec(
"UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3",
status, time.Now(), id,
)
return err
}

// UpdateOrder updates an order in the database
func (r *OrderRepository) UpdateOrder(order models.Order) error {
// Begin transaction
tx, err := r.db.Begin()
if err != nil {
return err
}
defer tx.Rollback()

// Update order
_, err = tx.Exec(
"UPDATE orders SET user_id = $1, status = $2, total_amount = $3, shipping_address = $4, updated_at = $5, inventory_locked = $6, payment_processed = $7, shipping_scheduled = $8, failure_reason = $9 WHERE id = $10",
order.CustomerID, order.Status, order.TotalAmount, order.ShippingAddress, order.UpdatedAt,
order.InventoryLocked, order.PaymentProcessed, order.ShippingScheduled, order.FailureReason, order.ID,
)
if err != nil {
return err
}

// Commit transaction
return tx.Commit()
}
