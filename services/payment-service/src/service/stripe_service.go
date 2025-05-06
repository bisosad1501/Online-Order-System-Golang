package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/online-order-system/payment-service/config"
	"github.com/online-order-system/payment-service/db"
	"github.com/online-order-system/payment-service/interfaces"
	"github.com/online-order-system/payment-service/models"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/webhook"
)

// StripePaymentService handles payment processing with Stripe
type StripePaymentService struct {
	config     *config.Config
	repository *db.PaymentRepository
	producer   interfaces.PaymentProducer
}

// NewStripePaymentService creates a new Stripe payment service
func NewStripePaymentService(cfg *config.Config, repo *db.PaymentRepository, producer interfaces.PaymentProducer) *StripePaymentService {
	// Initialize Stripe API key
	stripe.Key = cfg.StripeSecretKey

	return &StripePaymentService{
		config:     cfg,
		repository: repo,
		producer:   producer,
	}
}

// ProcessPayment processes a payment with Stripe
func (s *StripePaymentService) ProcessPayment(payment models.Payment) (bool, error) {
	log.Printf("Processing payment with Stripe: %+v", payment)

	// Convert amount to cents (Stripe requires amount in smallest currency unit)
	amountInCents := int64(payment.Amount * 100)

	// Set default currency if not provided
	currency := payment.Currency
	if currency == "" {
		currency = "vnd" // Default to VND for Vietnamese currency
	}

	// Create description if not provided
	description := payment.Description
	if description == "" {
		description = fmt.Sprintf("Payment for order %s", payment.OrderID)
	}

	// Create payment intent parameters
	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(amountInCents),
		Currency:    stripe.String(currency),
		Description: stripe.String(description),
		// Set automatic payment methods to enable all available payment methods
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// Add metadata
	params.AddMetadata("order_id", payment.OrderID)

	// If customer email is provided, add receipt email
	if payment.CustomerEmail != "" {
		params.ReceiptEmail = stripe.String(payment.CustomerEmail)
	}

	// Create payment intent without confirming it
	// This allows the frontend to handle the confirmation with Stripe Elements
	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Error creating Stripe payment intent: %v", err)
		return false, fmt.Errorf("failed to create payment intent: %v", err)
	}

	log.Printf("Created Stripe payment intent: %s with client_secret: %s", pi.ID, pi.ClientSecret)

	// Update payment with Stripe-specific information
	err = s.repository.UpdatePaymentStripeInfo(payment.ID, pi.ID, pi.ClientSecret)
	if err != nil {
		log.Printf("Error updating payment with Stripe info: %v", err)
		// Continue anyway, as the payment intent was created
	}

	// For backend-only integration, we would check the status here
	// But since we're using a frontend integration, we'll return the payment as pending
	// The actual payment status will be updated via webhook when the frontend completes the payment

	// Return false to indicate payment is pending and needs frontend action
	// This is not an error condition, just part of the flow
	return false, nil
}

// HandleStripeWebhook processes a webhook event from Stripe
func (s *StripePaymentService) HandleStripeWebhook(payload []byte, signature string) error {
	log.Printf("Received Stripe webhook with signature: %s", signature)

	// Verify webhook signature
	event, err := webhook.ConstructEvent(payload, signature, s.config.StripeWebhookSecret)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		return fmt.Errorf("webhook signature verification failed: %v", err)
	}

	log.Printf("Webhook signature verified successfully. Event type: %s", event.Type)

	// Process different event types
	switch event.Type {
	case "payment_intent.succeeded":
		log.Printf("Processing payment_intent.succeeded event")
		return s.handlePaymentIntentSucceeded(event)

	case "payment_intent.payment_failed":
		log.Printf("Processing payment_intent.payment_failed event")
		return s.handlePaymentIntentFailed(event)

	case "payment_intent.created":
		log.Printf("Received payment_intent.created event - no action needed")
		return nil

	case "payment_intent.processing":
		log.Printf("Received payment_intent.processing event - payment is being processed")
		return nil

	case "payment_intent.requires_action":
		log.Printf("Received payment_intent.requires_action event - additional action required")
		return nil

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return nil
}

