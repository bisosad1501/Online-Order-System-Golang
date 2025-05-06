package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/online-order-system/inventory-service/cache"
	"github.com/online-order-system/inventory-service/config"
	"github.com/online-order-system/inventory-service/db"
	"github.com/online-order-system/inventory-service/interfaces"
	"github.com/online-order-system/inventory-service/models"
)

// InventoryService handles business logic for inventory
type InventoryService struct {
	config     *config.Config
	repository *db.InventoryRepository
	producer   interfaces.InventoryProducer
	cache      *cache.RedisCache
}

// Ensure InventoryService implements InventoryService interface
var _ interfaces.InventoryService = (*InventoryService)(nil)

// NewInventoryService creates a new inventory service
func NewInventoryService(cfg *config.Config, repo *db.InventoryRepository, producer interfaces.InventoryProducer, redisCache *cache.RedisCache) *InventoryService {
	return &InventoryService{
		config:     cfg,
		repository: repo,
		producer:   producer,
		cache:      redisCache,
	}
}

// CreateProduct creates a new product with initial inventory
func (s *InventoryService) CreateProduct(req models.CreateProductRequest) (models.Product, error) {
	// Create product
	now := time.Now()
	product := models.Product{
		ID:          db.GenerateID(),
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
		Quantity:    req.Quantity, // Set initial quantity
		Tags:        req.Tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save product to database with initial inventory
	err := s.repository.CreateProduct(product, product.Quantity)
	if err != nil {
		return models.Product{}, err
	}

	// Save product tags if provided
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			productTag := models.ProductTag{
				ID:        db.GenerateID(),
				ProductID: product.ID,
				Tag:       tag,
				CreatedAt: now,
			}
			err := s.repository.CreateProductTag(productTag)
			if err != nil {
				// Log error but continue
				log.Printf("Failed to create product tag: %v", err)
			}
		}
	}

	// Cache the new product
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%s", product.ID)
	err = s.cache.Set(ctx, cacheKey, product)
	if err != nil {
		log.Printf("Failed to cache new product: %v", err)
	}

	// Invalidate products cache to ensure new product is included in list
	err = s.cache.Delete(ctx, "products:all")
	if err != nil {
		log.Printf("Failed to invalidate products cache: %v", err)
	} else {
		log.Printf("Successfully invalidated products cache after creating new product: %s", product.ID)
	}

	// Invalidate category cache
	categoryKey := fmt.Sprintf("category:%s", product.CategoryID)
	err = s.cache.Delete(ctx, categoryKey)
	if err != nil {
		log.Printf("Failed to invalidate category cache: %v", err)
	}

	// Invalidate inventory cache
	inventoryKey := fmt.Sprintf("inventory:%s", product.ID)
	err = s.cache.Delete(ctx, inventoryKey)
	if err != nil {
		log.Printf("Failed to invalidate inventory cache: %v", err)
	}

	// Publish inventory event if initial quantity > 0
	if product.Quantity > 0 {
		err = s.producer.PublishInventoryUpdated(product.ID, product.Quantity)
		if err != nil {
			// Log error but continue
			log.Printf("Failed to publish inventory updated event: %v", err)
		}
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (s *InventoryService) GetProductByID(id string) (models.Product, error) {
	// Try to get from cache first
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%s", id)
	var product models.Product

	err := s.cache.Get(ctx, cacheKey, &product)
	if err == nil {
		log.Printf("Cache hit for product: %s", id)
		return product, nil
	}

	// If not in cache, get from database
	product, err = s.repository.GetProductByID(id)
	if err != nil {
		return models.Product{}, err
	}

	// Store in cache for future requests
	err = s.cache.Set(ctx, cacheKey, product)
	if err != nil {
		log.Printf("Failed to cache product: %v", err)
	}

	return product, nil
}

// GetProducts retrieves all products
func (s *InventoryService) GetProducts() ([]models.Product, error) {
	// Try to get from cache first
	ctx := context.Background()
	cacheKey := "products:all"
	var products []models.Product

	err := s.cache.Get(ctx, cacheKey, &products)
	if err == nil {
		log.Printf("Cache hit for all products")
		return products, nil
	}

	// If not in cache, get from database
	products, err = s.repository.GetProducts()
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	err = s.cache.Set(ctx, cacheKey, products)
	if err != nil {
		log.Printf("Failed to cache products: %v", err)
	}

	return products, nil
}

// UpdateProduct updates a product and its inventory
func (s *InventoryService) UpdateProduct(id string, req models.UpdateProductRequest) (models.Product, error) {
	// Get existing product
	product, err := s.repository.GetProductByID(id)
	if err != nil {
		return models.Product{}, err
	}

	// Update product fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.CategoryID != "" {
		product.CategoryID = req.CategoryID
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if len(req.Tags) > 0 {
		product.Tags = req.Tags
	}

	// Update quantity if provided
	if req.Quantity != nil {
		product.Quantity = *req.Quantity

		// Publish inventory updated event
		err = s.producer.PublishInventoryUpdated(id, *req.Quantity)
		if err != nil {
			// Log error but continue
			log.Printf("Failed to publish inventory updated event: %v", err)
		}
	}

	product.UpdatedAt = time.Now()

	// Save product and inventory to database in a single transaction
	err = s.repository.UpdateProduct(id, product)
	if err != nil {
		return models.Product{}, err
	}

	// Update cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%s", id)
	err = s.cache.Set(ctx, cacheKey, product)
	if err != nil {
		log.Printf("Failed to update product in cache: %v", err)
	}

	// Invalidate products cache
	err = s.cache.Delete(ctx, "products:all")
	if err != nil {
		log.Printf("Failed to invalidate products cache: %v", err)
	}

	// Invalidate category cache
	categoryKey := fmt.Sprintf("category:%s", product.CategoryID)
	err = s.cache.Delete(ctx, categoryKey)
	if err != nil {
		log.Printf("Failed to invalidate category cache: %v", err)
	}

	// Invalidate inventory cache
	inventoryKey := fmt.Sprintf("inventory:%s", id)
	err = s.cache.Delete(ctx, inventoryKey)
	if err != nil {
		log.Printf("Failed to invalidate inventory cache: %v", err)
	}

	return product, nil
}

// DeleteProduct deletes a product
func (s *InventoryService) DeleteProduct(id string) error {
	// Get product to get category ID before deletion
	product, err := s.repository.GetProductByID(id)
	if err != nil {
		return err
	}

	// Delete from database
	err = s.repository.DeleteProduct(id)
	if err != nil {
		return err
	}

	// Delete from cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%s", id)
	err = s.cache.Delete(ctx, cacheKey)
	if err != nil {
		log.Printf("Failed to delete product from cache: %v", err)
	}

	// Invalidate products cache
	err = s.cache.Delete(ctx, "products:all")
	if err != nil {
		log.Printf("Failed to invalidate products cache: %v", err)
	}

	// Invalidate category cache
	categoryKey := fmt.Sprintf("category:%s", product.CategoryID)
	err = s.cache.Delete(ctx, categoryKey)
	if err != nil {
		log.Printf("Failed to invalidate category cache: %v", err)
	}

	return nil
}



// CheckInventory checks if all items are available in inventory
func (s *InventoryService) CheckInventory(req models.InventoryCheckRequest) (models.InventoryCheckResponse, error) {
	if len(req.Items) == 0 {
		return models.InventoryCheckResponse{}, errors.New("no items to check")
	}

	// Check inventory
	available, unavailableItems, err := s.repository.CheckInventory(req.Items)
	if err != nil {
		return models.InventoryCheckResponse{}, err
	}

	// Create response
	response := models.InventoryCheckResponse{
		Available: available,
	}

	if !available {
		response.UnavailableItems = unavailableItems
	}

	return response, nil
}

// GetProductRecommendations retrieves product recommendations based on product ID
func (s *InventoryService) GetProductRecommendations(productID string, limit int) ([]models.Product, error) {
	// Set default limit if not provided
	if limit <= 0 {
		limit = 10
	}

	// Try to get from cache first
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%s:recommendations:limit:%d", productID, limit)
	var recommendations []models.Product

	err := s.cache.Get(ctx, cacheKey, &recommendations)
	if err == nil {
		log.Printf("Cache hit for product recommendations: %s", productID)
		return recommendations, nil
	}

	// Get product to check if it exists
	product, err := s.repository.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	// Get similar products based on category
	similarProducts, err := s.repository.GetProductsByCategory(product.CategoryID, limit)
	if err != nil {
		return nil, err
	}

	// Filter out the original product
	recommendations = []models.Product{}
	for _, p := range similarProducts {
		if p.ID != productID {
			recommendations = append(recommendations, p)
		}
	}

	// If we don't have enough recommendations, get more products
	if len(recommendations) < limit {
		moreProducts, err := s.repository.GetProducts()
		if err != nil {
			// Store in cache for future requests
			err = s.cache.Set(ctx, cacheKey, recommendations)
			if err != nil {
				log.Printf("Failed to cache product recommendations: %v", err)
			}
			return recommendations, nil // Return what we have so far
		}

		// Add more products until we reach the limit
		for _, p := range moreProducts {
			if p.ID != productID && !containsProduct(recommendations, p.ID) {
				recommendations = append(recommendations, p)
				if len(recommendations) >= limit {
					break
				}
			}
		}
	}

	// Store in cache for future requests
	err = s.cache.Set(ctx, cacheKey, recommendations)
	if err != nil {
		log.Printf("Failed to cache product recommendations: %v", err)
	}

	return recommendations, nil
}

// GetCategoryRecommendations retrieves product recommendations based on category ID
func (s *InventoryService) GetCategoryRecommendations(categoryID string, limit int) ([]models.Product, error) {
	// Set default limit if not provided
	if limit <= 0 {
		limit = 10
	}

	// Try to get from cache first
	ctx := context.Background()
	cacheKey := fmt.Sprintf("category:%s:limit:%d", categoryID, limit)
	var products []models.Product

	err := s.cache.Get(ctx, cacheKey, &products)
	if err == nil {
		log.Printf("Cache hit for category recommendations: %s", categoryID)
		return products, nil
	}

	// If not in cache, get from database
	products, err = s.repository.GetProductsByCategory(categoryID, limit)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	err = s.cache.Set(ctx, cacheKey, products)
	if err != nil {
		log.Printf("Failed to cache category recommendations: %v", err)
	}

	return products, nil
}

// GetSimilarProducts retrieves similar products based on various criteria
func (s *InventoryService) GetSimilarProducts(req models.RecommendationRequest) ([]models.Product, error) {
	// Set default limit if not provided
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Try to get from cache first
	ctx := context.Background()
	reqJSON, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal recommendation request: %v", err)
	} else {
		cacheKey := fmt.Sprintf("similar:%s", string(reqJSON))
		var products []models.Product

		err := s.cache.Get(ctx, cacheKey, &products)
		if err == nil {
			log.Printf("Cache hit for similar products")
			return products, nil
		}
	}

	var products []models.Product

	// If product ID is provided, get recommendations based on product
	if req.ProductID != "" {
		products, err = s.GetProductRecommendations(req.ProductID, req.Limit)
		if err != nil {
			return nil, err
		}
	} else if req.CategoryID != "" {
		// If category ID is provided, get recommendations based on category
		products, err = s.GetCategoryRecommendations(req.CategoryID, req.Limit)
		if err != nil {
			return nil, err
		}
	} else {
		// Otherwise, get products based on other criteria
		products, err = s.repository.GetProducts()
		if err != nil {
			return nil, err
		}
	}

	// Filter by price range if provided
	if req.PriceMin > 0 || req.PriceMax > 0 {
		var filteredProducts []models.Product
		for _, p := range products {
			if (req.PriceMin <= 0 || p.Price >= req.PriceMin) &&
				(req.PriceMax <= 0 || p.Price <= req.PriceMax) {
				filteredProducts = append(filteredProducts, p)
			}
		}
		products = filteredProducts
	}

	// Filter by tags if provided
	if len(req.Tags) > 0 {
		var filteredProducts []models.Product
		for _, p := range products {
			if containsAnyTag(p.Tags, req.Tags) {
				filteredProducts = append(filteredProducts, p)
			}
		}
		products = filteredProducts
	}

	// Limit the number of products
	if len(products) > req.Limit {
		products = products[:req.Limit]
	}

	// Store in cache for future requests
	reqJSON, err = json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal recommendation request: %v", err)
	} else {
		cacheKey := fmt.Sprintf("similar:%s", string(reqJSON))
		err = s.cache.Set(ctx, cacheKey, products)
		if err != nil {
			log.Printf("Failed to cache similar products: %v", err)
		}
	}

	return products, nil
}

// Helper function to check if a product is already in the recommendations
func containsProduct(products []models.Product, id string) bool {
	for _, p := range products {
		if p.ID == id {
			return true
		}
	}
	return false
}

// Helper function to check if a product has any of the specified tags
func containsAnyTag(productTags []string, tags []string) bool {
	for _, pt := range productTags {
		for _, t := range tags {
			if pt == t {
				return true
			}
		}
	}
	return false
}
