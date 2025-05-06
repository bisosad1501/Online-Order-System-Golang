package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/online-order-system/order-service/config"
	"github.com/online-order-system/order-service/interfaces"
	"github.com/online-order-system/order-service/models"
	"github.com/segmentio/kafka-go"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	orderReader    *kafka.Reader
	paymentReader  *kafka.Reader
	shipmentReader *kafka.Reader
	dlqWriter      *kafka.Writer
	service        interfaces.OrderService
	maxRetries     int
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, service interfaces.OrderService) *Consumer {
	// Reader for orders topic
	orderReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBootstrapServers},
		Topic:    cfg.KafkaTopic,
		GroupID:  "order-service-orders",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	// Reader for payments topic
	paymentReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBootstrapServers},
		Topic:    "payments",
		GroupID:  "order-service-payments",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	// Reader for shipments topic
	shipmentReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBootstrapServers},
		Topic:    "shipments",
		GroupID:  "order-service-shipments",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	// Writer for DLQ
	dlqWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBootstrapServers),
		Topic:    "order-service-dlq",
		Balancer: &kafka.LeastBytes{},
	}

	return &Consumer{
		orderReader:    orderReader,
		paymentReader:  paymentReader,
		shipmentReader: shipmentReader,
		dlqWriter:      dlqWriter,
		service:        service,
		maxRetries:     2, // Retry 2 times as per design
	}
}

// StartConsuming starts consuming messages from Kafka
func (c *Consumer) StartConsuming(ctx context.Context) {
	// Start consuming from orders topic
	go c.consumeOrders(ctx)

	// Start consuming from payments topic
	go c.consumePayments(ctx)

	// Start consuming from shipments topic
	go c.consumeShipments(ctx)
}

// consumeOrders consumes messages from the orders topic
func (c *Consumer) consumeOrders(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer for orders topic")
			c.orderReader.Close()
			return
		default:
			m, err := c.orderReader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from Kafka orders topic: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			log.Printf("Received message from Kafka orders topic: %s", string(m.Value))

			// Process the message
			var event models.OrderEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				// Send to DLQ if unmarshaling fails
				c.sendToDLQ(m.Value, "orders", "unmarshal_error", err.Error())
				continue
			}

			// Process the event based on its type
			switch event.EventType {
			case "order_created":
				log.Printf("Processing order created event for order %s", event.OrderID)
				// Implement processing logic with retry
				var processErr error
				for i := 0; i <= c.maxRetries; i++ {
					if i > 0 {
						// Exponential backoff: 1s, 2s, 4s, ...
						backoffTime := time.Duration(1<<uint(i-1)) * time.Second
						log.Printf("Retrying after %v...", backoffTime)
						time.Sleep(backoffTime)
					}

					// Process order created event
					// This is just a placeholder as the actual implementation depends on your business logic
					// In a real implementation, you would call appropriate service methods

					// Simulate success for now
					processErr = nil
					break
				}

				// If all retries failed, send to DLQ
				if processErr != nil {
					log.Printf("All retries failed for order_created event for order %s: %v", event.OrderID, processErr)
					c.sendToDLQ(m.Value, "orders", "processing_error", processErr.Error())
				}
			}
		}
	}
}

