package service

import (
"errors"
"fmt"
"log"
"math/rand"
"time"

"github.com/online-order-system/payment-service/config"
"github.com/online-order-system/payment-service/db"
"github.com/online-order-system/payment-service/interfaces"
"github.com/online-order-system/payment-service/models"
)

// PaymentService handles business logic for payments
type PaymentService struct {
config     *config.Config
repository *db.PaymentRepository
producer   interfaces.PaymentProducer
stripeService *StripePaymentService
}

// Ensure PaymentService implements PaymentService interface
var _ interfaces.PaymentService = (*PaymentService)(nil)

// NewPaymentService creates a new payment service
func NewPaymentService(cfg *config.Config, repo *db.PaymentRepository, producer interfaces.PaymentProducer) *PaymentService {
	// Create Stripe service
	stripeService := NewStripePaymentService(cfg, repo, producer)

	return &PaymentService{
		config:        cfg,
		repository:    repo,
		producer:      producer,
		stripeService: stripeService,
	}
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(req models.CreatePaymentRequest) (models.Payment, error) {
// Log the request
log.Printf("Service: Creating payment with request: %+v", req)

// Normalize payment method to handle different naming conventions
paymentMethod := req.PaymentMethod
if paymentMethod == "credit_card" {
    paymentMethod = "card"
}

// Create payment
now := time.Now()
payment := models.Payment{
    ID:            db.GenerateID(),
    OrderID:       req.OrderID,
    Amount:        req.Amount,
    Status:        models.PaymentStatusPending,
    PaymentMethod: paymentMethod,
    CardNumber:    req.CardNumber,
    ExpiryMonth:   req.ExpiryMonth,
    ExpiryYear:    req.ExpiryYear,
    CVV:           req.CVV,
    CreatedAt:     now,
    UpdatedAt:     now,
    // Stripe specific fields
    Currency:      req.Currency,
    Description:   req.Description,
    CustomerEmail: req.CustomerEmail,
    CustomerName:  req.CustomerName,
    // For Stripe payment method
    PaymentMethodID: req.PaymentMethodID,
}

// Save payment to database
err := s.repository.CreatePayment(payment)
if err != nil {
    log.Printf("Error creating payment in database: %v", err)
    return models.Payment{}, err
}

// Process payment with payment gateway
success, err := s.processPayment(payment)
if err != nil {
    log.Printf("Error processing payment: %v", err)
    // Update payment status to failed
    updateErr := s.repository.UpdatePaymentStatus(payment.ID, models.PaymentStatusFailed)
    if updateErr != nil {
        log.Printf("Error updating payment status to failed: %v", updateErr)
    }

    // Publish payment failed event
    failedPayment, getErr := s.repository.GetPaymentByID(payment.ID)
    if getErr != nil {
        log.Printf("Error getting failed payment: %v", getErr)
    } else {
        publishErr := s.producer.PublishPaymentFailed(failedPayment)
        if publishErr != nil {
            log.Printf("Error publishing payment failed event: %v", publishErr)
        }
    }

    return models.Payment{}, err
}

// For Stripe integration, we don't immediately update the status to successful
// because the payment is processed asynchronously via webhooks
// Instead, we leave it as pending if we're using Stripe and success is false
var status models.PaymentStatus
if s.config.PaymentMode == "stripe" && !success {
    // For Stripe, if success is false, it means the payment is pending
    // and requires further action from the frontend
    status = models.PaymentStatusPending
    log.Printf("Stripe payment is pending and requires frontend action")
} else if success {
    status = models.PaymentStatusSuccessful
    log.Printf("Payment was successful")
} else {
    status = models.PaymentStatusFailed
    log.Printf("Payment failed")
}

// Update payment status in database
err = s.repository.UpdatePaymentStatus(payment.ID, status)
if err != nil {
    log.Printf("Error updating payment status: %v", err)
    return models.Payment{}, err
}

// Get updated payment
payment, err = s.repository.GetPaymentByID(payment.ID)
if err != nil {
    log.Printf("Error getting updated payment: %v", err)
    return models.Payment{}, err
}

// Publish event based on status (only for successful or failed payments)
if payment.Status == models.PaymentStatusSuccessful {
    err = s.producer.PublishPaymentSuccessful(payment)
    if err != nil {
        log.Printf("Failed to publish payment successful event: %v", err)
        // Continue anyway
    }
} else if payment.Status == models.PaymentStatusFailed {
    err = s.producer.PublishPaymentFailed(payment)
    if err != nil {
        log.Printf("Failed to publish payment failed event: %v", err)
        // Continue anyway
    }
}

log.Printf("Payment created with ID: %s, Status: %s, Stripe Payment ID: %s",
    payment.ID, payment.Status, payment.StripePaymentID)

return payment, nil
}

// GetPaymentByID retrieves a payment by ID
func (s *PaymentService) GetPaymentByID(id string) (models.Payment, error) {
return s.repository.GetPaymentByID(id)
}

// GetPaymentByOrderID retrieves a payment by order ID
func (s *PaymentService) GetPaymentByOrderID(orderID string) (models.Payment, error) {
return s.repository.GetPaymentByOrderID(orderID)
}

// GetPayments retrieves all payments
func (s *PaymentService) GetPayments() ([]models.Payment, error) {
return s.repository.GetPayments()
}

// UpdatePaymentStatus updates the status of a payment
func (s *PaymentService) UpdatePaymentStatus(id string, status models.PaymentStatus) (models.Payment, error) {
// Update payment status
err := s.repository.UpdatePaymentStatus(id, status)
if err != nil {
return models.Payment{}, err
}

// Get updated payment
payment, err := s.repository.GetPaymentByID(id)
if err != nil {
return models.Payment{}, err
}

// Publish event based on status
var publishErr error
switch status {
case models.PaymentStatusSuccessful:
publishErr = s.producer.PublishPaymentSuccessful(payment)
case models.PaymentStatusFailed:
publishErr = s.producer.PublishPaymentFailed(payment)
case models.PaymentStatusRefunded:
publishErr = s.producer.PublishPaymentRefunded(payment)
}

if publishErr != nil {
log.Printf("Failed to publish payment event: %v", publishErr)
// Continue anyway
}

return payment, nil
}

// processPayment processes a payment with a payment gateway
func (s *PaymentService) processPayment(payment models.Payment) (bool, error) {
	// Check payment mode
	if s.config.PaymentMode == "stripe" {
		// Use Stripe for payment processing
		return s.stripeService.ProcessPayment(payment)
	} else {
		// Use mock implementation for demonstration purposes
		return s.processMockPayment(payment)
	}
}

// processMockPayment is a mock implementation for demonstration purposes
func (s *PaymentService) processMockPayment(payment models.Payment) (bool, error) {
	// Simulate payment processing
	time.Sleep(500 * time.Millisecond)

	// Validate payment method
	log.Printf("Processing mock payment with method: '%s'", payment.PaymentMethod)

	// Accept both "card" and "credit_card" for consistency
	if payment.PaymentMethod != "card" && payment.PaymentMethod != "credit_card" && payment.PaymentMethod != "paypal" {
		return false, fmt.Errorf("unsupported payment method: %s", payment.PaymentMethod)
	}

	// Validate credit card details for card payments
	if payment.PaymentMethod == "card" || payment.PaymentMethod == "credit_card" {
		if payment.CardNumber == "" {
			log.Printf("Warning: Card number is empty for card payment")
			// Don't fail in mock mode, just log a warning
		}

		if payment.ExpiryMonth == "" || payment.ExpiryYear == "" || payment.CVV == "" {
			log.Printf("Warning: Card details incomplete (expiry or CVV missing)")
			// Don't fail in mock mode, just log a warning
		}
	}

	// Simulate success/failure (90% success rate)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	success := r.Float64() < 0.9

	if success {
		log.Printf("Mock payment successful for order %s", payment.OrderID)
	} else {
		log.Printf("Mock payment failed for order %s", payment.OrderID)
	}

	return success, nil
}

// HandleStripeWebhook processes a webhook event from Stripe
func (s *PaymentService) HandleStripeWebhook(payload []byte, signature string) error {
	if s.config.PaymentMode == "stripe" {
		return s.stripeService.HandleStripeWebhook(payload, signature)
	}
	return errors.New("stripe payment mode is not enabled")
}

// ConfirmPayment confirms a payment with Stripe
func (s *PaymentService) ConfirmPayment(paymentID string) (models.Payment, error) {
	// Get payment
	payment, err := s.repository.GetPaymentByID(paymentID)
	if err != nil {
		log.Printf("Error getting payment: %v", err)
		return models.Payment{}, fmt.Errorf("failed to get payment: %v", err)
	}

	// Check if payment is in PENDING state
	if payment.Status != models.PaymentStatusPending {
		log.Printf("Payment is not in PENDING state: %s", payment.Status)
		return payment, fmt.Errorf("payment is not in PENDING state: %s", payment.Status)
	}

	// Check if payment has Stripe payment ID
	if payment.StripePaymentID == "" {
		log.Printf("Payment does not have Stripe payment ID")
		return payment, errors.New("payment does not have Stripe payment ID")
	}

	// Confirm payment with Stripe
	if s.config.PaymentMode == "stripe" {
		// Use Stripe for payment confirmation
		success, err := s.stripeService.ConfirmPayment(payment)
		if err != nil {
			log.Printf("Error confirming payment with Stripe: %v", err)

			// Get updated payment after error (it should be in FAILED state now)
			updatedPayment, getErr := s.repository.GetPaymentByID(paymentID)
			if getErr != nil {
				log.Printf("Error getting updated payment after failure: %v", getErr)
				return payment, fmt.Errorf("failed to confirm payment with Stripe: %v", err)
			}

			// Publish payment failed event
			if updatedPayment.Status == models.PaymentStatusFailed {
				publishErr := s.producer.PublishPaymentFailed(updatedPayment)
				if publishErr != nil {
					log.Printf("Failed to publish payment failed event: %v", publishErr)
					// Continue anyway
				} else {
					log.Printf("Published payment failed event for payment %s", updatedPayment.ID)
				}
			}

			return updatedPayment, fmt.Errorf("failed to confirm payment with Stripe: %v", err)
		}

		if !success {
			log.Printf("Payment confirmation failed")
			return payment, errors.New("payment confirmation failed")
		}
	} else {
		// Use mock implementation for demonstration purposes
		// Update payment status to successful
		err = s.repository.UpdatePaymentStatus(payment.ID, models.PaymentStatusSuccessful)
		if err != nil {
			log.Printf("Error updating payment status: %v", err)
			return payment, fmt.Errorf("failed to update payment status: %v", err)
		}
	}

	// Get updated payment
	payment, err = s.repository.GetPaymentByID(paymentID)
	if err != nil {
		log.Printf("Error getting updated payment: %v", err)
		return models.Payment{}, fmt.Errorf("failed to get updated payment: %v", err)
	}

	// Publish payment successful event
	err = s.producer.PublishPaymentSuccessful(payment)
	if err != nil {
		log.Printf("Failed to publish payment successful event: %v", err)
		// Continue anyway
	} else {
		log.Printf("Published payment successful event for payment %s", payment.ID)
	}

	return payment, nil
}
