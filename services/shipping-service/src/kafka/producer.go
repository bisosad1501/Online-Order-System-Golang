package kafka

import (
"context"
"encoding/json"
"log"

"github.com/segmentio/kafka-go"
"github.com/online-order-system/shipping-service/config"
"github.com/online-order-system/shipping-service/interfaces"
"github.com/online-order-system/shipping-service/models"
)

// Producer represents a Kafka producer
type Producer struct {
writer *kafka.Writer
}

// Ensure Producer implements ShippingProducer interface
var _ interfaces.ShippingProducer = (*Producer)(nil)

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.Config) *Producer {
writer := &kafka.Writer{
Addr:     kafka.TCP(cfg.KafkaBootstrapServers),
Topic:    cfg.KafkaTopic,
Balancer: &kafka.LeastBytes{},
}

return &Producer{
writer: writer,
}
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
return p.writer.Close()
}

// PublishShipmentCreated publishes a shipment created event
func (p *Producer) PublishShipmentCreated(shipment models.Shipment) error {
event := models.ShipmentEvent{
EventType:      "shipment_created",
ShipmentID:     shipment.ID,
OrderID:        shipment.OrderID,
Status:         shipment.Status,
TrackingNumber: shipment.TrackingNumber,
Timestamp:      shipment.CreatedAt.Unix(),
CustomerID:     shipment.CustomerID, // Include customer ID
}

return p.publishEvent(event)
}

// PublishShipmentStatusUpdated publishes a shipment status updated event
func (p *Producer) PublishShipmentStatusUpdated(shipment models.Shipment) error {
event := models.ShipmentEvent{
EventType:      "shipment_status_updated",
ShipmentID:     shipment.ID,
OrderID:        shipment.OrderID,
Status:         shipment.Status,
TrackingNumber: shipment.TrackingNumber,
Timestamp:      shipment.UpdatedAt.Unix(),
CustomerID:     shipment.CustomerID, // Include customer ID
}

return p.publishEvent(event)
}

// PublishShipmentCompleted publishes a shipment completed event
func (p *Producer) PublishShipmentCompleted(shipment models.Shipment) error {
event := models.ShipmentEvent{
EventType:      "shipping_completed",
ShipmentID:     shipment.ID,
OrderID:        shipment.OrderID,
Status:         shipment.Status,
TrackingNumber: shipment.TrackingNumber,
Timestamp:      shipment.UpdatedAt.Unix(),
CustomerID:     shipment.CustomerID, // Include customer ID
}

return p.publishEvent(event)
}

// publishEvent publishes an event to Kafka
func (p *Producer) publishEvent(event models.ShipmentEvent) error {
// Convert event to JSON
eventJSON, err := json.Marshal(event)
if err != nil {
return err
}

// Create message
message := kafka.Message{
Key:   []byte(event.ShipmentID),
Value: eventJSON,
}

// Write message
err = p.writer.WriteMessages(context.Background(), message)
if err != nil {
return err
}

log.Printf("Published event: %s for shipment %s", event.EventType, event.ShipmentID)
return nil
}
