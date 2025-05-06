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

// PaymentConsumer represents a Kafka consumer for payment events
type PaymentConsumer struct {
	reader  *kafka.Reader
	service interfaces.OrderService
}

// NewPaymentConsumer creates a new Kafka consumer for payment events
func NewPaymentConsumer(cfg *config.Config, service interfaces.OrderService) *PaymentConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBootstrapServers},
		Topic:    "payments", // Listen to payments topic
		GroupID:  "order-service",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	log.Printf("Order Service Payment Kafka consumer created and subscribed to topic: payments")

	return &PaymentConsumer{
		reader:  reader,
		service: service,
	}
}

// StartConsuming starts consuming messages from Kafka
func (c *PaymentConsumer) StartConsuming(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping Payment Kafka consumer")
				c.reader.Close()
				return
			default:
				m, err := c.reader.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message from Kafka: %v", err)
					time.Sleep(1 * time.Second)
					continue
				}

				log.Printf("Received payment message from Kafka: %s", string(m.Value))

				// Process the message
				var event struct {
					EventType string `json:"event_type"`
				}
				if err := json.Unmarshal(m.Value, &event); err != nil {
					log.Printf("Error unmarshaling payment message: %v", err)
					continue
				}

				// Process the event based on its type
				switch event.EventType {
				case "payment_successful":
					log.Printf("Processing payment successful event")
					var paymentEvent struct {
						OrderID string `json:"order_id"`
					}
					if err := json.Unmarshal(m.Value, &paymentEvent); err != nil {
						log.Printf("Error unmarshaling payment successful event: %v", err)
						continue
					}

					// Update order status to CONFIRMED
					err = c.service.UpdateOrderStatus(paymentEvent.OrderID, models.OrderStatusConfirmed)
					if err != nil {
						log.Printf("Error updating order status for order %s: %v", paymentEvent.OrderID, err)
					}

				case "payment_failed":
					log.Printf("Processing payment failed event")
					var paymentEvent struct {
						OrderID string `json:"order_id"`
					}
					if err := json.Unmarshal(m.Value, &paymentEvent); err != nil {
						log.Printf("Error unmarshaling payment failed event: %v", err)
						continue
					}

					// Update order status to FAILED
					err = c.service.UpdateOrderStatus(paymentEvent.OrderID, models.OrderStatusFailed)
					if err != nil {
						log.Printf("Error updating order status for order %s: %v", paymentEvent.OrderID, err)
					}
				}
			}
		}
	}()
}
