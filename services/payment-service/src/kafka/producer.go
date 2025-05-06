package kafka

import (
"context"
"encoding/json"
"log"

"github.com/segmentio/kafka-go"
"github.com/online-order-system/payment-service/config"
"github.com/online-order-system/payment-service/interfaces"
"github.com/online-order-system/payment-service/models"
)

// Producer represents a Kafka producer
type Producer struct {
writer *kafka.Writer
}

// Ensure Producer implements PaymentProducer interface
var _ interfaces.PaymentProducer = (*Producer)(nil)

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

// PublishPaymentCreated publishes a payment created event
func (p *Producer) PublishPaymentCreated(payment models.Payment) error {
event := models.PaymentEvent{
EventType:     "payment_created",
PaymentID:     payment.ID,
OrderID:       payment.OrderID,
Amount:        payment.Amount,
Status:        payment.Status,
PaymentMethod: payment.PaymentMethod,
Timestamp:     payment.CreatedAt.Unix(),
}

return p.publishEvent(event)
}

// PublishPaymentSuccessful publishes a payment successful event
func (p *Producer) PublishPaymentSuccessful(payment models.Payment) error {
event := models.PaymentEvent{
EventType:     "payment_successful",
PaymentID:     payment.ID,
OrderID:       payment.OrderID,
Amount:        payment.Amount,
Status:        payment.Status,
PaymentMethod: payment.PaymentMethod,
Timestamp:     payment.UpdatedAt.Unix(),
}

return p.publishEvent(event)
}

// PublishPaymentFailed publishes a payment failed event
func (p *Producer) PublishPaymentFailed(payment models.Payment) error {
event := models.PaymentEvent{
EventType:     "payment_failed",
PaymentID:     payment.ID,
OrderID:       payment.OrderID,
Amount:        payment.Amount,
Status:        payment.Status,
PaymentMethod: payment.PaymentMethod,
Timestamp:     payment.UpdatedAt.Unix(),
}

return p.publishEvent(event)
}

// PublishPaymentRefunded publishes a payment refunded event
func (p *Producer) PublishPaymentRefunded(payment models.Payment) error {
event := models.PaymentEvent{
EventType:     "payment_refunded",
PaymentID:     payment.ID,
OrderID:       payment.OrderID,
Amount:        payment.Amount,
Status:        payment.Status,
PaymentMethod: payment.PaymentMethod,
Timestamp:     payment.UpdatedAt.Unix(),
}

return p.publishEvent(event)
}

// publishEvent publishes an event to Kafka
func (p *Producer) publishEvent(event models.PaymentEvent) error {
// Convert event to JSON
eventJSON, err := json.Marshal(event)
if err != nil {
return err
}

// Create message
message := kafka.Message{
Key:   []byte(event.PaymentID),
Value: eventJSON,
}

// Write message
err = p.writer.WriteMessages(context.Background(), message)
if err != nil {
return err
}

log.Printf("Published event: %s for payment %s", event.EventType, event.PaymentID)
return nil
}
