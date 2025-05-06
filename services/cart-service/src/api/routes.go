package api

import (
"github.com/gin-gonic/gin"
"github.com/online-order-system/cart-service/interfaces"
)

// SetupRouter sets up the router with all the necessary routes and middleware
func SetupRouter(service interfaces.CartService) *gin.Engine {
// Create router
router := gin.Default()

// Add middleware
SetupMiddleware(router)

// Create handlers
handlers := NewHandlers(service)

// Health check endpoint for monitoring and readiness checks
router.GET("/health", func(c *gin.Context) {
c.JSON(200, gin.H{
"status": "UP",
})
})

// Swagger documentation endpoints
router.StaticFile("/swagger", "./docs/swagger.html")
router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")

// Cart routes
carts := router.Group("/carts")
{
carts.POST("", handlers.CreateCart)
carts.GET("/:id", handlers.GetCart)
carts.GET("/user/:user_id", handlers.GetCartByUserID)
carts.POST("/:id/items", handlers.AddCartItem)
carts.PUT("/:id/items/:item_id", handlers.UpdateCartItem)
carts.DELETE("/:id/items/:item_id", handlers.RemoveCartItem)
carts.DELETE("/:id", handlers.DeleteCart)
}

return router
}
