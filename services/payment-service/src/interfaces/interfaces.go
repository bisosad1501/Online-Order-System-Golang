package interfaces

import (
"github.com/online-order-system/payment-service/models"
)

// PaymentService defines the interface for payment service
type PaymentService interface {
CreatePayment(req models.CreatePaymentRequest) (models.Payment, error)
GetPaymentByID(id string) (models.Payment, error)
GetPaymentByOrderID(orderID string) (models.Payment, error)
GetPayments() ([]models.Payment, error)
UpdatePaymentStatus(id string, status models.PaymentStatus) (models.Payment, error)
HandleStripeWebhook(payload []byte, signature string) error
ConfirmPayment(paymentID string) (models.Payment, error)
}

// PaymentProducer defines the interface for payment producer
type PaymentProducer interface {
PublishPaymentCreated(payment models.Payment) error
PublishPaymentSuccessful(payment models.Payment) error
PublishPaymentFailed(payment models.Payment) error
PublishPaymentRefunded(payment models.Payment) error
Close() error
}
