package interfaces

import (
"github.com/online-order-system/order-service/models"
)

// OrderService defines the interface for order service
type OrderService interface {
CreateOrder(req models.CreateOrderRequest) (models.Order, error)
GetOrderByID(id string) (models.Order, error)
GetOrders() ([]models.Order, error)
UpdateOrderStatus(id string, status models.OrderStatus) error
Compensate(order models.Order, failureReason string) error
RetryPayment(orderID string, req models.RetryPaymentRequest) (models.Order, error)
}

// OrderProducer defines the interface for order producer
type OrderProducer interface {
PublishOrderCreated(order models.Order) error
PublishOrderConfirmed(order models.Order) error
PublishOrderCancelled(order models.Order) error
PublishOrderCompleted(order models.Order, trackingNumber string) error
PublishOrderEvent(event models.OrderEvent) error
Close() error
}
