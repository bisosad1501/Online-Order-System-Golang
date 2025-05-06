package api

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/online-order-system/user-service/interfaces"
	"github.com/online-order-system/user-service/models"
)

// Handlers contains all API handlers
type Handlers struct {
	service interfaces.UserService
}

// NewHandlers creates new API handlers
func NewHandlers(service interfaces.UserService) *Handlers {
	return &Handlers{
		service: service,
	}
}

// HealthCheck handles health check requests
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

// CreateUser handles user creation requests
func (h *Handlers) CreateUser(c *gin.Context) {
	var req models.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[POST] /users - Bad request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[POST] /users - Creating user with email: %s", req.Email)
	user, err := h.service.CreateUser(req)
	if err != nil {
		log.Printf("[POST] /users - Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[POST] /users - User created successfully with ID: %s", user.ID)
	c.JSON(http.StatusCreated, user)
}

// GetUserByID handles user retrieval by ID requests
func (h *Handlers) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsers handles user list retrieval requests
func (h *Handlers) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser handles user update requests
func (h *Handlers) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateUser(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion requests
func (h *Handlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// VerifyUser handles user verification requests
func (h *Handlers) VerifyUser(c *gin.Context) {
	var req models.VerifyCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.VerifyUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUserOrders handles user orders retrieval requests
func (h *Handlers) GetUserOrders(c *gin.Context) {
	id := c.Param("id")
	orders, err := h.service.GetUserOrders(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// Login handles user login requests
func (h *Handlers) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.service.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token := generateToken(user.ID, user.Email)

	// Return response
	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  &user,
	})
}

// ValidateToken handles token validation requests
func (h *Handlers) ValidateToken(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
		return
	}

	// Check if the header has the Bearer prefix
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	// Extract the token
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate token
	claims, err := validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set user ID and role headers
	c.Header("X-User-ID", claims.UserID)
	c.Header("X-User-Role", claims.Role)

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

// generateToken creates a simple token for authentication
func generateToken(userID, email string) string {
	// In a real app, this would be a JWT with proper signing
	// For demo purposes, we'll just create a base64 encoded string
	tokenData := userID + ":" + email + ":" + time.Now().Format(time.RFC3339)
	return base64.StdEncoding.EncodeToString([]byte(tokenData))
}

// validateToken validates a token and returns the claims
func validateToken(token string) (*models.TokenClaims, error) {
	// In a real app, this would verify the JWT signature
	// For demo purposes, we'll just decode the base64 string
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(data), ":")
	if len(parts) < 2 {
		return nil, err
	}

	return &models.TokenClaims{
		UserID: parts[0],
		Email:  parts[1],
		Role:   "user", // Default role
	}, nil
}
