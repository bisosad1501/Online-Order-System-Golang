package kafka

import (
"context"
"encoding/json"
"log"
"time"

"github.com/segmentio/kafka-go"
"github.com/online-order-system/order-service/config"
"github.com/online-order-system/order-service/interfaces"
"github.com/online-order-system/order-service/models"
)

// Producer represents a Kafka producer
type Producer struct {
writer *kafka.Writer
}

// Ensure Producer implements OrderProducer interface
var _ interfaces.OrderProducer = (*Producer)(nil)

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

// PublishOrderCreated publishes an order created event
func (p *Producer) PublishOrderCreated(order models.Order) error {
	log.Printf("Publishing order_created event for order %s", order.ID)
	event := models.OrderEvent{
		EventType:       "order_created",
		OrderID:         order.ID,
		CustomerID:      order.CustomerID,
		Status:          order.Status,
		TotalAmount:     order.TotalAmount,
		Timestamp:       order.CreatedAt.Unix(),
		Items:           order.Items,
		ShippingAddress: order.ShippingAddress,
	}

	return p.publishEvent(event)
}

// PublishOrderConfirmed publishes an order confirmed event
func (p *Producer) PublishOrderConfirmed(order models.Order) error {
	log.Printf("Publishing order_confirmed event for order %s", order.ID)
	event := models.OrderEvent{
		EventType:   "order_confirmed",
		OrderID:     order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Timestamp:   order.UpdatedAt.Unix(),
		ShippingAddress: order.ShippingAddress, // Add shipping address for notification service
	}

	return p.publishEvent(event)
}

// PublishOrderCancelled publishes an order cancelled event
func (p *Producer) PublishOrderCancelled(order models.Order) error {
event := models.OrderEvent{
EventType:   "order_cancelled",
OrderID:     order.ID,
CustomerID:  order.CustomerID,
Status:      order.Status,
TotalAmount: order.TotalAmount,
Timestamp:   order.UpdatedAt.Unix(),
}

return p.publishEvent(event)
}

// PublishOrderCompleted publishes an order completed event
func (p *Producer) PublishOrderCompleted(order models.Order, trackingNumber string) error {
event := models.OrderEvent{
EventType:      "order_completed",
OrderID:        order.ID,
CustomerID:     order.CustomerID,
Status:         order.Status,
TotalAmount:    order.TotalAmount,
Timestamp:      order.UpdatedAt.Unix(),
}

return p.publishEvent(event)
}

// PublishOrderEvent publishes an order event
func (p *Producer) PublishOrderEvent(event models.OrderEvent) error {
return p.publishEvent(event)
}

// publishEvent publishes an event to Kafka
func (p *Producer) publishEvent(event models.OrderEvent) error {
// Marshal event to JSON
eventJSON, err := json.Marshal(event)
if err != nil {
return err
}

// Create message
msg := kafka.Message{
Key:   []byte(event.OrderID),
Value: eventJSON,
Time:  time.Now(),
}

// Write message to Kafka
err = p.writer.WriteMessages(context.Background(), msg)
if err != nil {
return err
}

log.Printf("Published event: %s for order %s", event.EventType, event.OrderID)
return nil
}