// handlePaymentIntentSucceeded processes a successful payment intent
func (s *StripePaymentService) handlePaymentIntentSucceeded(event stripe.Event) error {
	// Parse the payment intent from the event
	var pi stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &pi)
	if err != nil {
		log.Printf("Error parsing payment intent: %v", err)
		return fmt.Errorf("error parsing payment intent: %v", err)
	}

	log.Printf("Processing successful payment intent: %s", pi.ID)

	// Get order ID from metadata
	orderID, ok := pi.Metadata["order_id"]
	if !ok {
		log.Printf("Payment intent %s does not have order_id in metadata", pi.ID)
		return errors.New("payment intent does not have order_id in metadata")
	}

	log.Printf("Found order ID in metadata: %s", orderID)

	// Get payment by order ID
	payment, err := s.repository.GetPaymentByOrderID(orderID)
	if err != nil {
		log.Printf("Error getting payment for order %s: %v", orderID, err)
		return fmt.Errorf("error getting payment for order %s: %v", orderID, err)
	}

	log.Printf("Found payment with ID: %s for order: %s", payment.ID, orderID)

	// Update payment with additional Stripe information
	var receiptURL string
	if pi.LatestCharge != nil {
		// If we have a charge, get the receipt URL
		receiptURL = pi.LatestCharge.ReceiptURL
		if receiptURL != "" {
			log.Printf("Updating payment with receipt URL: %s", receiptURL)
			err = s.repository.UpdatePaymentReceiptURL(payment.ID, receiptURL)
			if err != nil {
				log.Printf("Error updating payment receipt URL: %v", err)
				// Continue anyway
			}
		}
	}

	// Update payment status
	log.Printf("Updating payment status to SUCCESSFUL")
	err = s.repository.UpdatePaymentStatus(payment.ID, models.PaymentStatusSuccessful)
	if err != nil {
		log.Printf("Error updating payment status: %v", err)
		return fmt.Errorf("error updating payment status: %v", err)
	}

	// Get updated payment
	payment, err = s.repository.GetPaymentByID(payment.ID)
	if err != nil {
		log.Printf("Error getting updated payment: %v", err)
		return fmt.Errorf("error getting updated payment: %v", err)
	}

	// Publish payment successful event
	log.Printf("Publishing payment_successful event for payment: %s, order: %s", payment.ID, payment.OrderID)
	err = s.producer.PublishPaymentSuccessful(payment)
	if err != nil {
		log.Printf("Error publishing payment successful event: %v", err)
		// Continue anyway
	}

	log.Printf("Successfully processed payment intent succeeded event")
	return nil
}

// handlePaymentIntentFailed processes a failed payment intent
func (s *StripePaymentService) handlePaymentIntentFailed(event stripe.Event) error {
	// Parse the payment intent from the event
	var pi stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &pi)
	if err != nil {
		log.Printf("Error parsing payment intent: %v", err)
		return fmt.Errorf("error parsing payment intent: %v", err)
	}

	log.Printf("Processing failed payment intent: %s", pi.ID)

	// Get order ID from metadata
	orderID, ok := pi.Metadata["order_id"]
	if !ok {
		log.Printf("Payment intent %s does not have order_id in metadata", pi.ID)
		return errors.New("payment intent does not have order_id in metadata")
	}

	log.Printf("Found order ID in metadata: %s", orderID)

	// Get payment by order ID
	payment, err := s.repository.GetPaymentByOrderID(orderID)
	if err != nil {
		log.Printf("Error getting payment for order %s: %v", orderID, err)
		return fmt.Errorf("error getting payment for order %s: %v", orderID, err)
	}

	log.Printf("Found payment with ID: %s for order: %s", payment.ID, orderID)

	// Get detailed error message if available
	var errorMessage string
	var errorCode string
	if pi.LastPaymentError != nil {
		errorMessage = pi.LastPaymentError.Error()
		errorCode = string(pi.LastPaymentError.Code)
		log.Printf("Payment error details - Code: %s, Message: %s", errorCode, errorMessage)
	} else {
		errorMessage = "Payment failed without specific error details"
		log.Printf("No specific error details available")
	}

	// Create a more user-friendly error message
	userFriendlyError := getUserFriendlyErrorMessage(errorCode, errorMessage)
	log.Printf("User-friendly error message: %s", userFriendlyError)

	// Update payment status and error message
	log.Printf("Updating payment status to FAILED with error message")
	err = s.repository.UpdatePaymentStatusWithError(payment.ID, models.PaymentStatusFailed, userFriendlyError)
	if err != nil {
		log.Printf("Error updating payment status: %v", err)
		return fmt.Errorf("error updating payment status: %v", err)
	}

	// Get updated payment
	payment, err = s.repository.GetPaymentByID(payment.ID)
	if err != nil {
		log.Printf("Error getting updated payment: %v", err)
		return fmt.Errorf("error getting updated payment: %v", err)
	}

	// Publish payment failed event
	log.Printf("Publishing payment_failed event for payment: %s, order: %s", payment.ID, payment.OrderID)
	err = s.producer.PublishPaymentFailed(payment)
	if err != nil {
		log.Printf("Error publishing payment failed event: %v", err)
		// Continue anyway
	}

	log.Printf("Successfully processed payment intent failed event")
	return nil
}

