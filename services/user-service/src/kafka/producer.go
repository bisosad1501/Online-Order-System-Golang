package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/online-order-system/user-service/config"
	"github.com/online-order-system/user-service/interfaces"
	"github.com/online-order-system/user-service/models"
	"github.com/segmentio/kafka-go"
)

// Producer handles Kafka message production
type Producer struct {
	writer *kafka.Writer
	topic  string
}

// Ensure Producer implements UserProducer interface
var _ interfaces.UserProducer = (*Producer)(nil)

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.Config) *Producer {
	// Create Kafka writer
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBootstrapServers),
		Topic:    cfg.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	log.Println("Kafka producer created")
	return &Producer{writer: w, topic: cfg.KafkaTopic}
}

// PublishUserVerified publishes a user verified event
func (p *Producer) PublishUserVerified(user models.Customer) error {
	if p.writer == nil {
		log.Println("Kafka producer not available, skipping event publishing")
		return nil
	}

	event := models.CustomerEvent{
		EventType:  "user_verified",
		CustomerID: user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Timestamp:  time.Now().Unix(),
	}

	return p.publishEvent(event)
}

// PublishUserUpdated publishes a user updated event
func (p *Producer) PublishUserUpdated(user models.Customer) error {
	if p.writer == nil {
		log.Println("Kafka producer not available, skipping event publishing")
		return nil
	}

	event := models.CustomerEvent{
		EventType:  "user_updated",
		CustomerID: user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Timestamp:  time.Now().Unix(),
	}

	return p.publishEvent(event)
}

// publishEvent publishes an event to Kafka
func (p *Producer) publishEvent(event interface{}) error {
	// Convert event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Create Kafka message
	message := kafka.Message{
		Value: eventJSON,
	}

	// Produce message
	return p.writer.WriteMessages(context.Background(), message)
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
