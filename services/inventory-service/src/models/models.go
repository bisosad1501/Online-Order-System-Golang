package models

import (
	"time"
)

// Product represents a product in the system with its inventory information
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryID  string    `json:"category_id"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`    // Integrated inventory quantity
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Inventory represents the inventory of a product
type Inventory struct {
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductTag represents a tag for a product
type ProductTag struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductSimilarity represents the similarity between two products
type ProductSimilarity struct {
	ID               string    `json:"id"`
	ProductID        string    `json:"product_id"`
	SimilarProductID string    `json:"similar_product_id"`
	SimilarityScore  float64   `json:"similarity_score"`
	CreatedAt        time.Time `json:"created_at"`
}

// CreateProductRequest represents a request to create a new product
type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	CategoryID  string   `json:"category_id" binding:"required"`
	Price       float64  `json:"price" binding:"required"`
	Tags        []string `json:"tags,omitempty"`
	Quantity    int      `json:"quantity"`  // Initial inventory quantity, defaults to 0 if not provided
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  string   `json:"category_id"`
	Price       float64  `json:"price"`
	Tags        []string `json:"tags,omitempty"`
	Quantity    *int     `json:"quantity,omitempty"`  // Optional inventory quantity update
}



// InventoryCheckRequest represents a request to check inventory
type InventoryCheckRequest struct {
	Items []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

// InventoryCheckResponse represents a response from inventory check
type InventoryCheckResponse struct {
	Available        bool `json:"available"`
	UnavailableItems []struct {
		ProductID   string `json:"product_id"`
		ProductName string `json:"product_name"`
		Requested   int    `json:"requested"`
		Available   int    `json:"available"`
	} `json:"unavailable_items,omitempty"`
}

// InventoryEvent represents an event related to inventory
type InventoryEvent struct {
	EventType string `json:"event_type"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Timestamp int64  `json:"timestamp"`
}

// RecommendationRequest represents a request for product recommendations
type RecommendationRequest struct {
	ProductID  string   `json:"product_id,omitempty"`
	CategoryID string   `json:"category_id,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	PriceMin   float64  `json:"price_min,omitempty"`
	PriceMax   float64  `json:"price_max,omitempty"`
	Limit      int      `json:"limit,omitempty"`
}

// RecommendationResponse represents a response with product recommendations
type RecommendationResponse struct {
	Products []Product `json:"products"`
}
