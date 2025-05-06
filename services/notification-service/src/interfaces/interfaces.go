package interfaces

import (
"github.com/online-order-system/notification-service/models"
)

// NotificationService defines the interface for notification service
type NotificationService interface {
CreateNotification(req models.CreateNotificationRequest) (models.Notification, error)
GetNotificationByID(id string) (models.Notification, error)
GetNotificationsByCustomerID(customerID string) ([]models.Notification, error)
GetNotifications() ([]models.Notification, error)
UpdateNotificationStatus(id string, status models.NotificationStatus) (models.Notification, error)
SendNotification(notification models.Notification) error
ProcessOrderEvent(event models.OrderEvent) error
ProcessPaymentEvent(event models.PaymentEvent) error
ProcessShipmentEvent(event models.ShipmentEvent) error
}

// NotificationProducer defines the interface for notification producer
type NotificationProducer interface {
PublishNotificationCreated(notification models.Notification) error
PublishNotificationSent(notification models.Notification) error
PublishNotificationFailed(notification models.Notification) error
Close() error
}

// NotificationConsumer defines the interface for notification consumer
type NotificationConsumer interface {
ProcessOrderEvent(event models.OrderEvent) error
ProcessPaymentEvent(event models.PaymentEvent) error
ProcessShipmentEvent(event models.ShipmentEvent) error
}

// EmailSender defines the interface for email sender
type EmailSender interface {
SendEmail(to string, subject string, content string) error
}

// SMSSender defines the interface for SMS sender
type SMSSender interface {
SendSMS(to string, content string) error
}

// PushSender defines the interface for push notification sender
type PushSender interface {
SendPush(to string, title string, content string) error
}
