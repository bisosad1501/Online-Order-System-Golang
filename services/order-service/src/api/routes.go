package api

import (
	"github.com/gin-gonic/gin"
	"github.com/online-order-system/order-service/interfaces"
	"log"
)

// SetupRouter sets up the router with all the necessary routes and middleware
func SetupRouter(service interfaces.OrderService) *gin.Engine {
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

	// Order routes
	orders := router.Group("/orders")
	{
		// Create a new order
		orders.POST("", handler.CreateOrder)

		// Get all orders
		orders.GET("", handler.GetOrders)

		// Get a specific order by ID
		orders.GET("/:id", handler.GetOrderByID)

		// Update the status of an order
		orders.PUT("/:id/status", handler.UpdateOrderStatus)

		// Retry payment for a failed order
		orders.POST("/:id/retry-payment", handler.RetryPayment)
	}

	log.Printf("Route registered: GET /orders")
	log.Printf("Route registered: GET /orders/:id")
	log.Printf("Route registered: GET /health")
	log.Printf("Route registered: POST /orders")
	log.Printf("Route registered: POST /orders/:id/retry-payment")
	log.Printf("Route registered: PUT /orders/:id/status")

	// Print all registered routes for debugging
	for _, route := range router.Routes() {
		log.Printf("Route registered: %s %s", route.Method, route.Path)
	}

	return router
}
