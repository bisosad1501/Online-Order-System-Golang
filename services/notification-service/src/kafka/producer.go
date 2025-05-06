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

// Producer represents a Kafka producer
type Producer struct {
writer *kafka.Writer
}

// Ensure Producer implements NotificationProducer interface
var _ interfaces.NotificationProducer = (*Producer)(nil)

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

// PublishNotificationCreated publishes a notification created event
func (p *Producer) PublishNotificationCreated(notification models.Notification) error {
event := models.NotificationEvent{
EventType:  "notification_created",
CustomerID: notification.CustomerID,
Type:       notification.Type,
Status:     notification.Status,
Subject:    notification.Subject,
Content:    notification.Content,
Recipient:  notification.Recipient,
Timestamp:  time.Now().Unix(),
}

return p.publishEvent(event)
}

// PublishNotificationSent publishes a notification sent event
func (p *Producer) PublishNotificationSent(notification models.Notification) error {
event := models.NotificationEvent{
EventType:  "notification_sent",
CustomerID: notification.CustomerID,
Type:       notification.Type,
Status:     notification.Status,
Subject:    notification.Subject,
Content:    notification.Content,
Recipient:  notification.Recipient,
Timestamp:  time.Now().Unix(),
}

return p.publishEvent(event)
}

// PublishNotificationFailed publishes a notification failed event
func (p *Producer) PublishNotificationFailed(notification models.Notification) error {
event := models.NotificationEvent{
EventType:  "notification_failed",
CustomerID: notification.CustomerID,
Type:       notification.Type,
Status:     notification.Status,
Subject:    notification.Subject,
Content:    notification.Content,
Recipient:  notification.Recipient,
Timestamp:  time.Now().Unix(),
}

return p.publishEvent(event)
}

// publishEvent publishes an event to Kafka
func (p *Producer) publishEvent(event models.NotificationEvent) error {
// Marshal event to JSON
eventJSON, err := json.Marshal(event)
if err != nil {
return err
}

// Create message
msg := kafka.Message{
Key:   []byte(event.CustomerID),
Value: eventJSON,
Time:  time.Now(),
}

// Write message to Kafka
err = p.writer.WriteMessages(context.Background(), msg)
if err != nil {
return err
}

log.Printf("Published event: %s for user %s", event.EventType, event.CustomerID)
return nil
}
