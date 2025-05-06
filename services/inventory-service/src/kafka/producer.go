package kafka

import (
"context"
"encoding/json"
"log"
"time"

"github.com/segmentio/kafka-go"
"github.com/online-order-system/inventory-service/config"
"github.com/online-order-system/inventory-service/interfaces"
"github.com/online-order-system/inventory-service/models"
)

// Producer represents a Kafka producer
type Producer struct {
writer *kafka.Writer
}

// Ensure Producer implements InventoryProducer interface
var _ interfaces.InventoryProducer = (*Producer)(nil)

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

// PublishInventoryUpdated publishes an inventory updated event
func (p *Producer) PublishInventoryUpdated(productID string, quantity int) error {
event := models.InventoryEvent{
EventType:  "inventory_updated",
ProductID:  productID,
Quantity:   quantity,
Timestamp:  time.Now().Unix(),
}

return p.publishEvent(event)
}

// publishEvent publishes an event to Kafka
func (p *Producer) publishEvent(event models.InventoryEvent) error {
// Convert event to JSON
eventJSON, err := json.Marshal(event)
if err != nil {
return err
}

// Create message
message := kafka.Message{
Key:   []byte(event.ProductID),
Value: eventJSON,
}

// Write message
err = p.writer.WriteMessages(context.Background(), message)
if err != nil {
return err
}

log.Printf("Published event: %s for product %s", event.EventType, event.ProductID)
return nil
}
