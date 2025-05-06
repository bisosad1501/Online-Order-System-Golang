package models

import (
"time"
)

// ShipmentStatus represents the status of a shipment
type ShipmentStatus string

// Shipment statuses
const (
ShipmentStatusPending   ShipmentStatus = "PENDING"
ShipmentStatusProcessing ShipmentStatus = "PROCESSING"
ShipmentStatusShipped   ShipmentStatus = "SHIPPED"
ShipmentStatusDelivered ShipmentStatus = "DELIVERED"
ShipmentStatusCancelled ShipmentStatus = "CANCELLED"
)

// Shipment represents a shipment in the system
type Shipment struct {
ID              string         `json:"id"`
OrderID         string         `json:"order_id"`
Status          ShipmentStatus `json:"status"`
TrackingNumber  string         `json:"tracking_number"`
ShippingAddress string         `json:"shipping_address"`
Carrier         string         `json:"carrier"`
EstimatedDelivery time.Time    `json:"estimated_delivery"`
CreatedAt       time.Time      `json:"created_at"`
UpdatedAt       time.Time      `json:"updated_at"`
CustomerID      string         `json:"customer_id,omitempty"` // Added for notification purposes
}

// CreateShipmentRequest represents a request to create a new shipment
type CreateShipmentRequest struct {
OrderID         string `json:"order_id" binding:"required"`
ShippingAddress string `json:"shipping_address" binding:"required"`
Carrier         string `json:"carrier,omitempty"`
CustomerID      string `json:"customer_id,omitempty"` // Added for notification purposes
}

// UpdateShipmentStatusRequest represents a request to update a shipment's status
type UpdateShipmentStatusRequest struct {
Status ShipmentStatus `json:"status" binding:"required"`
}

// UpdateTrackingRequest represents a request to update a shipment's tracking number
type UpdateTrackingRequest struct {
TrackingNumber string `json:"tracking_number" binding:"required"`
}

// ShipmentEvent represents an event related to a shipment
type ShipmentEvent struct {
EventType      string         `json:"event_type"`
ShipmentID     string         `json:"shipment_id"`
OrderID        string         `json:"order_id"`
Status         ShipmentStatus `json:"status"`
TrackingNumber string         `json:"tracking_number,omitempty"`
Timestamp      int64          `json:"timestamp"`
CustomerID     string         `json:"customer_id,omitempty"` // Added for notification purposes
}
