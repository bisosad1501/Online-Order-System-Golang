package service

import (
"log"
"time"

"github.com/online-order-system/notification-service/config"
"github.com/online-order-system/notification-service/db"
"github.com/online-order-system/notification-service/interfaces"
"github.com/online-order-system/notification-service/models"
)

// NotificationService handles business logic for notifications
type NotificationService struct {
config     *config.Config
repository *db.NotificationRepository
producer   interfaces.NotificationProducer
orderClient *OrderClient
}

// Ensure NotificationService implements NotificationService interface
var _ interfaces.NotificationService = (*NotificationService)(nil)

// NewNotificationService creates a new notification service
func NewNotificationService(cfg *config.Config, repo *db.NotificationRepository, producer interfaces.NotificationProducer) *NotificationService {
return &NotificationService{
config:     cfg,
repository: repo,
producer:   producer,
orderClient: NewOrderClient(cfg),
}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(req models.CreateNotificationRequest) (models.Notification, error) {
// Create notification
now := time.Now()
notification := models.Notification{
ID:         db.GenerateID(),
CustomerID: req.CustomerID,
Type:       req.Type,
Status:     models.NotificationStatusUnread, // Use UNREAD status for user-facing notifications
Subject:    req.Subject,
Content:    req.Content,
Recipient:  req.Recipient,
CreatedAt:  now,
UpdatedAt:  now,
}

// Save notification to database
err := s.repository.CreateNotification(notification)
if err != nil {
return models.Notification{}, err
}

// Publish notification created event
err = s.producer.PublishNotificationCreated(notification)
if err != nil {
log.Printf("Failed to publish notification created event: %v", err)
// Continue anyway
}

// Send notification
go func() {
err := s.SendNotification(notification)
if err != nil {
log.Printf("Failed to send notification: %v", err)
}
}()

return notification, nil
}

// GetNotificationByID retrieves a notification by ID
func (s *NotificationService) GetNotificationByID(id string) (models.Notification, error) {
return s.repository.GetNotificationByID(id)
}

// GetNotifications retrieves all notifications
func (s *NotificationService) GetNotifications() ([]models.Notification, error) {
return s.repository.GetNotifications()
}

// GetNotificationsByCustomerID retrieves all notifications for a customer
func (s *NotificationService) GetNotificationsByCustomerID(customerID string) ([]models.Notification, error) {
return s.repository.GetNotificationsByCustomerID(customerID)
}

// UpdateNotificationStatus updates the status of a notification
func (s *NotificationService) UpdateNotificationStatus(id string, status models.NotificationStatus) (models.Notification, error) {
// Update notification status
err := s.repository.UpdateNotificationStatus(id, status)
if err != nil {
return models.Notification{}, err
}

// Get updated notification
notification, err := s.GetNotificationByID(id)
if err != nil {
return models.Notification{}, err
}

// Publish notification status event when it's read
if status == models.NotificationStatusRead {
err = s.producer.PublishNotificationSent(notification)
if err != nil {
log.Printf("Failed to publish notification read event: %v", err)
// Continue anyway
}
}

return notification, nil
}

// SendNotification sends a notification
func (s *NotificationService) SendNotification(notification models.Notification) error {
// In a real application, this would send the notification via email, SMS, etc.
// For this demo, we'll just log it
log.Printf("Sending notification: %s to %s", notification.Subject, notification.Recipient)

// Mark notification as READ after sending
_, err := s.UpdateNotificationStatus(notification.ID, models.NotificationStatusRead)
if err != nil {
return err
}
return nil
}

// ProcessOrderEvent processes an order event
func (s *NotificationService) ProcessOrderEvent(event models.OrderEvent) error {
    // Create notification based on order event
    var notificationType models.NotificationType
    var subject, content string

    switch event.EventType {
    // Bước 8: Xác nhận đơn hàng
    case "order_confirmed":
        notificationType = models.NotificationTypeOrderConfirmed
        subject = "Order Confirmed"
        content = "Your order " + event.OrderID + " has been confirmed and is being prepared for shipping."
        log.Printf("Creating notification for order_confirmed event: %s", event.OrderID)

    // Bước 11: Thông báo trạng thái giao hàng
    case "shipping_status_updated":
        notificationType = models.NotificationTypeShippingUpdate
        subject = "Shipment Status Updated"
        content = "The status of your shipment for order " + event.OrderID + " has been updated to: " + event.Status
        log.Printf("Creating notification for shipping_status_updated event: %s, status: %s", event.OrderID, event.Status)

    // Bước 13: Thông báo hoàn thành đơn hàng
    case "shipping_completed":
        notificationType = models.NotificationTypeShippingDone
        subject = "Order Delivered"
        content = "Your order " + event.OrderID + " has been delivered. Thank you for shopping with us!"
        log.Printf("Creating notification for shipping_completed event: %s", event.OrderID)

    // Thông báo khi đơn hàng hoàn thành (từ order-service)
    case "order_completed":
        notificationType = models.NotificationTypeOrderCompleted
        subject = "Order Completed"
        content = "Your order " + event.OrderID + " has been completed. Thank you for shopping with us!"
        log.Printf("Creating notification for order_completed event: %s", event.OrderID)

    // Thông báo khi đơn hàng bị hủy
    case "order_cancelled":
        notificationType = models.NotificationTypeOrderCancelled
        subject = "Order Cancelled"

        // Kiểm tra lý do hủy đơn hàng
        reason := "Unknown reason"
        if event.FailureReason == "inventory_unavailable" {
            reason = "Some items in your order are not available in inventory"
        } else if event.FailureReason == "payment_failed" {
            reason = "Payment processing failed. You can retry payment by visiting our website or contacting customer service."
        }

        content = "Your order " + event.OrderID + " has been cancelled: " + reason
        log.Printf("Creating notification for order_cancelled event: %s, reason: %s", event.OrderID, reason)

    default:
        log.Printf("Ignoring order event type: %s", event.EventType)
        return nil // Ignore unknown event types
    }

// Create notification request
req := models.CreateNotificationRequest{
CustomerID: event.CustomerID,
Type:      notificationType,
Subject:   subject,
Content:   content,
Recipient: event.CustomerID + "@example.com", // In a real app, get from customer service
}

// Create notification
_, err := s.CreateNotification(req)
return err
}

// ProcessPaymentEvent processes a payment event
func (s *NotificationService) ProcessPaymentEvent(event models.PaymentEvent) error {
    // Create notification based on payment event
    var notificationType models.NotificationType
    var subject, content string

    switch event.EventType {
    // Bước 6-7: Xử lý thanh toán
    case "payment_succeeded":
        notificationType = models.NotificationTypeOrderConfirmed
        subject = "Payment Successful"
        content = "Your payment for order " + event.OrderID + " has been successfully processed."
        log.Printf("Creating notification for payment_succeeded event: %s", event.OrderID)

    // Bước 7: Thông báo lỗi nếu thanh toán thất bại
    case "payment_failed":
        notificationType = models.NotificationTypePaymentFailed
        subject = "Payment Failed"
        content = "Your payment for order " + event.OrderID + " has failed. Please try again or contact customer service."
        log.Printf("Creating notification for payment_failed event: %s", event.OrderID)

    default:
        log.Printf("Ignoring payment event type: %s", event.EventType)
        return nil // Ignore unknown event types
    }

// Create notification request
req := models.CreateNotificationRequest{
CustomerID: event.CustomerID,
Type:      notificationType,
Subject:   subject,
Content:   content,
Recipient: event.CustomerID + "@example.com", // In a real app, get from customer service
}

// Create notification
_, err := s.CreateNotification(req)
return err
}

// ProcessShipmentEvent processes a shipment event
func (s *NotificationService) ProcessShipmentEvent(event models.ShipmentEvent) error {
    // Create notification based on shipment event
    var notificationType models.NotificationType
    var subject, content string

    switch event.EventType {
    // Bước 9: Lập lịch giao hàng
    case "shipping_scheduled":
        notificationType = models.NotificationTypeShippingUpdate
        subject = "Shipment Scheduled"
        content = "Your order " + event.OrderID + " has been scheduled for shipping."
        log.Printf("Creating notification for shipping_scheduled event: %s", event.OrderID)

    // Bước 11: Thông báo trạng thái giao hàng
    case "shipping_status_updated":
        notificationType = models.NotificationTypeShippingUpdate
        subject = "Shipment Updated"
        content = "Your shipment for order " + event.OrderID + " has been updated. Status: " + event.Status
        log.Printf("Creating notification for shipping_status_updated event: %s, status: %s", event.OrderID, event.Status)

    // Bước 13: Thông báo hoàn thành đơn hàng
    case "shipping_completed":
        notificationType = models.NotificationTypeShippingDone
        subject = "Shipment Completed"
        content = "Your order " + event.OrderID + " has been delivered."
        log.Printf("Creating notification for shipping_completed event: %s", event.OrderID)

    default:
        log.Printf("Ignoring shipment event type: %s", event.EventType)
        return nil // Ignore unknown event types
    }

// Check if we have customer ID in the event
if event.CustomerID == "" {
    // If not, try to get it from order service
    order, err := s.orderClient.GetOrderByID(event.OrderID)
    if err != nil {
        log.Printf("Failed to get order details: %v", err)
        // Continue with a fallback approach instead of returning error
        // We'll use the order ID as a placeholder for customer ID
        event.CustomerID = event.OrderID
    } else {
        event.CustomerID = order.CustomerID
    }
}

// Create notification request with customer ID
req := models.CreateNotificationRequest{
    CustomerID: event.CustomerID,
    Type:      notificationType,
    Subject:   subject,
    Content:   content,
    Recipient: event.CustomerID + "@example.com", // Use customer ID for email
}

// Create notification
_, err := s.CreateNotification(req)
return err
}
