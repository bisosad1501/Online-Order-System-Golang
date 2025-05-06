package api

import (
"github.com/gin-gonic/gin"
"github.com/online-order-system/inventory-service/interfaces"
)

// SetupRouter sets up the router with all the necessary routes and middleware
func SetupRouter(service interfaces.InventoryService) *gin.Engine {
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

// Product routes
products := router.Group("/products")
{
// Create a new product
products.POST("", handler.CreateProduct)

// Get all products
products.GET("", handler.GetProducts)

// Get a specific product by ID
products.GET("/:id", handler.GetProductByID)

// Update a product (full update)
products.PUT("/:id", handler.UpdateProduct)

// Update a product (partial update)
products.PATCH("/:id", handler.UpdateProduct)

// Delete a product
products.DELETE("/:id", handler.DeleteProduct)
}

// Inventory routes
inventory := router.Group("/inventory")
{
// Check inventory
inventory.POST("/check", handler.CheckInventory)
}

// Recommendation routes
recommendations := router.Group("/recommendations")
{
// Get recommendations by product ID
recommendations.GET("/product/:id", handler.GetProductRecommendations)

// Get recommendations by category ID
recommendations.GET("/category/:id", handler.GetCategoryRecommendations)

// Get similar products
recommendations.POST("/similar", handler.GetSimilarProducts)
}

return router
}
