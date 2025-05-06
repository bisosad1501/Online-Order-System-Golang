package api

import (
"net/http"

"github.com/gin-gonic/gin"
"github.com/online-order-system/shipping-service/interfaces"
"github.com/online-order-system/shipping-service/models"
)

// Handler handles HTTP requests
type Handler struct {
service interfaces.ShippingService
}

// NewHandler creates a new handler
func NewHandler(service interfaces.ShippingService) *Handler {
return &Handler{
service: service,
}
}

// CreateShipment handles the creation of a new shipment
func (h *Handler) CreateShipment(c *gin.Context) {
var req models.CreateShipmentRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

shipment, err := h.service.CreateShipment(req)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusCreated, shipment)
}

// GetShipmentByID handles retrieving a shipment by ID
func (h *Handler) GetShipmentByID(c *gin.Context) {
id := c.Param("id")

shipment, err := h.service.GetShipmentByID(id)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "Shipment not found"})
return
}

c.JSON(http.StatusOK, shipment)
}

// GetShipmentByOrderID handles retrieving a shipment by order ID
func (h *Handler) GetShipmentByOrderID(c *gin.Context) {
orderID := c.Param("order_id")

shipment, err := h.service.GetShipmentByOrderID(orderID)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "Shipment not found"})
return
}

c.JSON(http.StatusOK, shipment)
}

// GetShipments handles retrieving all shipments
func (h *Handler) GetShipments(c *gin.Context) {
shipments, err := h.service.GetShipments()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get shipments"})
return
}

c.JSON(http.StatusOK, shipments)
}

// UpdateShipmentStatus handles updating a shipment's status
func (h *Handler) UpdateShipmentStatus(c *gin.Context) {
id := c.Param("id")

var req models.UpdateShipmentStatusRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

shipment, err := h.service.UpdateShipmentStatus(id, req.Status)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, shipment)
}

// UpdateTrackingNumber handles updating a shipment's tracking number
func (h *Handler) UpdateTrackingNumber(c *gin.Context) {
id := c.Param("id")

var req models.UpdateTrackingRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

shipment, err := h.service.UpdateTrackingNumber(id, req.TrackingNumber)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, shipment)
}