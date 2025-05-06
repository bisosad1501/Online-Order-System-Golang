package api

import (
"net/http"

"github.com/gin-gonic/gin"
"github.com/online-order-system/notification-service/interfaces"
"github.com/online-order-system/notification-service/models"
)

// Handler handles HTTP requests
type Handler struct {
service interfaces.NotificationService
}

// NewHandler creates a new handler
func NewHandler(service interfaces.NotificationService) *Handler {
return &Handler{
service: service,
}
}

// CreateNotification handles creating a new notification
func (h *Handler) CreateNotification(c *gin.Context) {
var req models.CreateNotificationRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

notification, err := h.service.CreateNotification(req)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
return
}

c.JSON(http.StatusCreated, notification)
}

// GetNotificationByID handles retrieving a notification by ID
func (h *Handler) GetNotificationByID(c *gin.Context) {
id := c.Param("id")

notification, err := h.service.GetNotificationByID(id)
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
return
}

c.JSON(http.StatusOK, notification)
}

// GetNotificationsByUserID handles retrieving notifications by user ID
func (h *Handler) GetNotificationsByUserID(c *gin.Context) {
userID := c.Param("user_id")

notifications, err := h.service.GetNotificationsByCustomerID(userID)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
return
}

c.JSON(http.StatusOK, notifications)
}

// GetNotifications handles retrieving all notifications
func (h *Handler) GetNotifications(c *gin.Context) {
notifications, err := h.service.GetNotifications()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
return
}

c.JSON(http.StatusOK, notifications)
}

// UpdateNotificationStatus handles updating a notification's status
func (h *Handler) UpdateNotificationStatus(c *gin.Context) {
id := c.Param("id")

var req models.UpdateNotificationStatusRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

notification, err := h.service.UpdateNotificationStatus(id, req.Status)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification status"})
return
}

c.JSON(http.StatusOK, notification)
}
