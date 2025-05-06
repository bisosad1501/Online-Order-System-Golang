package api

import (
"log"
"net/http"

"github.com/gin-gonic/gin"
"github.com/online-order-system/order-service/interfaces"
"github.com/online-order-system/order-service/models"
)

// Handler handles HTTP requests
type Handler struct {
service interfaces.OrderService
}

// NewHandler creates a new handler
func NewHandler(service interfaces.OrderService) *Handler {
return &Handler{
service: service,
}
}

// CreateOrder handles the creation of a new order
// @Summary Create a new order
// @Description Create a new order with the provided details
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order details"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
var req models.CreateOrderRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

order, err := h.service.CreateOrder(req)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusCreated, order)
}

// GetOrderByID handles retrieving an order by ID
// @Summary Get order by ID
// @Description Get a specific order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 404 {object} map[string]string "Not Found"
// @Router /orders/{id} [get]
func (h *Handler) GetOrderByID(c *gin.Context) {
id := c.Param("id")

order, err := h.service.GetOrderByID(id)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
return
}

c.JSON(http.StatusOK, order)
}

// GetOrders handles retrieving all orders
// @Summary Get all orders
// @Description Get all orders in the system
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /orders [get]
func (h *Handler) GetOrders(c *gin.Context) {
orders, err := h.service.GetOrders()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
return
}

c.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus handles updating the status of an order
// @Summary Update order status
// @Description Update the status of an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body models.UpdateOrderStatusRequest true "Status details"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /orders/{id}/status [put]
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
id := c.Param("id")

var req models.UpdateOrderStatusRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

err := h.service.UpdateOrderStatus(id, req.Status)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

order, err := h.service.GetOrderByID(id)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, order)
}

// RetryPayment handles retrying payment for a failed order
// @Summary Retry payment for a failed order
// @Description Retry payment for an order that failed due to payment issues
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param payment body models.RetryPaymentRequest true "Payment details"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /orders/{id}/retry-payment [post]
func (h *Handler) RetryPayment(c *gin.Context) {
id := c.Param("id")
log.Printf("RetryPayment handler called for order ID: %s", id)

var req models.RetryPaymentRequest
if err := c.ShouldBindJSON(&req); err != nil {
log.Printf("Error binding JSON: %v", err)
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

log.Printf("Retrying payment for order %s with payment method: %s", id, req.PaymentMethod)
order, err := h.service.RetryPayment(id, req)
if err != nil {
log.Printf("Error retrying payment: %v", err)
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

log.Printf("Payment retried successfully for order %s", id)
c.JSON(http.StatusOK, order)
}
