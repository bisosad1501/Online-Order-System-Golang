package kafka

import (
"context"
"encoding/json"
"log"
"time"

"github.com/segmentio/kafka-go"
"github.com/online-order-system/cart-service/config"
"github.com/online-order-system/cart-service/interfaces"
"github.com/online-order-system/cart-service/models"
)

// Producer represents a Kafka producer
type Producer struct {
writer *kafka.Writer
}

// Ensure Producer implements CartProducer interface
var _ interfaces.CartProducer = (*Producer)(nil)

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.Config) *Producer {
writer := kafka.NewWriter(kafka.WriterConfig{
Brokers:  []string{cfg.KafkaBootstrapServers},
Topic:    cfg.KafkaTopic,
Balancer: &kafka.LeastBytes{},
})

return &Producer{
writer: writer,
}
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
return p.writer.Close()
}

// PublishCartUpdated publishes a cart updated event
func (p *Producer) PublishCartUpdated(cart models.Cart) error {
event := models.CartEvent{
EventType:  "cart_updated",
CartID:     cart.ID,
CustomerID: cart.CustomerID,
Items:      cart.Items,
Timestamp:  time.Now().Unix(),
}

return p.publishEvent(event)
}

// PublishCartCleared publishes a cart cleared event
func (p *Producer) PublishCartCleared(cartID string, userID string) error {
event := models.CartEvent{
EventType:  "cart_cleared",
CartID:     cartID,
CustomerID: userID,
Timestamp:  time.Now().Unix(),
}

return p.publishEvent(event)
}

// publishEvent publishes an event to Kafka
func (p *Producer) publishEvent(event models.CartEvent) error {
// Marshal event to JSON
eventJSON, err := json.Marshal(event)
if err != nil {
return err
}

// Create message
msg := kafka.Message{
Key:   []byte(event.CartID),
Value: eventJSON,
Time:  time.Now(),
}

// Write message to Kafka
err = p.writer.WriteMessages(context.Background(), msg)
if err != nil {
return err
}

log.Printf("Published event: %s for cart %s", event.EventType, event.CartID)
return nil
}
