package interfaces

import (
	"github.com/online-order-system/inventory-service/models"
)

// InventoryService defines the interface for inventory service
type InventoryService interface {
	// Product methods with integrated inventory
	CreateProduct(req models.CreateProductRequest) (models.Product, error)
	GetProductByID(id string) (models.Product, error)
	GetProducts() ([]models.Product, error)
	UpdateProduct(id string, req models.UpdateProductRequest) (models.Product, error)
	DeleteProduct(id string) error

	// Inventory check method
	CheckInventory(req models.InventoryCheckRequest) (models.InventoryCheckResponse, error)

	// Recommendation methods
	GetProductRecommendations(productID string, limit int) ([]models.Product, error)
	GetCategoryRecommendations(categoryID string, limit int) ([]models.Product, error)
	GetSimilarProducts(req models.RecommendationRequest) ([]models.Product, error)
}

// InventoryProducer defines the interface for inventory producer
type InventoryProducer interface {
	PublishInventoryUpdated(productID string, quantity int) error
	Close() error
}
