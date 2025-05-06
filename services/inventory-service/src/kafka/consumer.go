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

// Consumer represents a Kafka consumer
type Consumer struct {
reader  *kafka.Reader
service interfaces.InventoryService
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, service interfaces.InventoryService) *Consumer {
// Create a reader for the orders topic
ordersReader := kafka.NewReader(kafka.ReaderConfig{
Brokers:  []string{cfg.KafkaBootstrapServers},
Topic:    "orders", // Listen to orders topic
GroupID:  "inventory-service-orders",
MinBytes: 10e3, // 10KB
MaxBytes: 10e6, // 10MB
})

// Create a reader for the payments topic
paymentsReader := kafka.NewReader(kafka.ReaderConfig{
Brokers:  []string{cfg.KafkaBootstrapServers},
Topic:    "payments", // Listen to payments topic
GroupID:  "inventory-service-payments",
MinBytes: 10e3, // 10KB
MaxBytes: 10e6, // 10MB
})

log.Printf("Inventory Service Kafka consumer created and subscribed to topics: orders, payments")

// Start consuming from both topics
consumer := &Consumer{
reader:  ordersReader, // We'll use this as the primary reader
service: service,
}

// Start a separate goroutine to consume from the payments topic
go func() {
ctx := context.Background()
for {
m, err := paymentsReader.ReadMessage(ctx)
if err != nil {
log.Printf("Error reading message from Kafka payments topic: %v", err)
time.Sleep(1 * time.Second)
continue
}

log.Printf("Received message from Kafka payments topic: %s", string(m.Value))

// Process the message
var event struct {
EventType string `json:"event_type"`
}
if err := json.Unmarshal(m.Value, &event); err != nil {
log.Printf("Error unmarshaling message from payments topic: %v", err)
continue
}

// Process the event based on its type
if event.EventType == "payment_failed" {
log.Printf("Processing payment failed event from payments topic")
var paymentEvent struct {
OrderID string `json:"order_id"`
}
if err := json.Unmarshal(m.Value, &paymentEvent); err != nil {
log.Printf("Error unmarshaling payment failed event: %v", err)
continue
}

// Get order details to restore inventory
// We need to make a request to the order service to get the order details
// For now, we'll just log that we received the event
log.Printf("Payment failed for order %s, should restore inventory", paymentEvent.OrderID)

// TODO: Implement inventory restoration logic
// This would require getting the order details from the order service
// and then updating the inventory for each item in the order
}
}
}()

return consumer
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
if event.EventType == "order_created" {
log.Printf("Processing order created event")
var orderEvent struct {
OrderID string `json:"order_id"`
Items   []struct {
ProductID string `json:"product_id"`
Quantity  int    `json:"quantity"`
} `json:"items"`
}
if err := json.Unmarshal(m.Value, &orderEvent); err != nil {
log.Printf("Error unmarshaling order created event: %v", err)
continue
}

// Update inventory for each item
for _, item := range orderEvent.Items {
// Get current product with inventory
product, err := c.service.GetProductByID(item.ProductID)
if err != nil {
log.Printf("Error getting product %s: %v", item.ProductID, err)
continue
}

// Update inventory
newQuantity := product.Quantity - item.Quantity
if newQuantity < 0 {
newQuantity = 0
}

// Create update request with only quantity
quantityPtr := new(int)
*quantityPtr = newQuantity

// UpdateProduct returns (models.Product, error)
_, err = c.service.UpdateProduct(item.ProductID, models.UpdateProductRequest{
Quantity: quantityPtr,
})
if err != nil {
log.Printf("Error updating inventory for product %s: %v", item.ProductID, err)
}
}
}
}
}
}()
}
