package models

import (
"time"
)

// Cart represents a shopping cart in the system
type Cart struct {
ID         string     `json:"id"`
CustomerID string     `json:"customer_id"`
Items      []CartItem `json:"items,omitempty"`
CreatedAt  time.Time  `json:"created_at"`
UpdatedAt  time.Time  `json:"updated_at"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
ID        string    `json:"id"`
CartID    string    `json:"cart_id"`
ProductID string    `json:"product_id"`
Quantity  int       `json:"quantity"`
Price     float64   `json:"price"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}

// CreateCartRequest represents a request to create a new cart
type CreateCartRequest struct {
CustomerID string `json:"customer_id" binding:"required"`
}

// AddCartItemRequest represents a request to add an item to a cart
type AddCartItemRequest struct {
ProductID string  `json:"product_id" binding:"required"`
Quantity  int     `json:"quantity" binding:"required,min=1"`
Price     float64 `json:"price" binding:"required,min=0"`
}

// UpdateCartItemRequest represents a request to update an item in a cart
type UpdateCartItemRequest struct {
Quantity int `json:"quantity" binding:"required,min=1"`
}

// CartEvent represents an event related to a cart
type CartEvent struct {
EventType  string     `json:"event_type"`
CartID     string     `json:"cart_id"`
CustomerID string     `json:"customer_id"`
Items      []CartItem `json:"items,omitempty"`
Timestamp  int64      `json:"timestamp"`
}

// OrderEvent represents an event related to an order
type OrderEvent struct {
EventType   string    `json:"event_type"`
OrderID     string    `json:"order_id"`
CustomerID  string    `json:"customer_id"`
Status      string    `json:"status"`
TotalAmount float64   `json:"total_amount"`
Timestamp   int64     `json:"timestamp"`
Items       []CartItem `json:"items,omitempty"`
}
