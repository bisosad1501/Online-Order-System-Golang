package kafka

import (
"context"
"encoding/json"
"log"
"time"

"github.com/segmentio/kafka-go"
"github.com/online-order-system/payment-service/config"
"github.com/online-order-system/payment-service/interfaces"
"github.com/online-order-system/payment-service/models"
)

// Consumer represents a Kafka consumer
type Consumer struct {
reader  *kafka.Reader
service interfaces.PaymentService
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, service interfaces.PaymentService) *Consumer {
reader := kafka.NewReader(kafka.ReaderConfig{
Brokers:  []string{cfg.KafkaBootstrapServers},
Topic:    "orders", // Listen to orders topic
GroupID:  "payment-service",
MinBytes: 10e3, // 10KB
MaxBytes: 10e6, // 10MB
})

log.Printf("Payment Service Kafka consumer created and subscribed to topic: orders")

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
case "order_created":
	log.Printf("Processing order created event")
	var orderEvent struct {
		OrderID     string  `json:"order_id"`
		CustomerID  string  `json:"customer_id"`
		TotalAmount float64 `json:"total_amount"`
	}
	if err := json.Unmarshal(m.Value, &orderEvent); err != nil {
		log.Printf("Error unmarshaling order created event: %v", err)
		continue
	}

	// Check if payment already exists for this order
	_, err := c.service.GetPaymentByOrderID(orderEvent.OrderID)
	if err == nil {
		// Payment already exists, skip
		log.Printf("Payment already exists for order %s, skipping", orderEvent.OrderID)
		continue
	}

	// Create payment request
	paymentReq := struct {
		OrderID       string  `json:"order_id"`
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"payment_method"`
		Currency      string  `json:"currency"`
		Description   string  `json:"description"`
	}{
		OrderID:       orderEvent.OrderID,
		Amount:        orderEvent.TotalAmount,
		PaymentMethod: "STRIPE", // Default to STRIPE
		Currency:      "USD",    // Default to USD
		Description:   "Payment for order " + orderEvent.OrderID,
	}

	// Convert to models.CreatePaymentRequest
	createPaymentReq := models.CreatePaymentRequest{
		OrderID:       paymentReq.OrderID,
		Amount:        paymentReq.Amount,
		PaymentMethod: paymentReq.PaymentMethod,
		Currency:      paymentReq.Currency,
		Description:   paymentReq.Description,
	}

	// Create payment
	_, err = c.service.CreatePayment(createPaymentReq)
	if err != nil {
		log.Printf("Error creating payment for order %s: %v", orderEvent.OrderID, err)
	}

case "order_cancelled":
	log.Printf("Processing order cancelled event")
	var orderEvent struct {
		OrderID string `json:"order_id"`
	}
	if err := json.Unmarshal(m.Value, &orderEvent); err != nil {
		log.Printf("Error unmarshaling order cancelled event: %v", err)
		continue
	}

	// Get payment for order
	payment, err := c.service.GetPaymentByOrderID(orderEvent.OrderID)
	if err != nil {
		log.Printf("Error getting payment for order %s: %v", orderEvent.OrderID, err)
		continue
	}

	// Refund payment if it was successful
	if payment.Status == "SUCCESSFUL" {
		_, err = c.service.UpdatePaymentStatus(payment.ID, "REFUNDED")
		if err != nil {
			log.Printf("Error refunding payment %s: %v", payment.ID, err)
		}
	}
}
}
}
}()
}
