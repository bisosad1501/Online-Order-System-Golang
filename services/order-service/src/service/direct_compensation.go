package service

import (
"fmt"
"log"

"github.com/online-order-system/order-service/models"
)

// DirectCompensate is a direct compensation function that can be called from outside the service
func DirectCompensate(orderID string, failureReason string) error {
log.Printf("Direct compensation for order %s with reason: %s", orderID, failureReason)

// Get the order service instance
orderService := GetOrderServiceInstance()
if orderService == nil {
return fmt.Errorf("order service instance is nil")
}

// Get the order
order, err := orderService.GetOrderByID(orderID)
if err != nil {
return fmt.Errorf("failed to get order: %v", err)
}

// Skip if order is already in FAILED state
if order.Status == models.OrderStatusFailed {
log.Printf("Order %s is already in FAILED state, skipping", orderID)
return nil
}

// Skip if order is in DELIVERED state (too late to compensate)
if order.Status == models.OrderStatusDelivered {
log.Printf("Order %s is already in DELIVERED state, too late to compensate", orderID)
return nil
}

// Call compensate method
return orderService.Compensate(order, failureReason)
}

// orderServiceInstance is a global variable to store the order service instance
var orderServiceInstance *OrderService

// SetOrderServiceInstance sets the order service instance
func SetOrderServiceInstance(service *OrderService) {
orderServiceInstance = service
}

// GetOrderServiceInstance gets the order service instance
func GetOrderServiceInstance() *OrderService {
return orderServiceInstance
}
