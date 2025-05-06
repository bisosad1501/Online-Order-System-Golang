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

// Consumer represents a Kafka consumer
type Consumer struct {
reader  *kafka.Reader
service interfaces.CartService
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, service interfaces.CartService) *Consumer {
reader := kafka.NewReader(kafka.ReaderConfig{
Brokers:         []string{cfg.KafkaBootstrapServers},
Topic:           "orders", // Listen to orders topic
GroupID:         "cart-service",
MinBytes:        10e3, // 10KB
MaxBytes:        10e6, // 10MB
MaxWait:         1 * time.Second,
ReadLagInterval: -1,
})

return &Consumer{
reader:  reader,
service: service,
}
}

// StartConsuming starts consuming messages from Kafka
func (c *Consumer) StartConsuming(ctx context.Context) {
go func() {
for {
select {
case <-ctx.Done():
log.Println("Stopping Kafka consumer")
c.reader.Close()
return
default:
// Read message
msg, err := c.reader.ReadMessage(ctx)
if err != nil {
log.Printf("Error reading message: %v", err)
continue
}

// Process message
err = c.processMessage(msg)
if err != nil {
log.Printf("Error processing message: %v", err)
continue
}
}
}
}()

log.Println("Kafka consumer started")
}

// processMessage processes a Kafka message
func (c *Consumer) processMessage(msg kafka.Message) error {
// Parse message
var event models.OrderEvent
err := json.Unmarshal(msg.Value, &event)
if err != nil {
return err
}

log.Printf("Received message from Kafka orders topic: %s", string(msg.Value))

// Process event
return c.processOrderEvent(event)
}

// processOrderEvent processes an order event
func (c *Consumer) processOrderEvent(event models.OrderEvent) error {
// Delete cart when order is created
if event.EventType == "order_created" {
log.Printf("Deleting cart for user %s after order %s was created", event.CustomerID, event.OrderID)
return c.service.DeleteCartByUserID(event.CustomerID)
}
return nil
}
