package models

import (
	"time"
)

// Customer represents a customer in the system
type Customer struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CustomerOrder represents an order associated with a customer
type CustomerOrder struct {
	CustomerID  string    `json:"customer_id"`
	OrderID     string    `json:"order_id"`
	OrderDate   time.Time `json:"order_date"`
	OrderStatus string    `json:"order_status"`
}

// CreateCustomerRequest represents a request to create a new customer
type CreateCustomerRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

// UpdateCustomerRequest represents a request to update a customer
type UpdateCustomerRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

// VerifyCustomerRequest represents a request to verify a customer
type VerifyCustomerRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Address string `json:"address" binding:"omitempty"`
	ID      string `json:"id" binding:"omitempty"`
}

// VerifyCustomerResponse represents a response from customer verification
type VerifyCustomerResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message,omitempty"`
	Token    string `json:"token,omitempty"`
	User     *Customer `json:"user,omitempty"`
}

// LoginRequest represents a request to login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty"`
}

// LoginResponse represents a response from login
type LoginResponse struct {
	Token string    `json:"token"`
	User  *Customer `json:"user"`
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// CustomerEvent represents an event related to a customer
type CustomerEvent struct {
	EventType  string `json:"event_type"`
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Timestamp  int64  `json:"timestamp"`
}

// OrderEvent represents an event related to an order
type OrderEvent struct {
	EventType      string  `json:"event_type"`
	OrderID        string  `json:"order_id"`
	CustomerID     string  `json:"customer_id"`
	Status         string  `json:"status"`
	TotalAmount    float64 `json:"total_amount"`
	TrackingNumber string  `json:"tracking_number,omitempty"`
	Timestamp      int64   `json:"timestamp"`
}
