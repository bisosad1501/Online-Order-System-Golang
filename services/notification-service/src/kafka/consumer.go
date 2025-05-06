package kafka

import (
"context"
"encoding/json"
"log"
"time"

"github.com/segmentio/kafka-go"
"github.com/online-order-system/notification-service/config"
"github.com/online-order-system/notification-service/interfaces"
"github.com/online-order-system/notification-service/models"
)

// Consumer represents a Kafka consumer
type Consumer struct {
orderReader    *kafka.Reader
paymentReader  *kafka.Reader
shipmentReader *kafka.Reader
service        interfaces.NotificationService
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, service interfaces.NotificationService) *Consumer {
orderReader := kafka.NewReader(kafka.ReaderConfig{
Brokers:         []string{cfg.KafkaBootstrapServers},
Topic:           "orders",
GroupID:         "notification-service",
MinBytes:        10e3, // 10KB
MaxBytes:        10e6, // 10MB
MaxWait:         1 * time.Second,
ReadLagInterval: -1,
})

paymentReader := kafka.NewReader(kafka.ReaderConfig{
Brokers:         []string{cfg.KafkaBootstrapServers},
Topic:           "payments",
GroupID:         "notification-service",
MinBytes:        10e3, // 10KB
MaxBytes:        10e6, // 10MB
MaxWait:         1 * time.Second,
ReadLagInterval: -1,
})

shipmentReader := kafka.NewReader(kafka.ReaderConfig{
Brokers:         []string{cfg.KafkaBootstrapServers},
Topic:           "shipments",
GroupID:         "notification-service",
MinBytes:        10e3, // 10KB
MaxBytes:        10e6, // 10MB
MaxWait:         1 * time.Second,
ReadLagInterval: -1,
})

return &Consumer{
orderReader:    orderReader,
paymentReader:  paymentReader,
shipmentReader: shipmentReader,
service:        service,
}
}

// StartConsuming starts consuming messages from Kafka
func (c *Consumer) StartConsuming(ctx context.Context) {
// Start consuming order events
go func() {
for {
select {
case <-ctx.Done():
log.Println("Stopping order Kafka consumer")
c.orderReader.Close()
return
default:
// Read message
msg, err := c.orderReader.ReadMessage(ctx)
if err != nil {
log.Printf("Error reading order message: %v", err)
continue
}

// Process message
err = c.processOrderMessage(msg)
if err != nil {
log.Printf("Error processing order message: %v", err)
continue
}
}
}
}()

// Start consuming payment events
go func() {
for {
select {
case <-ctx.Done():
log.Println("Stopping payment Kafka consumer")
c.paymentReader.Close()
return
default:
// Read message
msg, err := c.paymentReader.ReadMessage(ctx)
if err != nil {
log.Printf("Error reading payment message: %v", err)
continue
}

// Process message
err = c.processPaymentMessage(msg)
if err != nil {
log.Printf("Error processing payment message: %v", err)
continue
}
}
}
}()

// Start consuming shipment events
go func() {
for {
select {
case <-ctx.Done():
log.Println("Stopping shipment Kafka consumer")
c.shipmentReader.Close()
return
default:
// Read message
msg, err := c.shipmentReader.ReadMessage(ctx)
if err != nil {
log.Printf("Error reading shipment message: %v", err)
continue
}

// Process message
err = c.processShipmentMessage(msg)
if err != nil {
log.Printf("Error processing shipment message: %v", err)
continue
}
}
}
}()

log.Println("Kafka consumers started")
}

// processOrderMessage processes an order message
func (c *Consumer) processOrderMessage(msg kafka.Message) error {
// Parse message
var event models.OrderEvent
err := json.Unmarshal(msg.Value, &event)
if err != nil {
return err
}

log.Printf("Received message from Kafka orders topic: %s", string(msg.Value))

// Process event
return c.service.ProcessOrderEvent(event)
}

// processPaymentMessage processes a payment message
func (c *Consumer) processPaymentMessage(msg kafka.Message) error {
// Parse message
var event models.PaymentEvent
err := json.Unmarshal(msg.Value, &event)
if err != nil {
return err
}

log.Printf("Received message from Kafka payments topic: %s", string(msg.Value))

// Process event
return c.service.ProcessPaymentEvent(event)
}

// processShipmentMessage processes a shipment message
func (c *Consumer) processShipmentMessage(msg kafka.Message) error {
// Parse message
var event models.ShipmentEvent
err := json.Unmarshal(msg.Value, &event)
if err != nil {
return err
}

log.Printf("Received message from Kafka shipments topic: %s", string(msg.Value))

// Process event
return c.service.ProcessShipmentEvent(event)
}