// ConfirmPayment confirms a payment with Stripe
func (s *StripePaymentService) ConfirmPayment(payment models.Payment) (bool, error) {
	log.Printf("Confirming payment with Stripe: %+v", payment)

	// Get the payment intent
	pi, err := paymentintent.Get(payment.StripePaymentID, nil)
	if err != nil {
		log.Printf("Error getting Stripe payment intent: %v", err)
		return false, fmt.Errorf("failed to get payment intent: %v", err)
	}

	// Check if payment intent is already succeeded
	if pi.Status == stripe.PaymentIntentStatusSucceeded {
		log.Printf("Payment intent already succeeded: %s", pi.ID)

		// Update payment status
		err = s.repository.UpdatePaymentStatus(payment.ID, models.PaymentStatusSuccessful)
		if err != nil {
			log.Printf("Error updating payment status: %v", err)
			return false, fmt.Errorf("error updating payment status: %v", err)
		}

		return true, nil
	}

	// For testing purposes, we'll use a test card based on the card number
	params := &stripe.PaymentIntentConfirmParams{
		ReturnURL: stripe.String("http://localhost:3000/payment/success"),
	}

	// Create a payment method based on the card details
	if payment.CardNumber != "" {
		// Use the provided card details
		pmParams := &stripe.PaymentMethodParams{
			Card: &stripe.PaymentMethodCardParams{
				Number:   stripe.String(payment.CardNumber),
				ExpMonth: stripe.String(payment.ExpiryMonth),
				ExpYear:  stripe.String(payment.ExpiryYear),
				CVC:      stripe.String(payment.CVV),
			},
			Type: stripe.String("card"),
		}

		pm, err := paymentmethod.New(pmParams)
		if err != nil {
			log.Printf("Error creating payment method: %v", err)
			return payment, fmt.Errorf("failed to create payment method: %v", err)
		}

		params.PaymentMethod = stripe.String(pm.ID)
	} else if payment.PaymentMethod == "STRIPE_TEST_SUCCESS" {
		// Use a test card that always succeeds
		params.PaymentMethod = stripe.String("pm_card_visa")
	} else {
		// Default to a test card that always fails with insufficient funds
		params.PaymentMethod = stripe.String("pm_card_chargeCustomerFail")
	}

	// Confirm the payment intent
	pi, err = paymentintent.Confirm(payment.StripePaymentID, params)
	if err != nil {
		log.Printf("Error confirming Stripe payment intent: %v", err)

		// Update payment status to FAILED
		errorMessage := fmt.Sprintf("Payment failed: %v", err)
		updateErr := s.repository.UpdatePaymentStatusWithError(payment.ID, models.PaymentStatusFailed, errorMessage)
		if updateErr != nil {
			log.Printf("Error updating payment status: %v", updateErr)
			// Continue anyway
		}

		return false, fmt.Errorf("failed to confirm payment intent: %v", err)
	}

	log.Printf("Confirmed Stripe payment intent: %s with status: %s", pi.ID, pi.Status)

	// Check if payment intent succeeded
	if pi.Status == stripe.PaymentIntentStatusSucceeded {
		// Update payment status
		err = s.repository.UpdatePaymentStatus(payment.ID, models.PaymentStatusSuccessful)
		if err != nil {
			log.Printf("Error updating payment status: %v", err)
			return false, fmt.Errorf("error updating payment status: %v", err)
		}

		// Update receipt URL if available
		if pi.LatestCharge != nil && pi.LatestCharge.ReceiptURL != "" {
			err = s.repository.UpdatePaymentReceiptURL(payment.ID, pi.LatestCharge.ReceiptURL)
			if err != nil {
				log.Printf("Error updating payment receipt URL: %v", err)
				// Continue anyway
			}
		}

		return true, nil
	}

	// If payment intent requires action, return false
	if pi.Status == stripe.PaymentIntentStatusRequiresAction {
		log.Printf("Payment intent requires action: %s", pi.ID)
		return false, nil
	}

	// If payment intent failed, update payment status and return error
	if pi.Status == stripe.PaymentIntentStatusCanceled || pi.Status == stripe.PaymentIntentStatusRequiresPaymentMethod {
		// Get error message
		var errorMessage string
		if pi.LastPaymentError != nil {
			errorMessage = pi.LastPaymentError.Error()
		} else {
			errorMessage = "Payment failed without specific error details"
		}

		// Update payment status
		err = s.repository.UpdatePaymentStatusWithError(payment.ID, models.PaymentStatusFailed, errorMessage)
		if err != nil {
			log.Printf("Error updating payment status: %v", err)
			return false, fmt.Errorf("error updating payment status: %v", err)
		}

		return false, fmt.Errorf("payment failed: %s", errorMessage)
	}

	// For any other status, return false
	return false, fmt.Errorf("payment intent has unexpected status: %s", pi.Status)
}

// getUserFriendlyErrorMessage converts Stripe error codes to user-friendly messages
func getUserFriendlyErrorMessage(errorCode string, defaultMessage string) string {
	switch errorCode {
	case "card_declined":
		return "Your card was declined. Please try another payment method."
	case "expired_card":
		return "Your card has expired. Please try another card."
	case "incorrect_cvc":
		return "The security code (CVV) you entered is incorrect. Please check and try again."
	case "processing_error":
		return "An error occurred while processing your card. Please try again later."
	case "insufficient_funds":
		return "Your card has insufficient funds. Please try another payment method."
	case "invalid_card_number":
		return "The card number you entered is invalid. Please check and try again."
	case "invalid_expiry_month":
		return "The expiration month you entered is invalid. Please check and try again."
	case "invalid_expiry_year":
		return "The expiration year you entered is invalid. Please check and try again."
	default:
		if defaultMessage != "" {
			return defaultMessage
		}
		return "An error occurred while processing your payment. Please try again."
	}
}
