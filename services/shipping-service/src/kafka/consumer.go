package kafka

import (
"context"
"encoding/json"
"log"
"time"

"github.com/segmentio/kafka-go"
"github.com/online-order-system/shipping-service/config"
"github.com/online-order-system/shipping-service/interfaces"
"github.com/online-order-system/shipping-service/models"
)

// Consumer represents a Kafka consumer
type Consumer struct {
reader  *kafka.Reader
service interfaces.ShippingService
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, service interfaces.ShippingService) *Consumer {
reader := kafka.NewReader(kafka.ReaderConfig{
Brokers:  []string{cfg.KafkaBootstrapServers},
Topic:    "orders", // Listen to orders topic
GroupID:  "shipping-service",
MinBytes: 10e3, // 10KB
MaxBytes: 10e6, // 10MB
})

log.Printf("Shipping Service Kafka consumer created and subscribed to topic: orders")

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
m, err := c.reader.ReadMessage(ctx)
if err != nil {
log.Printf("Error reading message from Kafka: %v", err)
time.Sleep(1 * time.Second)
continue
}

log.Printf("Received message from Kafka: %s", string(m.Value))

// Process the message
var event struct {
EventType string `json:"event_type"`
}
if err := json.Unmarshal(m.Value, &event); err != nil {
log.Printf("Error unmarshaling message: %v", err)
continue
}

// Process the event based on its type
switch event.EventType {
case "order_confirmed":
log.Printf("Processing order confirmed event")
var orderEvent struct {
OrderID         string `json:"order_id"`
ShippingAddress string `json:"shipping_address"`
}
if err := json.Unmarshal(m.Value, &orderEvent); err != nil {
log.Printf("Error unmarshaling order confirmed event: %v", err)
continue
}

// Create shipment
_, err = c.service.CreateShipment(models.CreateShipmentRequest{
OrderID:         orderEvent.OrderID,
ShippingAddress: orderEvent.ShippingAddress,
})
if err != nil {
log.Printf("Error creating shipment for order %s: %v", orderEvent.OrderID, err)
}
}
}
}
}()
}
