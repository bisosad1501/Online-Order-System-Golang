package service

import (
"fmt"
"log"
"time"

"github.com/online-order-system/shipping-service/config"
"github.com/online-order-system/shipping-service/db"
"github.com/online-order-system/shipping-service/interfaces"
"github.com/online-order-system/shipping-service/models"
)

// ShippingService handles business logic for shipments
type ShippingService struct {
config     *config.Config
repository *db.ShippingRepository
producer   interfaces.ShippingProducer
orderClient *OrderClient
}

// Ensure ShippingService implements ShippingService interface
var _ interfaces.ShippingService = (*ShippingService)(nil)

// NewShippingService creates a new shipping service
func NewShippingService(cfg *config.Config, repo *db.ShippingRepository, producer interfaces.ShippingProducer) *ShippingService {
return &ShippingService{
config:     cfg,
repository: repo,
producer:   producer,
orderClient: NewOrderClient(cfg),
}
}

// CreateShipment creates a new shipment
func (s *ShippingService) CreateShipment(req models.CreateShipmentRequest) (models.Shipment, error) {
// Create shipment
now := time.Now()
estimatedDelivery := now.Add(3 * 24 * time.Hour) // 3 days from now

// Use default carrier if not provided
carrier := req.Carrier
if carrier == "" {
carrier = "Standard Shipping"
}

// Generate tracking number
trackingNumber := fmt.Sprintf("TRK-%s", db.GenerateID()[:8])

// Get customer ID from order service if not provided
customerID := req.CustomerID
if customerID == "" {
    // Get order details to get customer ID
    order, err := s.orderClient.GetOrderByID(req.OrderID)
    if err != nil {
        log.Printf("Failed to get order details: %v", err)
        // Continue without customer ID
    } else {
        customerID = order.CustomerID
    }
}

shipment := models.Shipment{
ID:               db.GenerateID(),
OrderID:          req.OrderID,
Status:           models.ShipmentStatusPending,
TrackingNumber:   trackingNumber,
ShippingAddress:  req.ShippingAddress,
Carrier:          carrier,
EstimatedDelivery: estimatedDelivery,
CreatedAt:        now,
UpdatedAt:        now,
CustomerID:       customerID, // Save customer ID
}

// Save shipment to database
err := s.repository.CreateShipment(shipment)
if err != nil {
return models.Shipment{}, err
}

// Publish shipment created event
err = s.producer.PublishShipmentCreated(shipment)
if err != nil {
log.Printf("Failed to publish shipment created event: %v", err)
// Continue anyway for demo purposes
}

return shipment, nil
}

// GetShipmentByID retrieves a shipment by ID
func (s *ShippingService) GetShipmentByID(id string) (models.Shipment, error) {
return s.repository.GetShipmentByID(id)
}

// GetShipmentByOrderID retrieves a shipment by order ID
func (s *ShippingService) GetShipmentByOrderID(orderID string) (models.Shipment, error) {
return s.repository.GetShipmentByOrderID(orderID)
}

// GetShipments retrieves all shipments
func (s *ShippingService) GetShipments() ([]models.Shipment, error) {
return s.repository.GetShipments()
}

// UpdateShipmentStatus updates the status of a shipment
func (s *ShippingService) UpdateShipmentStatus(id string, status models.ShipmentStatus) (models.Shipment, error) {
// Update shipment status
err := s.repository.UpdateShipmentStatus(id, status)
if err != nil {
return models.Shipment{}, err
}

// Get updated shipment
shipment, err := s.repository.GetShipmentByID(id)
if err != nil {
return models.Shipment{}, err
}

// Publish event based on status
var publishErr error
if status == models.ShipmentStatusDelivered {
publishErr = s.producer.PublishShipmentCompleted(shipment)
} else {
publishErr = s.producer.PublishShipmentStatusUpdated(shipment)
}

if publishErr != nil {
log.Printf("Failed to publish shipment event: %v", publishErr)
// Continue anyway for demo purposes
}

return shipment, nil
}

// UpdateTrackingNumber updates the tracking number of a shipment
func (s *ShippingService) UpdateTrackingNumber(id string, trackingNumber string) (models.Shipment, error) {
// Update tracking number
err := s.repository.UpdateTrackingNumber(id, trackingNumber)
if err != nil {
return models.Shipment{}, err
}

// Get updated shipment
shipment, err := s.repository.GetShipmentByID(id)
if err != nil {
return models.Shipment{}, err
}

// Publish shipment status updated event
err = s.producer.PublishShipmentStatusUpdated(shipment)
if err != nil {
log.Printf("Failed to publish shipment status updated event: %v", err)
// Continue anyway for demo purposes
}

return shipment, nil
}