// consumePayments consumes messages from the payments topic
func (c *Consumer) consumePayments(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer for payments topic")
			c.paymentReader.Close()
			return
		default:
			m, err := c.paymentReader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from Kafka payments topic: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			log.Printf("Received message from Kafka payments topic: %s", string(m.Value))

			// Process the message
			var event map[string]interface{}
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				// Send to DLQ if unmarshaling fails
				c.sendToDLQ(m.Value, "payments", "unmarshal_error", err.Error())
				continue
			}

			// Process the event based on its type
			eventType, ok := event["event_type"].(string)
			if !ok {
				log.Printf("Error: event_type is not a string")
				// Send to DLQ if event_type is invalid
				c.sendToDLQ(m.Value, "payments", "invalid_event_type", "event_type is not a string")
				continue
			}

			orderID, ok := event["order_id"].(string)
			if !ok {
				log.Printf("Error: order_id is not a string")
				// Send to DLQ if order_id is invalid
				c.sendToDLQ(m.Value, "payments", "invalid_order_id", "order_id is not a string")
				continue
			}

			if eventType == "payment_successful" {
				log.Printf("Processing payment successful event for order %s", orderID)

				// Implement retry with exponential backoff
				var processErr error
				for i := 0; i <= c.maxRetries; i++ {
					if i > 0 {
						// Exponential backoff: 1s, 2s, 4s, ...
						backoffTime := time.Duration(1<<uint(i-1)) * time.Second
						log.Printf("Retrying after %v...", backoffTime)
						time.Sleep(backoffTime)
					}

					// Update order status to CONFIRMED to trigger shipping
					err = c.service.UpdateOrderStatus(orderID, models.OrderStatusConfirmed)
					if err != nil {
						processErr = err
						log.Printf("Error updating order status to CONFIRMED (attempt %d/%d): %v", i+1, c.maxRetries+1, err)
						continue
					}

					// Success, break the retry loop
					processErr = nil
					break
				}

				// If all retries failed, send to DLQ
				if processErr != nil {
					log.Printf("All retries failed for payment_successful event for order %s: %v", orderID, processErr)
					c.sendToDLQ(m.Value, "payments", "processing_error", processErr.Error())
				}
			} else if eventType == "payment_failed" {
				log.Printf("Processing payment failed event for order %s", orderID)

				// Implement retry with exponential backoff
				var processErr error
				for i := 0; i <= c.maxRetries; i++ {
					if i > 0 {
						// Exponential backoff: 1s, 2s, 4s, ...
						backoffTime := time.Duration(1<<uint(i-1)) * time.Second
						log.Printf("Retrying after %v...", backoffTime)
						time.Sleep(backoffTime)
					}

					// Get the order
					order, err := c.service.GetOrderByID(orderID)
					if err != nil {
						processErr = err
						log.Printf("Error getting order %s (attempt %d/%d): %v", orderID, i+1, c.maxRetries+1, err)
						continue
					}

					// Skip if order is already in FAILED state
					if order.Status == models.OrderStatusFailed {
						log.Printf("Order %s is already in FAILED state, skipping", orderID)
						processErr = nil
						break
					}

					// Skip if order is in DELIVERED state (too late to compensate)
					if order.Status == models.OrderStatusDelivered {
						log.Printf("Order %s is already in DELIVERED state, too late to compensate", orderID)
						processErr = nil
						break
					}

					// Skip if order is in CONFIRMED state (payment already successful)
					if order.Status == models.OrderStatusConfirmed {
						log.Printf("Order %s is already in CONFIRMED state, payment was successful, skipping", orderID)
						processErr = nil
						break
					}

					// Call Compensate directly
					err = c.service.Compensate(order, "payment_failed")
					if err != nil {
						processErr = err
						log.Printf("Error compensating for order %s (attempt %d/%d): %v", orderID, i+1, c.maxRetries+1, err)
						continue
					}

					log.Printf("Successfully compensated for order %s", orderID)
					processErr = nil
					break
				}

				// If all retries failed, send to DLQ
				if processErr != nil {
					log.Printf("All retries failed for payment_failed event for order %s: %v", orderID, processErr)
					c.sendToDLQ(m.Value, "payments", "processing_error", processErr.Error())
				}
			}
		}
	}
}

// consumeShipments consumes messages from the shipments topic
func (c *Consumer) consumeShipments(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer for shipments topic")
			c.shipmentReader.Close()
			return
		default:
			m, err := c.shipmentReader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from Kafka shipments topic: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			log.Printf("Received message from Kafka shipments topic: %s", string(m.Value))

			// Process the message
			var event struct {
				EventType string `json:"event_type"`
				OrderID   string `json:"order_id"`
			}
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				// Send to DLQ if unmarshaling fails
				c.sendToDLQ(m.Value, "shipments", "unmarshal_error", err.Error())
				continue
			}

			// Process the event based on its type
			switch event.EventType {
			case "shipping_completed":
				log.Printf("Processing shipping completed event for order %s", event.OrderID)
				// Implement retry with exponential backoff
				var processErr error
				for i := 0; i <= c.maxRetries; i++ {
					if i > 0 {
						// Exponential backoff: 1s, 2s, 4s, ...
						backoffTime := time.Duration(1<<uint(i-1)) * time.Second
						log.Printf("Retrying after %v...", backoffTime)
						time.Sleep(backoffTime)
					}

					err := c.service.UpdateOrderStatus(event.OrderID, models.OrderStatusDelivered)
					if err == nil {
						// Success, break the retry loop
						processErr = nil
						break
					}

					processErr = err
					log.Printf("Error updating order status (attempt %d/%d): %v", i+1, c.maxRetries+1, err)
				}

				// If all retries failed, send to DLQ
				if processErr != nil {
					log.Printf("All retries failed for shipping_completed event for order %s: %v", event.OrderID, processErr)
					c.sendToDLQ(m.Value, "shipments", "processing_error", processErr.Error())
				}
			}
		}
	}
}

// sendToDLQ sends a failed message to the Dead Letter Queue
func (c *Consumer) sendToDLQ(value []byte, sourceTopic, errorType, errorDetails string) {
	// Create DLQ message with metadata
	dlqMessage := map[string]interface{}{
		"original_message": string(value),
		"source_topic":     sourceTopic,
		"error_type":       errorType,
		"error_details":    errorDetails,
		"timestamp":        time.Now().Unix(),
	}

	// Marshal DLQ message
	dlqValue, err := json.Marshal(dlqMessage)
	if err != nil {
		log.Printf("Error marshaling DLQ message: %v", err)
		return
	}

	// Send to DLQ
	err = c.dlqWriter.WriteMessages(context.Background(),
		kafka.Message{
			Value: dlqValue,
		},
	)

	if err != nil {
		log.Printf("Error sending message to DLQ: %v", err)
	} else {
		log.Printf("Message sent to DLQ successfully")
	}
}
