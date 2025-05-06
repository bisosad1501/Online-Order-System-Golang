package interfaces

import (
"github.com/online-order-system/shipping-service/models"
)

// ShippingService defines the interface for shipping service
type ShippingService interface {
CreateShipment(req models.CreateShipmentRequest) (models.Shipment, error)
GetShipmentByID(id string) (models.Shipment, error)
GetShipmentByOrderID(orderID string) (models.Shipment, error)
GetShipments() ([]models.Shipment, error)
UpdateShipmentStatus(id string, status models.ShipmentStatus) (models.Shipment, error)
UpdateTrackingNumber(id string, trackingNumber string) (models.Shipment, error)
}

// ShippingProducer defines the interface for shipping producer
type ShippingProducer interface {
PublishShipmentCreated(shipment models.Shipment) error
PublishShipmentStatusUpdated(shipment models.Shipment) error
PublishShipmentCompleted(shipment models.Shipment) error
Close() error
}
