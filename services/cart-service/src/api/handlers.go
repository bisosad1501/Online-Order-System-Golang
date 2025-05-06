package api

import (
"net/http"

"github.com/gin-gonic/gin"
"github.com/online-order-system/cart-service/interfaces"
"github.com/online-order-system/cart-service/models"
)

// Handlers handles HTTP requests
type Handlers struct {
service interfaces.CartService
}

// NewHandlers creates a new handlers
func NewHandlers(service interfaces.CartService) *Handlers {
return &Handlers{
service: service,
}
}

// CreateCart handles cart creation requests
func (h *Handlers) CreateCart(c *gin.Context) {
var req models.CreateCartRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

cart, err := h.service.CreateCart(req)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusCreated, cart)
}

// GetCart handles cart retrieval requests
func (h *Handlers) GetCart(c *gin.Context) {
id := c.Param("id")
cart, err := h.service.GetCartByID(id)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, cart)
}

// GetCartByUserID handles cart retrieval by user ID requests
func (h *Handlers) GetCartByUserID(c *gin.Context) {
userID := c.Param("user_id")
cart, err := h.service.GetCartByUserID(userID)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, cart)
}

// AddCartItem handles adding an item to a cart
func (h *Handlers) AddCartItem(c *gin.Context) {
id := c.Param("id")
var req models.AddCartItemRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

cart, err := h.service.AddCartItem(id, req)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, cart)
}

// UpdateCartItem handles updating an item in a cart
func (h *Handlers) UpdateCartItem(c *gin.Context) {
id := c.Param("id")
itemID := c.Param("item_id")
var req models.UpdateCartItemRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

cart, err := h.service.UpdateCartItem(id, itemID, req)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, cart)
}

// RemoveCartItem handles removing an item from a cart
func (h *Handlers) RemoveCartItem(c *gin.Context) {
id := c.Param("id")
itemID := c.Param("item_id")

cart, err := h.service.RemoveCartItem(id, itemID)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, cart)
}

// DeleteCart handles cart deletion requests
func (h *Handlers) DeleteCart(c *gin.Context) {
id := c.Param("id")

err := h.service.DeleteCart(id)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{"message": "Cart deleted successfully"})
}
