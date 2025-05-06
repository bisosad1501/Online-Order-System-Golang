package api

import (
	"github.com/gin-gonic/gin"
	"github.com/online-order-system/user-service/interfaces"
)

// SetupRouter sets up the router
func SetupRouter(service interfaces.UserService) *gin.Engine {
	router := gin.Default()

	// Set up middleware
	SetupMiddleware(router)

	// Create handlers
	handlers := NewHandlers(service)

	// Serve Swagger UI
	router.StaticFile("/swagger", "./docs/swagger.html")
	router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// Users
	users := router.Group("/users")
	{
		users.POST("", handlers.CreateUser)
		users.GET("", handlers.GetUsers)
		users.GET("/:id", handlers.GetUserByID)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
		users.GET("/:id/orders", handlers.GetUserOrders)
		users.POST("/verify", handlers.VerifyUser)
	}

	// Auth
	auth := router.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.GET("/validate", handlers.ValidateToken)
	}

	return router
}
