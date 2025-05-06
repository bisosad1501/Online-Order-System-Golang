package api

import (
"github.com/gin-gonic/gin"
"github.com/online-order-system/shipping-service/interfaces"
)

// SetupRouter sets up the router with all the necessary routes and middleware
func SetupRouter(service interfaces.ShippingService) *gin.Engine {
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

// Shipment routes
shipments := router.Group("/shipments")
{
// Create a new shipment
shipments.POST("", handler.CreateShipment)

// Get all shipments
shipments.GET("", handler.GetShipments)

// Get a specific shipment by ID
shipments.GET("/:id", handler.GetShipmentByID)

// Get a shipment by order ID
shipments.GET("/order/:order_id", handler.GetShipmentByOrderID)

// Update a shipment's status
shipments.PUT("/:id/status", handler.UpdateShipmentStatus)

// Update a shipment's tracking number
shipments.PUT("/:id/tracking", handler.UpdateTrackingNumber)
}

return router
}
