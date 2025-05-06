package models

import (
"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

// Order statuses
const (
OrderStatusCreated   OrderStatus = "CREATED"
OrderStatusConfirmed OrderStatus = "CONFIRMED"
OrderStatusDelivered OrderStatus = "DELIVERED"
OrderStatusCancelled OrderStatus = "CANCELLED"
// Additional statuses for internal use
OrderStatusFailed    OrderStatus = "FAILED"
)

// Order represents an order in the system
type Order struct {
ID                string       `json:"id"`
CustomerID        string       `json:"customer_id"`
Status            OrderStatus  `json:"status"`
TotalAmount       float64      `json:"total_amount"`
ShippingAddress   string       `json:"shipping_address"`
Items             []OrderItem  `json:"items,omitempty"`
CreatedAt         time.Time    `json:"created_at"`
UpdatedAt         time.Time    `json:"updated_at"`
InventoryLocked   bool         `json:"inventory_locked,omitempty"`
PaymentProcessed  bool         `json:"payment_processed,omitempty"`
ShippingScheduled bool         `json:"shipping_scheduled,omitempty"`
FailureReason     string       `json:"failure_reason,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
ID        string  `json:"id"`
OrderID   string  `json:"order_id"`
ProductID string  `json:"product_id"`
Quantity  int     `json:"quantity"`
Price     float64 `json:"price"`
}

// CreateOrderRequest represents a request to create a new order
type CreateOrderRequest struct {
CustomerID      string      `json:"customer_id" binding:"required"`
Items           []OrderItem `json:"items" binding:"required"`
ShippingAddress string      `json:"shipping_address" binding:"required"`
}

// UpdateOrderStatusRequest represents a request to update an order's status
type UpdateOrderStatusRequest struct {
Status OrderStatus `json:"status" binding:"required"`
}

// OrderEvent represents an event related to an order
type OrderEvent struct {
	EventType       string      `json:"event_type"`
	OrderID         string      `json:"order_id"`
	CustomerID      string      `json:"customer_id"`
	Status          OrderStatus `json:"status"`
	TotalAmount     float64     `json:"total_amount"`
	Timestamp       int64       `json:"timestamp"`
	Items           []OrderItem `json:"items,omitempty"`
	FailureReason   string      `json:"failure_reason,omitempty"`
	ShippingAddress string      `json:"shipping_address,omitempty"`
}

// InventoryCheckRequest represents a request to check inventory
type InventoryCheckRequest struct {
Items []struct {
ProductID string `json:"product_id"`
Quantity  int    `json:"quantity"`
} `json:"items"`
}

// InventoryCheckResponse represents a response from inventory check
type InventoryCheckResponse struct {
Available       bool `json:"available"`
UnavailableItems []struct {
ProductID   string `json:"product_id"`
ProductName string `json:"product_name"`
Requested   int    `json:"requested"`
Available   int    `json:"available"`
} `json:"unavailable_items,omitempty"`
}

// InventoryRestoreRequest represents a request to restore inventory
type InventoryRestoreRequest struct {
ProductID string `json:"product_id"`
Quantity  int    `json:"quantity"`
}

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
OrderID       string  `json:"order_id"`
Amount        float64 `json:"amount"`
PaymentMethod string  `json:"payment_method"`
CardNumber    string  `json:"card_number,omitempty"`
ExpiryMonth   string  `json:"expiry_month,omitempty"`
ExpiryYear    string  `json:"expiry_year,omitempty"`
CVV           string  `json:"cvv,omitempty"`
// Stripe specific fields
Currency      string  `json:"currency,omitempty"`
Description   string  `json:"description,omitempty"`
CustomerEmail string  `json:"customer_email,omitempty"`
CustomerName  string  `json:"customer_name,omitempty"`
// For Stripe payment method
PaymentMethodID string `json:"payment_method_id,omitempty"`
}

// RefundRequest represents a request to refund a payment
type RefundRequest struct {
OrderID string  `json:"order_id"`
Amount  float64 `json:"amount"`
}

// CreateShipmentRequest represents a request to create a shipment
type CreateShipmentRequest struct {
OrderID         string `json:"order_id"`
ShippingAddress string `json:"shipping_address"`
Carrier         string `json:"carrier,omitempty"`
}

// RecommendationResponse represents a response from recommendation service
type RecommendationResponse struct {
Products []Product `json:"products"`
}

// Product represents a product in the system
type Product struct {
ID          string    `json:"id"`
Name        string    `json:"name"`
Description string    `json:"description"`
CategoryID  string    `json:"category_id"`
Price       float64   `json:"price"`
Tags        []string  `json:"tags,omitempty"`
CreatedAt   time.Time `json:"created_at"`
UpdatedAt   time.Time `json:"updated_at"`
}

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
CustomerID string `json:"customer_id"`
Type       string `json:"type"`
Channel    string `json:"channel"`
Subject    string `json:"subject"`
Content    string `json:"content"`
Recipient  string `json:"recipient"`
}

// RetryPaymentRequest represents a request to retry payment for a failed order
type RetryPaymentRequest struct {
PaymentMethod string `json:"payment_method" binding:"required"`
CardNumber    string `json:"card_number,omitempty"`
ExpiryMonth   string `json:"expiry_month,omitempty"`
ExpiryYear    string `json:"expiry_year,omitempty"`
CVV           string `json:"cvv,omitempty"`
}

// AuditLogAction represents the action type for audit logs
type AuditLogAction string

// Audit log actions
const (
AuditLogActionCreateOrder    AuditLogAction = "create_order"
AuditLogActionProcessPayment AuditLogAction = "process_payment"
AuditLogActionInventoryError AuditLogAction = "inventory_error"
AuditLogActionShipmentUpdate AuditLogAction = "shipment_updated"
)

// AuditLog represents an audit log entry in the system
type AuditLog struct {
ID          string        `json:"id"`
ServiceName string        `json:"service_name"`
Action      AuditLogAction `json:"action"`
CustomerID  string        `json:"customer_id"`
Timestamp   time.Time     `json:"timestamp"`
Details     string        `json:"details"`
}
