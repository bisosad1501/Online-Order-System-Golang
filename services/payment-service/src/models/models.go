package models

import (
"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

// Payment statuses
const (
PaymentStatusPending   PaymentStatus = "PENDING"
PaymentStatusSuccessful PaymentStatus = "SUCCESSFUL"
PaymentStatusFailed    PaymentStatus = "FAILED"
PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

// Payment represents a payment in the system
type Payment struct {
ID                string        `json:"id"`
OrderID           string        `json:"order_id"`
Amount            float64       `json:"amount"`
Status            PaymentStatus `json:"status"`
PaymentMethod     string        `json:"payment_method"`
CardNumber        string        `json:"card_number,omitempty"`
ExpiryMonth       string        `json:"expiry_month,omitempty"`
ExpiryYear        string        `json:"expiry_year,omitempty"`
CVV               string        `json:"cvv,omitempty"`
CreatedAt         time.Time     `json:"created_at"`
UpdatedAt         time.Time     `json:"updated_at"`
// Stripe specific fields
StripePaymentID   string        `json:"stripe_payment_id,omitempty"`
StripeClientSecret string       `json:"stripe_client_secret,omitempty"`
Currency          string        `json:"currency,omitempty"`
Description       string        `json:"description,omitempty"`
CustomerEmail     string        `json:"customer_email,omitempty"`
CustomerName      string        `json:"customer_name,omitempty"`
PaymentMethodID   string        `json:"payment_method_id,omitempty"`
ReceiptURL        string        `json:"receipt_url,omitempty"`
ErrorMessage      string        `json:"error_message,omitempty"`
}

// CreatePaymentRequest represents a request to create a new payment
type CreatePaymentRequest struct {
OrderID       string  `json:"order_id" binding:"required"`
Amount        float64 `json:"amount" binding:"required"`
PaymentMethod string  `json:"payment_method" binding:"required"`
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

// UpdatePaymentStatusRequest represents a request to update a payment's status
type UpdatePaymentStatusRequest struct {
Status PaymentStatus `json:"status" binding:"required"`
}

// PaymentEvent represents an event related to a payment
type PaymentEvent struct {
EventType     string        `json:"event_type"`
PaymentID     string        `json:"payment_id"`
OrderID       string        `json:"order_id"`
Amount        float64       `json:"amount"`
Status        PaymentStatus `json:"status"`
PaymentMethod string        `json:"payment_method"`
Timestamp     int64         `json:"timestamp"`
}
