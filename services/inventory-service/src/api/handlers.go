package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/online-order-system/inventory-service/interfaces"
	"github.com/online-order-system/inventory-service/models"
)

// Handler handles HTTP requests
type Handler struct {
	service interfaces.InventoryService
}

// NewHandler creates a new handler
func NewHandler(service interfaces.InventoryService) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateProduct handles the creation of a new product
func (h *Handler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProductByID handles retrieving a product by ID
func (h *Handler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	// Validate product ID
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID cannot be empty"})
		return
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		// Check if error contains "not found" message
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// Log the error
		log.Printf("Error getting product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts handles retrieving all products
func (h *Handler) GetProducts(c *gin.Context) {
	products, err := h.service.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProduct handles updating a product
func (h *Handler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.UpdateProduct(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct handles deleting a product
func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}



// CheckInventory handles checking inventory
func (h *Handler) CheckInventory(c *gin.Context) {
	var req models.InventoryCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no items to check"})
		return
	}

	// Check for invalid product IDs
	for _, item := range req.Items {
		if item.ProductID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product ID cannot be empty"})
			return
		}
		if item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be greater than 0"})
			return
		}
	}

	response, err := h.service.CheckInventory(req)
	if err != nil {
		// Check if error contains "not found" message
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// Log the error
		log.Printf("Error checking inventory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProductRecommendations handles retrieving product recommendations by product ID
func (h *Handler) GetProductRecommendations(c *gin.Context) {
	productID := c.Param("id")
	limit := 10 // Default limit

	// Parse limit from query parameter if provided
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	recommendations, err := h.service.GetProductRecommendations(productID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": recommendations})
}

// GetCategoryRecommendations handles retrieving product recommendations by category ID
func (h *Handler) GetCategoryRecommendations(c *gin.Context) {
	categoryID := c.Param("id")
	limit := 10 // Default limit

	// Parse limit from query parameter if provided
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	recommendations, err := h.service.GetCategoryRecommendations(categoryID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": recommendations})
}

// GetSimilarProducts handles retrieving similar products
func (h *Handler) GetSimilarProducts(c *gin.Context) {
	var req models.RecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recommendations, err := h.service.GetSimilarProducts(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": recommendations})
}
