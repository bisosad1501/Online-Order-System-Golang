package api

import (
"github.com/gin-gonic/gin"
"github.com/online-order-system/notification-service/interfaces"
)

// SetupRouter sets up the router with all the necessary routes and middleware
func SetupRouter(service interfaces.NotificationService) *gin.Engine {
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

// Notification routes
notifications := router.Group("/notifications")
{
// Create a new notification
notifications.POST("", handler.CreateNotification)

// Get all notifications
notifications.GET("", handler.GetNotifications)

// Get a specific notification by ID
notifications.GET("/:id", handler.GetNotificationByID)

// Get notifications by user ID
notifications.GET("/user/:user_id", handler.GetNotificationsByUserID)

// Update a notification's status
notifications.PUT("/:id/status", handler.UpdateNotificationStatus)
}

return router
}
