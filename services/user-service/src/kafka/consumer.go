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

// Consumer handles Kafka message consumption
type Consumer struct {
	reader  *kafka.Reader
	service interfaces.UserService
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, service interfaces.UserService) *Consumer {
	// Create Kafka reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBootstrapServers},
		GroupID:  "user-service",
		Topic:    cfg.KafkaTopic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	log.Println("Kafka consumer created and subscribed to topics")
	return &Consumer{reader: r, service: service}
}

// StartConsuming starts consuming messages from Kafka
func (c *Consumer) StartConsuming(ctx context.Context) {
	if c.reader == nil {
		log.Println("Kafka consumer not available, skipping message consumption")
		return
	}

	// Start consuming in a goroutine
	go func() {
		log.Println("Starting to consume messages from Kafka")
		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, stopping Kafka consumer")
				c.reader.Close()
				return
			default:
				// Read message
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message: %v", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}

				// Process message
				c.processMessage(msg)
			}
		}
	}()
}

// processMessage processes a Kafka message
func (c *Consumer) processMessage(msg kafka.Message) {
	// Parse message value
	var event models.OrderEvent
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	// Process event based on type
	switch event.EventType {
	case "order_created":
		err = c.service.AddUserOrder(event.CustomerID, event.OrderID, event.Status)
		if err != nil {
			log.Printf("Failed to add user order: %v", err)
		}
	case "order_confirmed", "order_paid", "order_shipped", "order_delivered", "order_cancelled":
		err = c.service.UpdateUserOrderStatus(event.CustomerID, event.OrderID, event.Status)
		if err != nil {
			log.Printf("Failed to update user order status: %v", err)
		}
	default:
		log.Printf("Unknown event type: %s", event.EventType)
	}
}
