package service

import (
"log"
"time"

"github.com/google/uuid"
"github.com/online-order-system/cart-service/config"
"github.com/online-order-system/cart-service/db"
"github.com/online-order-system/cart-service/interfaces"
"github.com/online-order-system/cart-service/models"
)

// CartService handles business logic for carts
type CartService struct {
config     *config.Config
repository *db.CartRepository
producer   interfaces.CartProducer
}

// Ensure CartService implements CartService interface
var _ interfaces.CartService = (*CartService)(nil)

// NewCartService creates a new cart service
func NewCartService(cfg *config.Config, repo *db.CartRepository, producer interfaces.CartProducer) *CartService {
return &CartService{
config:     cfg,
repository: repo,
producer:   producer,
}
}

// CreateCart creates a new cart
func (s *CartService) CreateCart(req models.CreateCartRequest) (models.Cart, error) {
// Check if cart already exists for user
existingCart, err := s.GetCartByUserID(req.CustomerID)
if err == nil {
// Cart already exists
return existingCart, nil
}

// Create cart
now := time.Now()
cart := models.Cart{
ID:         uuid.New().String(),
CustomerID: req.CustomerID,
Items:      []models.CartItem{},
CreatedAt:  now,
UpdatedAt:  now,
}

// Save cart to database
err = s.repository.CreateCart(cart)
if err != nil {
return models.Cart{}, err
}

return cart, nil
}

// GetCartByID retrieves a cart by ID
func (s *CartService) GetCartByID(id string) (models.Cart, error) {
return s.repository.GetCartByID(id)
}

// GetCartByUserID retrieves a cart by user ID
func (s *CartService) GetCartByUserID(userID string) (models.Cart, error) {
return s.repository.GetCartByUserID(userID)
}

// AddCartItem adds an item to a cart
func (s *CartService) AddCartItem(cartID string, req models.AddCartItemRequest) (models.Cart, error) {
// Check if cart exists
_, err := s.GetCartByID(cartID)
if err != nil {
return models.Cart{}, err
}

// Create cart item
now := time.Now()
item := models.CartItem{
ID:        uuid.New().String(),
CartID:    cartID,
ProductID: req.ProductID,
Quantity:  req.Quantity,
Price:     req.Price,
CreatedAt: now,
UpdatedAt: now,
}

// Add item to cart
err = s.repository.AddCartItem(item)
if err != nil {
return models.Cart{}, err
}

// Update cart updated_at
err = s.repository.UpdateCartUpdatedAt(cartID)
if err != nil {
return models.Cart{}, err
}

// Get updated cart
cart, err := s.GetCartByID(cartID)
if err != nil {
return models.Cart{}, err
}

// Publish cart updated event
err = s.producer.PublishCartUpdated(cart)
if err != nil {
log.Printf("Failed to publish cart updated event: %v", err)
// Continue anyway
}

return cart, nil
}

// UpdateCartItem updates an item in a cart
func (s *CartService) UpdateCartItem(cartID string, itemID string, req models.UpdateCartItemRequest) (models.Cart, error) {
// Check if cart exists
_, err := s.GetCartByID(cartID)
if err != nil {
return models.Cart{}, err
}

// Update cart item
err = s.repository.UpdateCartItem(cartID, itemID, req.Quantity)
if err != nil {
return models.Cart{}, err
}

// Update cart updated_at
err = s.repository.UpdateCartUpdatedAt(cartID)
if err != nil {
return models.Cart{}, err
}

// Get updated cart
cart, err := s.GetCartByID(cartID)
if err != nil {
return models.Cart{}, err
}

// Publish cart updated event
err = s.producer.PublishCartUpdated(cart)
if err != nil {
log.Printf("Failed to publish cart updated event: %v", err)
// Continue anyway
}

return cart, nil
}

// RemoveCartItem removes an item from a cart
func (s *CartService) RemoveCartItem(cartID string, itemID string) (models.Cart, error) {
// Check if cart exists
_, err := s.GetCartByID(cartID)
if err != nil {
return models.Cart{}, err
}

// Remove cart item
err = s.repository.RemoveCartItem(cartID, itemID)
if err != nil {
return models.Cart{}, err
}

// Update cart updated_at
err = s.repository.UpdateCartUpdatedAt(cartID)
if err != nil {
return models.Cart{}, err
}

// Get updated cart
cart, err := s.GetCartByID(cartID)
if err != nil {
return models.Cart{}, err
}

// Publish cart updated event
err = s.producer.PublishCartUpdated(cart)
if err != nil {
log.Printf("Failed to publish cart updated event: %v", err)
// Continue anyway
}

return cart, nil
}

// DeleteCart deletes a cart
func (s *CartService) DeleteCart(id string) error {
// Check if cart exists
cart, err := s.GetCartByID(id)
if err != nil {
return err
}

// Delete cart
err = s.repository.DeleteCart(id)
if err != nil {
return err
}

// Publish cart cleared event
err = s.producer.PublishCartCleared(id, cart.CustomerID)
if err != nil {
log.Printf("Failed to publish cart cleared event: %v", err)
// Continue anyway
}

return nil
}

// DeleteCartByUserID deletes a cart by user ID
func (s *CartService) DeleteCartByUserID(userID string) error {
// Check if cart exists
_, err := s.GetCartByUserID(userID)
if err != nil {
// Cart doesn't exist, nothing to delete
return nil
}

// Delete cart
return s.repository.DeleteCartByUserID(userID)
}
