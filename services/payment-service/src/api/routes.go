package api

import (
"github.com/gin-gonic/gin"
"github.com/online-order-system/payment-service/interfaces"
)

// SetupRouter sets up the router with all the necessary routes and middleware
func SetupRouter(service interfaces.PaymentService) *gin.Engine {
// Create router
router := gin.Default()

// Add middleware
router.Use(Logger())
router.Use(CORS())

// Create handler
handler := NewHandler(service)

// Health check endpoint for monitoring and readiness checks
router.GET("/health", func(c *gin.Context) {
c.JSON(200, gin.H{
"status": "UP",
})
})

// Swagger documentation endpoints
router.StaticFile("/swagger", "./docs/swagger.html")
router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")

// Payment routes
payments := router.Group("/payments")
{
// Create a new payment
payments.POST("", handler.CreatePayment)

// Get all payments
payments.GET("", handler.GetPayments)

// Get a specific payment by ID
payments.GET("/:id", handler.GetPaymentByID)

// Get a payment by order ID
payments.GET("/order/:order_id", handler.GetPaymentByOrderID)

// Update a payment's status
payments.PUT("/:id/status", handler.UpdatePaymentStatus)

// Confirm a payment
payments.POST("/:id/confirm", handler.ConfirmPayment)

// Test successful payment with Stripe
payments.POST("/:id/test-success", handler.TestSuccessfulPayment)
}

// Stripe webhook route
// This endpoint receives webhook events from Stripe
router.POST("/webhooks/stripe", handler.HandleStripeWebhook)

return router
}
