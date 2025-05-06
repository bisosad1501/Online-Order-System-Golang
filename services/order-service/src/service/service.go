package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/online-order-system/order-service/config"
	"github.com/online-order-system/order-service/db"
	"github.com/online-order-system/order-service/interfaces"
	"github.com/online-order-system/order-service/models"
	"github.com/online-order-system/order-service/utils"
)

// OrderService handles business logic for orders
type OrderService struct {
	config     *config.Config
	repository *db.OrderRepository
	producer   interfaces.OrderProducer
	httpClient *utils.HTTPClient
}

// Ensure OrderService implements OrderService interface
var _ interfaces.OrderService = (*OrderService)(nil)

// NewOrderService creates a new order service
func NewOrderService(cfg *config.Config, repo *db.OrderRepository, producer interfaces.OrderProducer) *OrderService {
	return &OrderService{
		config:     cfg,
		repository: repo,
		producer:   producer,
		httpClient: utils.NewHTTPClient(),
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(req models.CreateOrderRequest) (models.Order, error) {
	// Create order
	now := time.Now()
	order := models.Order{
		ID:              uuid.New().String(),
		CustomerID:      req.CustomerID,
		Status:          models.OrderStatusCreated,
		TotalAmount:     calculateTotalAmount(req.Items),
		ShippingAddress: req.ShippingAddress,
		Items:           req.Items,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Verify customer (Step 2 in design)
	err := s.verifyCustomer(req.CustomerID)
	if err != nil {
		log.Printf("Failed to verify customer: %v", err)
		return models.Order{}, fmt.Errorf("failed to verify customer: %v", err)
	}

	// Create audit log for order creation
	auditLog := models.AuditLog{
		ID:          uuid.New().String(),
		ServiceName: "order-service",
		Action:      models.AuditLogActionCreateOrder,
		CustomerID:  req.CustomerID,
		Timestamp:   now,
		Details:     fmt.Sprintf("Order created with ID: %s, Total Amount: %.2f", order.ID, order.TotalAmount),
	}

	err = s.repository.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("Failed to create audit log: %v", err)
		// Continue anyway
	}

	// Get cart items (Step 1-2 in design)
	if len(req.Items) == 0 {
		cartItems, err := s.getCartItems(req.CustomerID)
		if err != nil {
			log.Printf("Failed to get cart items: %v", err)
			return models.Order{}, fmt.Errorf("failed to get cart items: %v", err)
		}

		if len(cartItems) == 0 {
			return models.Order{}, errors.New("cart is empty")
		}

		order.Items = cartItems
		order.TotalAmount = calculateTotalAmount(cartItems)
	}

	// Set order ID for each item
	for i := range order.Items {
		order.Items[i].ID = uuid.New().String()
		order.Items[i].OrderID = order.ID
	}

	// Save order to database
	err = s.repository.CreateOrder(order)
	if err != nil {
		return models.Order{}, err
	}

	// Check inventory
	available, err := s.checkInventory(req.Items)
	if err != nil {
		log.Printf("Failed to check inventory: %v", err)

		// Create audit log for inventory check failure
		auditLog := models.AuditLog{
			ID:          uuid.New().String(),
			ServiceName: "order-service",
			Action:      models.AuditLogActionInventoryError,
			CustomerID:  order.CustomerID,
			Timestamp:   time.Now(),
			Details:     fmt.Sprintf("Failed to check inventory for order ID: %s, Error: %v", order.ID, err),
		}

		s.repository.CreateAuditLog(auditLog)

		return order, fmt.Errorf("failed to check inventory: %v", err)
	}
	if !available {
		// Create audit log for inventory unavailable
		auditLog := models.AuditLog{
			ID:          uuid.New().String(),
			ServiceName: "order-service",
			Action:      models.AuditLogActionInventoryError,
			CustomerID:  order.CustomerID,
			Timestamp:   time.Now(),
			Details:     fmt.Sprintf("Inventory unavailable for order ID: %s", order.ID),
		}

		s.repository.CreateAuditLog(auditLog)

		// Compensate with reason "inventory_unavailable"
		s.Compensate(order, "inventory_unavailable")
		return order, errors.New("some items are not available in inventory")
	}

	// Mark inventory as locked
	order.InventoryLocked = true
	err = s.repository.UpdateOrder(order)
	if err != nil {
		log.Printf("Failed to update order: %v", err)
	}

	// Process payment
	log.Printf("Processing payment for order %s", order.ID)

	// Create audit log for payment processing
	auditLog2 := models.AuditLog{
		ID:          uuid.New().String(),
		ServiceName: "order-service",
		Action:      models.AuditLogActionProcessPayment,
		CustomerID:  order.CustomerID,
		Timestamp:   time.Now(),
		Details:     fmt.Sprintf("Payment processing initiated for order ID: %s, Amount: %.2f", order.ID, order.TotalAmount),
	}

	s.repository.CreateAuditLog(auditLog2)

	// Payment will be processed by payment-service when it receives the order_created event
	// We'll mark the order as CREATED until payment is confirmed
	order.Status = models.OrderStatusCreated
	err = s.repository.UpdateOrder(order)
	if err != nil {
		log.Printf("Failed to update order: %v", err)
	}

	// Create audit log for successful payment
	auditLog3 := models.AuditLog{
		ID:          uuid.New().String(),
		ServiceName: "order-service",
		Action:      models.AuditLogActionProcessPayment,
		CustomerID:  order.CustomerID,
		Timestamp:   time.Now(),
		Details:     fmt.Sprintf("Payment processed successfully for order ID: %s, Amount: %.2f", order.ID, order.TotalAmount),
	}

	s.repository.CreateAuditLog(auditLog3)

	// Clear cart after successful order (Step 4, 7 in design)
	err = s.clearCart(order.CustomerID)
	if err != nil {
		log.Printf("Failed to clear cart: %v", err)
		// Continue anyway for demo purposes
	}

	// Publish order created event
	err = s.producer.PublishOrderCreated(order)
	if err != nil {
		log.Printf("Failed to publish order created event: %v", err)
		// Continue anyway for demo purposes
	}

	// Publish order confirmed event if order is confirmed
	if order.Status == models.OrderStatusConfirmed {
		log.Printf("Order status is CONFIRMED, publishing order_confirmed event for order %s", order.ID)
		err = s.producer.PublishOrderConfirmed(order)
		if err != nil {
			log.Printf("Failed to publish order confirmed event: %v", err)
			// Continue anyway for demo purposes
		}
	} else {
		log.Printf("Order status is not CONFIRMED, not publishing order_confirmed event for order %s", order.ID)
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (s *OrderService) GetOrderByID(id string) (models.Order, error) {
	return s.repository.GetOrderByID(id)
}

// GetOrders retrieves all orders
func (s *OrderService) GetOrders() ([]models.Order, error) {
	return s.repository.GetOrders()
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(id string, status models.OrderStatus) error {
	// Get order
	order, err := s.repository.GetOrderByID(id)
	if err != nil {
		return err
	}

	// Update status
	order.Status = status
	order.UpdatedAt = time.Now()

	// Save to database
	err = s.repository.UpdateOrderStatus(id, status)
	if err != nil {
		return err
	}

	// If order is confirmed, schedule shipping
	if status == models.OrderStatusConfirmed {
		err = s.scheduleShipping(order)
		if err != nil {
			log.Printf("Failed to schedule shipping: %v", err)

			// Create audit log for shipping failure
			auditLog := models.AuditLog{
				ID:          uuid.New().String(),
				ServiceName: "order-service",
				Action:      models.AuditLogActionShipmentUpdate,
				CustomerID:  order.CustomerID,
				Timestamp:   time.Now(),
				Details:     fmt.Sprintf("Failed to schedule shipping for order ID: %s, Error: %v", order.ID, err),
			}

			s.repository.CreateAuditLog(auditLog)

			// Continue anyway for demo purposes
		} else {
			// Create audit log for successful shipping schedule
			auditLog := models.AuditLog{
				ID:          uuid.New().String(),
				ServiceName: "order-service",
				Action:      models.AuditLogActionShipmentUpdate,
				CustomerID:  order.CustomerID,
				Timestamp:   time.Now(),
				Details:     fmt.Sprintf("Shipping scheduled for order ID: %s", order.ID),
			}

			s.repository.CreateAuditLog(auditLog)
		}
	}

	// Publish event based on status
	var publishErr error
	switch status {
	case models.OrderStatusConfirmed:
		publishErr = s.producer.PublishOrderConfirmed(order)
	case models.OrderStatusCancelled:
		publishErr = s.producer.PublishOrderCancelled(order)
	case models.OrderStatusDelivered:
		publishErr = s.producer.PublishOrderCompleted(order, "")
	}

	if publishErr != nil {
		log.Printf("Failed to publish order event: %v", publishErr)
		// Continue anyway for demo purposes
	}

	return nil
}

// checkInventory checks if all items are available in inventory
func (s *OrderService) checkInventory(items []models.OrderItem) (bool, error) {
	// Prepare request
	var checkItems []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	for _, item := range items {
		checkItems = append(checkItems, struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		}{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	checkRequest := models.InventoryCheckRequest{
		Items: checkItems,
	}

	// Send request to inventory service with timeout and retry (2 retries as per design)
	var checkResponse models.InventoryCheckResponse
	inventoryClient := utils.NewHTTPClientWithOptions(2, 100*time.Millisecond, 5*time.Second)
	err := inventoryClient.Post(
		fmt.Sprintf("%s/inventory/check", s.config.InventoryServiceURL),
		checkRequest,
		&checkResponse,
	)
	if err != nil {
		log.Printf("Error checking inventory: %v", err)
		// Check if error is a client error (4xx)
		if err.Error() == "client error: 404" {
			return false, fmt.Errorf("one or more products not found")
		}
		return false, err
	}

	// If items are not available, get recommendations and send notification
	if !checkResponse.Available && len(checkResponse.UnavailableItems) > 0 {
		// Get customer information
		// For demo purposes, we'll use a placeholder email
		customerEmail := "customer@example.com"

		// For each unavailable item, get recommendations and send notification
		for _, item := range checkResponse.UnavailableItems {
			// Get recommendations
			recommendations, err := s.getRecommendations(item.ProductID)
			if err != nil {
				log.Printf("Failed to get recommendations for product %s: %v", item.ProductID, err)
				continue
			}

			// Prepare notification content
			content := fmt.Sprintf("The product '%s' you wanted is out of stock (requested: %d, available: %d). Here are some alternatives:\n\n",
				item.ProductName, item.Requested, item.Available)

			// Add recommendations to content
			for i, product := range recommendations.Products {
				if i >= 3 { // Limit to 3 recommendations
					break
				}
				content += fmt.Sprintf("- %s: $%.2f\n", product.Name, product.Price)
			}

			// Send notification
			err = s.sendNotification(item.ProductID, customerEmail, content)
			if err != nil {
				log.Printf("Failed to send notification for product %s: %v", item.ProductID, err)
			}
		}
	}

	return checkResponse.Available, nil
}

// getRecommendations gets product recommendations from the recommendation service
func (s *OrderService) getRecommendations(productID string) (models.RecommendationResponse, error) {
	// Send request to recommendation service with timeout and retry
	var recommendationResponse models.RecommendationResponse
	err := s.httpClient.Get(
		fmt.Sprintf("%s/recommendations/products/%s", s.config.RecommendationServiceURL, productID),
		&recommendationResponse,
	)
	if err != nil {
		log.Printf("Error getting recommendations: %v", err)
		return models.RecommendationResponse{}, err
	}

	return recommendationResponse, nil
}

// sendNotification sends a notification to the customer
func (s *OrderService) sendNotification(productID string, customerEmail string, content string) error {
	// Prepare request
	notificationRequest := models.CreateNotificationRequest{
		CustomerID: productID, // Using productID as CustomerID for demo purposes
		Type:       "LOW_INVENTORY",
		Channel:    "EMAIL",
		Subject:    "Product Out of Stock - Alternatives Available",
		Content:    content,
		Recipient:  customerEmail,
	}

	// Send request to notification service with timeout and retry
	err := s.httpClient.Post(
		fmt.Sprintf("%s/notifications", s.config.NotificationServiceURL),
		notificationRequest,
		nil,
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
		return err
	}

	return nil
}

// processPayment processes payment for an order
func (s *OrderService) processPayment(order models.Order) error {
	// Get customer information (in a real app, this would come from the user service)
	// For now, we'll use placeholder data
	customerEmail := "customer@example.com"
	customerName := "Test Customer"

	// Prepare request with more detailed information for Stripe
	paymentRequest := models.CreatePaymentRequest{
		OrderID:       order.ID,
		Amount:        order.TotalAmount,
		PaymentMethod: "card", // Use "card" as the standard payment method
		Currency:      "vnd",  // Vietnamese Dong
		Description:   fmt.Sprintf("Payment for order %s", order.ID),
		CustomerEmail: customerEmail,
		CustomerName:  customerName,
	}

	log.Printf("Sending payment request to payment service: %+v", paymentRequest)

	// Send request to payment service with timeout and retry (2 retries as per design)
	var paymentResponse struct {
		ID                string `json:"id"`
		OrderID           string `json:"order_id"`
		Status            string `json:"status"`
		StripePaymentID   string `json:"stripe_payment_id,omitempty"`
		StripeClientSecret string `json:"stripe_client_secret,omitempty"`
	}

	paymentClient := utils.NewHTTPClientWithOptions(2, 100*time.Millisecond, 5*time.Second)
	err := paymentClient.Post(
		fmt.Sprintf("%s/payments", s.config.PaymentServiceURL),
		paymentRequest,
		&paymentResponse,
	)
	if err != nil {
		log.Printf("Error processing payment: %v", err)
		return err
	}

	log.Printf("Payment processed. Payment ID: %s, Status: %s, Stripe Payment ID: %s",
		paymentResponse.ID, paymentResponse.Status, paymentResponse.StripePaymentID)

	// In a real implementation with frontend integration, we would return the client_secret
	// to the frontend to complete the payment with Stripe Elements

	return nil
}

// verifyCustomer verifies that a customer exists and is valid
func (s *OrderService) verifyCustomer(customerID string) error {
	// Send request to user service with timeout and retry
	var customerResponse struct {
		Verified bool   `json:"verified"`
		Message  string `json:"message"`
	}

	// Create a client with 2 retries as per design
	customerClient := utils.NewHTTPClientWithOptions(2, 100*time.Millisecond, 5*time.Second)

	// Prepare request body for user verification
	verifyRequest := struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Address string `json:"address,omitempty"`
	}{
		ID:    customerID,
		Email: "testuser@example.com", // Sử dụng email của người dùng test
	}

	// Gọi đúng endpoint /users/verify với phương thức POST
	err := customerClient.Post(
		fmt.Sprintf("%s/users/verify", s.config.UserServiceURL),
		verifyRequest,
		&customerResponse,
	)
	if err != nil {
		log.Printf("Error verifying customer: %v", err)
		return err
	}

	if !customerResponse.Verified {
		return errors.New("customer does not exist")
	}

	return nil
}

// getCartItems gets the items from a customer's cart
func (s *OrderService) getCartItems(customerID string) ([]models.OrderItem, error) {
	// Send request to cart service with timeout and retry
	var cartResponse struct {
		Items []struct {
			ProductID string  `json:"product_id"`
			Quantity  int     `json:"quantity"`
			Price     float64 `json:"price"`
		} `json:"items"`
	}

	// Create a client with 2 retries as per design
	cartClient := utils.NewHTTPClientWithOptions(2, 100*time.Millisecond, 5*time.Second)

	err := cartClient.Get(
		fmt.Sprintf("%s/carts/%s", s.config.CartServiceURL, customerID),
		&cartResponse,
	)
	if err != nil {
		log.Printf("Error getting cart items: %v", err)
		return nil, err
	}

	// Convert cart items to order items
	var orderItems []models.OrderItem
	for _, item := range cartResponse.Items {
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return orderItems, nil
}

// clearCart clears a customer's cart after order is processed or compensated
func (s *OrderService) clearCart(customerID string) error {
	// Create a client with 2 retries as per design
	cartClient := utils.NewHTTPClientWithOptions(2, 100*time.Millisecond, 5*time.Second)

	err := cartClient.Post(
		fmt.Sprintf("%s/carts/%s/clear", s.config.CartServiceURL, customerID),
		nil,
		nil,
	)
	if err != nil {
		log.Printf("Error clearing cart: %v", err)
		return err
	}

	return nil
}

// scheduleShipping schedules shipping for an order
func (s *OrderService) scheduleShipping(order models.Order) error {
	// Prepare request
	shipmentRequest := models.CreateShipmentRequest{
		OrderID:         order.ID,
		ShippingAddress: order.ShippingAddress,
	}

	// Create a special HTTP client with 3 retries for shipments as per design
	shipmentClient := utils.NewHTTPClientWithOptions(3, 1*time.Second, 5*time.Second)

	// Send request to shipping service with timeout and retry
	err := shipmentClient.Post(
		fmt.Sprintf("%s/shipments", s.config.ShippingServiceURL),
		shipmentRequest,
		nil,
	)
	if err != nil {
		log.Printf("Error scheduling shipping: %v", err)
		return err
	}

	return nil
}

// Compensate performs compensation actions when an order fails
func (s *OrderService) Compensate(order models.Order, failureReason string) error {
	log.Printf("Compensating for order %s with reason: %s", order.ID, failureReason)

	// 1. Update order status and failure reason
	order.Status = models.OrderStatusFailed
	order.FailureReason = failureReason
	order.UpdatedAt = time.Now()

	err := s.repository.UpdateOrder(order)
	if err != nil {
		log.Printf("Failed to update order status during compensation: %v", err)
		// Continue with other compensation actions
	}

	// 2. Restore inventory if it was locked
	if order.InventoryLocked {
		log.Printf("Restoring inventory for order %s", order.ID)
		for _, item := range order.Items {
			// Call Inventory Service to restore inventory
			restoreRequest := models.InventoryRestoreRequest{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}

			inventoryClient := utils.NewHTTPClientWithOptions(2, 100*time.Millisecond, 5*time.Second)
			err := inventoryClient.Post(
				fmt.Sprintf("%s/inventory/restore", s.config.InventoryServiceURL),
				restoreRequest,
				nil,
			)
			if err != nil {
				log.Printf("Failed to restore inventory for product %s: %v",
					item.ProductID, err)
				// Continue with other items
			}
		}
	}

	// 3. Refund payment if it was processed
	if order.PaymentProcessed {
		log.Printf("Refunding payment for order %s", order.ID)
		// Call Payment Service to refund payment
		refundRequest := models.RefundRequest{
			OrderID: order.ID,
			Amount:  order.TotalAmount,
		}

		paymentClient := utils.NewHTTPClientWithOptions(2, 100*time.Millisecond, 5*time.Second)
		err := paymentClient.Post(
			fmt.Sprintf("%s/payments/refund", s.config.PaymentServiceURL),
			refundRequest,
			nil,
		)
		if err != nil {
			log.Printf("Failed to refund payment for order %s: %v",
				order.ID, err)
		}
	}

	// 4. Send notification to customer
	customerEmail := "customer@example.com" // In a real app, get from Customer Service

	var notificationContent string
	if failureReason == "payment_failed" {
		// Include information about how to retry payment
		notificationContent = fmt.Sprintf("Your order %s has failed: %s. "+
			"You can retry payment by visiting http://example.com/orders/%s/payment or "+
			"by calling our customer service at 1-800-123-4567.",
			order.ID, getFailureMessage(failureReason), order.ID)
	} else {
		notificationContent = fmt.Sprintf("Your order %s has failed: %s",
			order.ID, getFailureMessage(failureReason))
	}

	err = s.sendNotification(order.CustomerID, customerEmail, notificationContent)
	if err != nil {
		log.Printf("Failed to send notification for order %s: %v",
			order.ID, err)
	}

	// 5. Clear cart (as per design)
	log.Printf("Clearing cart for customer %s", order.CustomerID)
	err = s.clearCart(order.CustomerID)
	if err != nil {
		log.Printf("Failed to clear cart for customer %s: %v", order.CustomerID, err)
		// Continue with other compensation actions
	}

	// 6. Publish order_cancelled event
	orderEvent := models.OrderEvent{
		EventType:     "order_cancelled",
		OrderID:       order.ID,
		CustomerID:    order.CustomerID,
		Status:        models.OrderStatusFailed,
		TotalAmount:   order.TotalAmount,
		Timestamp:     time.Now().Unix(),
		FailureReason: failureReason,
	}

	err = s.producer.PublishOrderEvent(orderEvent)
	if err != nil {
		log.Printf("Failed to publish order_cancelled event: %v", err)
	}

	return nil
}

// getFailureMessage returns a user-friendly message for a failure reason
func getFailureMessage(reason string) string {
	switch reason {
	case "inventory_unavailable":
		return "Some items in your order are not available in inventory"
	case "payment_failed":
		return "Payment processing failed"
	case "shipping_failed":
		return "Failed to schedule shipping"
	default:
		return "An error occurred while processing your order"
	}
}

// calculateTotalAmount calculates the total amount of an order
func calculateTotalAmount(items []models.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}

// RetryPayment retries payment for a failed order
func (s *OrderService) RetryPayment(orderID string, req models.RetryPaymentRequest) (models.Order, error) {
	// Get order
	order, err := s.repository.GetOrderByID(orderID)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to get order: %v", err)
	}

	// Check if order is in FAILED state with payment_failed reason
	if order.Status != models.OrderStatusFailed {
		return models.Order{}, errors.New("order is not in FAILED state")
	}

	if order.FailureReason != "payment_failed" {
		return models.Order{}, errors.New("order did not fail due to payment issues")
	}

	// Get customer information (in a real app, this would come from the user service)
	// For now, we'll use placeholder data
	customerEmail := "customer@example.com"
	customerName := "Test Customer"

	// Prepare payment request with more detailed information for Stripe
	paymentRequest := models.CreatePaymentRequest{
		OrderID:       order.ID,
		Amount:        order.TotalAmount,
		PaymentMethod: req.PaymentMethod,
		CardNumber:    req.CardNumber,
		ExpiryMonth:   req.ExpiryMonth,
		ExpiryYear:    req.ExpiryYear,
		CVV:           req.CVV,
		Currency:      "vnd", // Vietnamese Dong
		Description:   fmt.Sprintf("Payment retry for order %s", order.ID),
		CustomerEmail: customerEmail,
		CustomerName:  customerName,
	}

	log.Printf("Sending payment retry request to payment service: %+v", paymentRequest)

	// Process payment (2 retries as per design)
	var paymentResponse struct {
		ID                string `json:"id"`
		OrderID           string `json:"order_id"`
		Status            string `json:"status"`
		StripePaymentID   string `json:"stripe_payment_id,omitempty"`
		StripeClientSecret string `json:"stripe_client_secret,omitempty"`
	}

	paymentClient := utils.NewHTTPClientWithOptions(2, 100*time.Millisecond, 5*time.Second)
	err = paymentClient.Post(
		fmt.Sprintf("%s/payments", s.config.PaymentServiceURL),
		paymentRequest,
		&paymentResponse,
	)
	if err != nil {
		log.Printf("Error retrying payment: %v", err)
		return order, fmt.Errorf("failed to process payment: %v", err)
	}

	log.Printf("Payment retry processed. Payment ID: %s, Status: %s, Stripe Payment ID: %s",
		paymentResponse.ID, paymentResponse.Status, paymentResponse.StripePaymentID)

	// Update order status
	order.Status = models.OrderStatusConfirmed
	order.PaymentProcessed = true
	order.FailureReason = ""
	order.UpdatedAt = time.Now()

	err = s.repository.UpdateOrder(order)
	if err != nil {
		log.Printf("Failed to update order: %v", err)
		return order, fmt.Errorf("failed to update order: %v", err)
	}

	// Continue with order processing (schedule shipping)
	err = s.scheduleShipping(order)
	if err != nil {
		log.Printf("Failed to schedule shipping: %v", err)

		// Create audit log for shipping failure
		auditLog := models.AuditLog{
			ID:          uuid.New().String(),
			ServiceName: "order-service",
			Action:      models.AuditLogActionShipmentUpdate,
			CustomerID:  order.CustomerID,
			Timestamp:   time.Now(),
			Details:     fmt.Sprintf("Failed to schedule shipping for order ID: %s after payment retry, Error: %v", order.ID, err),
		}

		s.repository.CreateAuditLog(auditLog)

		// Continue anyway
	} else {
		// Create audit log for successful shipping schedule
		auditLog := models.AuditLog{
			ID:          uuid.New().String(),
			ServiceName: "order-service",
			Action:      models.AuditLogActionShipmentUpdate,
			CustomerID:  order.CustomerID,
			Timestamp:   time.Now(),
			Details:     fmt.Sprintf("Shipping scheduled for order ID: %s after payment retry", order.ID),
		}

		s.repository.CreateAuditLog(auditLog)
	}

	// Send notification to customer
	// Use the same customerEmail that was used for payment
	notificationContent := fmt.Sprintf("Your payment for order %s has been successfully processed. Your order is now confirmed.", order.ID)

	err = s.sendNotification(order.CustomerID, customerEmail, notificationContent)
	if err != nil {
		log.Printf("Failed to send notification for order %s: %v", order.ID, err)
		// Continue anyway
	}

	return order, nil
}
