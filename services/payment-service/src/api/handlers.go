package api

import (
"log"
"net/http"

"github.com/gin-gonic/gin"
"github.com/online-order-system/payment-service/interfaces"
"github.com/online-order-system/payment-service/models"
)

// Handler handles HTTP requests
type Handler struct {
service interfaces.PaymentService
}

// NewHandler creates a new handler
func NewHandler(service interfaces.PaymentService) *Handler {
return &Handler{
service: service,
}
}

// CreatePayment handles the creation of a new payment
func (h *Handler) CreatePayment(c *gin.Context) {
var req models.CreatePaymentRequest
if err := c.ShouldBindJSON(&req); err != nil {
log.Printf("Error binding JSON: %v", err)
log.Printf("Request body: %v", c.Request.Body)
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

log.Printf("Creating payment with request: %+v", req)
log.Printf("Payment method (raw): '%s'", req.PaymentMethod)
log.Printf("Payment method length: %d", len(req.PaymentMethod))

// Print each character and its ASCII code
for i, char := range req.PaymentMethod {
	log.Printf("Character %d: '%c' (ASCII: %d)", i, char, char)
}

payment, err := h.service.CreatePayment(req)
if err != nil {
log.Printf("Error creating payment: %v", err)
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

// Mask sensitive data
payment.CardNumber = maskCardNumber(payment.CardNumber)
payment.CVV = "***"

c.JSON(http.StatusCreated, payment)
}

// GetPaymentByID handles retrieving a payment by ID
func (h *Handler) GetPaymentByID(c *gin.Context) {
id := c.Param("id")

payment, err := h.service.GetPaymentByID(id)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
return
}

// Mask sensitive data
payment.CardNumber = maskCardNumber(payment.CardNumber)
payment.CVV = "***"

c.JSON(http.StatusOK, payment)
}

// GetPaymentByOrderID handles retrieving a payment by order ID
func (h *Handler) GetPaymentByOrderID(c *gin.Context) {
orderID := c.Param("order_id")

payment, err := h.service.GetPaymentByOrderID(orderID)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
return
}

// Mask sensitive data
payment.CardNumber = maskCardNumber(payment.CardNumber)
payment.CVV = "***"

c.JSON(http.StatusOK, payment)
}

// GetPayments handles retrieving all payments
func (h *Handler) GetPayments(c *gin.Context) {
payments, err := h.service.GetPayments()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments"})
return
}

// Mask sensitive data
for i := range payments {
payments[i].CardNumber = maskCardNumber(payments[i].CardNumber)
payments[i].CVV = "***"
}

c.JSON(http.StatusOK, payments)
}

// UpdatePaymentStatus handles updating a payment's status
func (h *Handler) UpdatePaymentStatus(c *gin.Context) {
id := c.Param("id")

var req models.UpdatePaymentStatusRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

payment, err := h.service.UpdatePaymentStatus(id, req.Status)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

// Mask sensitive data
payment.CardNumber = maskCardNumber(payment.CardNumber)
payment.CVV = "***"

c.JSON(http.StatusOK, payment)
}

// HandleStripeWebhook handles webhook events from Stripe
func (h *Handler) HandleStripeWebhook(c *gin.Context) {
	// Read the request body
	payload, err := c.GetRawData()
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	// Get the Stripe signature from the header
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		log.Printf("Missing Stripe-Signature header")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe-Signature header"})
		return
	}

	// Process the webhook
	err = h.service.HandleStripeWebhook(payload, signature)
	if err != nil {
		log.Printf("Error processing Stripe webhook: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"received": true})
}

// ConfirmPayment handles confirming a payment
func (h *Handler) ConfirmPayment(c *gin.Context) {
	id := c.Param("id")

	// Get the payment first
	payment, err := h.service.GetPaymentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Check if card details are provided
	var cardDetails struct {
		CardNumber  string `json:"card_number"`
		ExpiryMonth string `json:"expiry_month"`
		ExpiryYear  string `json:"expiry_year"`
		CVV         string `json:"cvv"`
	}

	if err := c.ShouldBindJSON(&cardDetails); err == nil && cardDetails.CardNumber != "" {
		// Update payment with card details
		payment.CardNumber = cardDetails.CardNumber
		payment.ExpiryMonth = cardDetails.ExpiryMonth
		payment.ExpiryYear = cardDetails.ExpiryYear
		payment.CVV = cardDetails.CVV
	}

	// Confirm the payment
	payment, err = h.service.ConfirmPayment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mask sensitive data
	payment.CardNumber = maskCardNumber(payment.CardNumber)
	payment.CVV = "***"

	c.JSON(http.StatusOK, payment)
}

// TestSuccessfulPayment handles testing a successful payment with Stripe
func (h *Handler) TestSuccessfulPayment(c *gin.Context) {
	id := c.Param("id")

	// Get the payment first
	payment, err := h.service.GetPaymentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Set payment method to STRIPE_TEST_SUCCESS to trigger successful payment
	payment.PaymentMethod = "STRIPE_TEST_SUCCESS"

	// Confirm the payment
	payment, err = h.service.ConfirmPayment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mask sensitive data
	payment.CardNumber = maskCardNumber(payment.CardNumber)
	payment.CVV = "***"

	c.JSON(http.StatusOK, payment)
}

// maskCardNumber masks a credit card number, showing only the last 4 digits
func maskCardNumber(cardNumber string) string {
	if len(cardNumber) <= 4 {
		return cardNumber
	}
	return "************" + cardNumber[len(cardNumber)-4:]
}
