package models

import (
"time"
)

// NotificationType represents the type of a notification
type NotificationType string

// Notification types
const (
	// As per design, only these three types are defined
	NotificationTypeOrderConfirmed NotificationType = "ORDER_CONFIRMED"
	NotificationTypeShippingUpdate NotificationType = "SHIPPING_UPDATED"
	NotificationTypeShippingDone   NotificationType = "SHIPPING_COMPLETED"

	// Additional types for comprehensive notification handling
	NotificationTypeOrderCompleted NotificationType = "ORDER_COMPLETED"
	NotificationTypeOrderCancelled NotificationType = "ORDER_CANCELLED"
	NotificationTypePaymentFailed  NotificationType = "PAYMENT_FAILED"

	// Channel types for internal use
	NotificationTypeEmail   NotificationType = "EMAIL"
	NotificationTypeSMS     NotificationType = "SMS"
	NotificationTypePush    NotificationType = "PUSH"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

// Notification statuses
const (
// As per design, only these two types are defined
NotificationStatusUnread   NotificationStatus = "UNREAD"
NotificationStatusRead     NotificationStatus = "READ"

// Additional statuses for internal use
NotificationStatusPending   NotificationStatus = "PENDING"
NotificationStatusSent      NotificationStatus = "SENT"
NotificationStatusFailed    NotificationStatus = "FAILED"
)

// Notification represents a notification in the system
type Notification struct {
ID          string             `json:"id"`
CustomerID  string             `json:"customer_id"` // Kept as CustomerID for code consistency, maps to user_id in DB
Type        NotificationType   `json:"type"`
Status      NotificationStatus `json:"status"`
Subject     string             `json:"subject,omitempty"` // Not in DB schema as per design
Content     string             `json:"content"`
Recipient   string             `json:"recipient,omitempty"` // Not in DB schema as per design
CreatedAt   time.Time          `json:"created_at"`
UpdatedAt   time.Time          `json:"updated_at"`
SentAt      *time.Time         `json:"sent_at,omitempty"`
}

// CreateNotificationRequest represents a request to create a new notification
type CreateNotificationRequest struct {
CustomerID  string           `json:"customer_id" binding:"required"`
Type        NotificationType `json:"type" binding:"required"`
Subject     string           `json:"subject" binding:"required"`
Content     string           `json:"content" binding:"required"`
Recipient   string           `json:"recipient" binding:"required"`
}

// UpdateNotificationStatusRequest represents a request to update a notification's status
type UpdateNotificationStatusRequest struct {
Status NotificationStatus `json:"status" binding:"required"`
}

// NotificationEvent represents an event related to a notification
type NotificationEvent struct {
EventType   string             `json:"event_type"`
CustomerID  string             `json:"customer_id"`
Type        NotificationType   `json:"type"`
Status      NotificationStatus `json:"status"`
Subject     string             `json:"subject"`
Content     string             `json:"content"`
Recipient   string             `json:"recipient"`
Timestamp   int64              `json:"timestamp"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	ProductID   string  `json:"product_id"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// OrderEvent represents an order event from Kafka
type OrderEvent struct {
	EventType       string      `json:"event_type"`
	OrderID         string      `json:"order_id"`
	CustomerID      string      `json:"customer_id"`
	Status          string      `json:"status"`
	TotalAmount     float64     `json:"total_amount"`
	Timestamp       int64       `json:"timestamp"`
	FailureReason   string      `json:"failure_reason,omitempty"`
	Items           []OrderItem `json:"items,omitempty"`
	ShippingAddress string      `json:"shipping_address,omitempty"`
}

// PaymentEvent represents a payment event from Kafka
type PaymentEvent struct {
EventType     string  `json:"event_type"`
PaymentID     string  `json:"payment_id"`
OrderID       string  `json:"order_id"`
CustomerID    string  `json:"customer_id"`
Amount        float64 `json:"amount"`
Status        string  `json:"status"`
PaymentMethod string  `json:"payment_method"`
Timestamp     int64   `json:"timestamp"`
}

// ShipmentEvent represents a shipment event from Kafka
type ShipmentEvent struct {
EventType      string `json:"event_type"`
ShipmentID     string `json:"shipment_id"`
OrderID        string `json:"order_id"`
Status         string `json:"status"`
TrackingNumber string `json:"tracking_number,omitempty"`
Timestamp      int64  `json:"timestamp"`
CustomerID     string `json:"customer_id,omitempty"` // Added for notification purposes
}
