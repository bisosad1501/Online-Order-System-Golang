package service

import (
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/online-order-system/user-service/config"
	"github.com/online-order-system/user-service/db"
	"github.com/online-order-system/user-service/interfaces"
	"github.com/online-order-system/user-service/models"
)

// UserService handles business logic for users
type UserService struct {
	config     *config.Config
	repository *db.UserRepository
	producer   interfaces.UserProducer
}

// Ensure UserService implements UserService interface
var _ interfaces.UserService = (*UserService)(nil)

// NewUserService creates a new user service
func NewUserService(cfg *config.Config, repo *db.UserRepository, producer interfaces.UserProducer) *UserService {
	return &UserService{
		config:     cfg,
		repository: repo,
		producer:   producer,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req models.CreateCustomerRequest) (models.Customer, error) {
	// Check if user with email already exists
	_, err := s.GetUserByEmail(req.Email)
	if err == nil {
		return models.Customer{}, errors.New("user with this email already exists")
	}

	// Create user
	now := time.Now()
	user := models.Customer{
		ID:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Address:   req.Address,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save user to database
	err = s.repository.CreateUser(user)
	if err != nil {
		return models.Customer{}, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id string) (models.Customer, error) {
	return s.repository.GetUserByID(id)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (models.Customer, error) {
	return s.repository.GetUserByEmail(email)
}

// GetUsers retrieves all users
func (s *UserService) GetUsers() ([]models.Customer, error) {
	return s.repository.GetUsers()
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id string, req models.UpdateCustomerRequest) (models.Customer, error) {
	// Get existing user
	user, err := s.GetUserByID(id)
	if err != nil {
		return models.Customer{}, err
	}

	// Update fields if provided
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	user.UpdatedAt = time.Now()

	// Save updated user to database
	err = s.repository.UpdateUser(user)
	if err != nil {
		return models.Customer{}, err
	}

	// Publish user updated event
	err = s.producer.PublishUserUpdated(user)
	if err != nil {
		log.Printf("Failed to publish user updated event: %v", err)
		// Continue anyway for demo purposes
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	// Check if user exists
	_, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// Delete user from database
	return s.repository.DeleteUser(id)
}

// VerifyUser verifies a user's information
func (s *UserService) VerifyUser(req models.VerifyCustomerRequest) (models.VerifyCustomerResponse, error) {
	var user models.Customer
	var err error

	// If ID is provided, get user by ID
	if req.ID != "" {
		user, err = s.GetUserByID(req.ID)
		if err != nil {
			return models.VerifyCustomerResponse{
				Verified: false,
				Message:  "User not found",
			}, err
		}
	} else {
		// Otherwise, get user by email
		user, err = s.GetUserByEmail(req.Email)
		if err != nil {
			return models.VerifyCustomerResponse{
				Verified: false,
				Message:  "User not found",
			}, err
		}
	}

	// If address is provided, verify it
	if req.Address != "" && user.Address != req.Address {
		return models.VerifyCustomerResponse{
			Verified: false,
			Message:  "Address does not match",
		}, nil
	}

	// Generate a simple token (in a real app, this would be a JWT)
	token := generateToken(user.ID, user.Email)

	// Publish user verified event
	err = s.producer.PublishUserVerified(user)
	if err != nil {
		log.Printf("Failed to publish user verified event: %v", err)
		// Continue anyway for demo purposes
	}

	return models.VerifyCustomerResponse{
		Verified: true,
		Message:  "User verified successfully",
		Token:    token,
		User:     &user,
	}, nil
}

// generateToken creates a simple token for authentication
func generateToken(userID, email string) string {
	// In a real app, this would be a JWT with proper signing
	// For demo purposes, we'll just create a base64 encoded string
	tokenData := userID + ":" + email + ":" + time.Now().Format(time.RFC3339)
	return base64.StdEncoding.EncodeToString([]byte(tokenData))
}

// GetUserOrders retrieves all orders for a user
func (s *UserService) GetUserOrders(userID string) ([]models.CustomerOrder, error) {
	// Check if user exists
	_, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.repository.GetUserOrders(userID)
}

// AddUserOrder adds an order to a user
func (s *UserService) AddUserOrder(userID string, orderID string, orderStatus string) error {
	// Check if user exists
	_, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	return s.repository.AddUserOrder(userID, orderID, orderStatus)
}

// UpdateUserOrderStatus updates the status of a user's order
func (s *UserService) UpdateUserOrderStatus(userID string, orderID string, orderStatus string) error {
	// Check if user exists
	_, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	return s.repository.UpdateUserOrderStatus(userID, orderID, orderStatus)
}
