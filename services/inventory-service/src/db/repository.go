package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/online-order-system/inventory-service/models"
)

// InventoryRepository handles database operations for inventory
type InventoryRepository struct {
	db *Database
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *Database) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// CreateProduct creates a new product in the database with initial inventory
func (r *InventoryRepository) CreateProduct(product models.Product, initialQuantity int) error {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert product
	_, err = tx.Exec(
		"INSERT INTO products (id, name, description, category_id, price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		product.ID, product.Name, product.Description, product.CategoryID, product.Price, product.CreatedAt, product.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert inventory with specified initial quantity
	_, err = tx.Exec(
		"INSERT INTO inventory (product_id, quantity, updated_at) VALUES ($1, $2, $3)",
		product.ID, initialQuantity, product.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// CreateProductTag creates a new product tag in the database
func (r *InventoryRepository) CreateProductTag(tag models.ProductTag) error {
	_, err := r.db.Exec(
		"INSERT INTO product_tags (id, product_id, tag, created_at) VALUES ($1, $2, $3, $4)",
		tag.ID, tag.ProductID, tag.Tag, tag.CreatedAt,
	)
	return err
}

// GetProductByID retrieves a product by ID with its inventory information
func (r *InventoryRepository) GetProductByID(id string) (models.Product, error) {
	var product models.Product
	var createdAt, updatedAt time.Time

	// Get product with inventory information using JOIN
	err := r.db.QueryRow(
		`SELECT p.id, p.name, p.description, p.category_id, p.price, p.created_at, p.updated_at,
		COALESCE(i.quantity, 0) as quantity
		FROM products p
		LEFT JOIN inventory i ON p.id = i.product_id
		WHERE p.id = $1`,
		id,
	).Scan(&product.ID, &product.Name, &product.Description, &product.CategoryID, &product.Price,
		&createdAt, &updatedAt, &product.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return product, fmt.Errorf("product with ID %s not found", id)
		}
		return product, fmt.Errorf("error getting product: %w", err)
	}

	product.CreatedAt = createdAt
	product.UpdatedAt = updatedAt

	// Get product tags
	rows, err := r.db.Query("SELECT tag FROM product_tags WHERE product_id = $1", id)
	if err != nil {
		// Log error but continue
		// log.Printf("Failed to get product tags: %v", err)
	} else {
		defer rows.Close()
		var tags []string
		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err == nil {
				tags = append(tags, tag)
			}
		}
		product.Tags = tags
	}

	return product, nil
}

// GetProducts retrieves all products with their inventory information
func (r *InventoryRepository) GetProducts() ([]models.Product, error) {
	// Get all products with inventory information using JOIN
	rows, err := r.db.Query(
		`SELECT p.id, p.name, p.description, p.category_id, p.price, p.created_at, p.updated_at,
		COALESCE(i.quantity, 0) as quantity
		FROM products p
		LEFT JOIN inventory i ON p.id = i.product_id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var createdAt, updatedAt time.Time

		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.CategoryID, &product.Price,
			&createdAt, &updatedAt, &product.Quantity)
		if err != nil {
			return nil, err
		}

		product.CreatedAt = createdAt
		product.UpdatedAt = updatedAt

		// Get product tags
		tagRows, err := r.db.Query("SELECT tag FROM product_tags WHERE product_id = $1", product.ID)
		if err == nil {
			defer tagRows.Close()
			var tags []string
			for tagRows.Next() {
				var tag string
				if err := tagRows.Scan(&tag); err == nil {
					tags = append(tags, tag)
				}
			}
			product.Tags = tags
		}

		products = append(products, product)
	}

	return products, nil
}

// UpdateProduct updates a product and its inventory in the database
func (r *InventoryRepository) UpdateProduct(id string, product models.Product) error {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update product
	_, err = tx.Exec(
		"UPDATE products SET name = $1, description = $2, category_id = $3, price = $4, updated_at = $5 WHERE id = $6",
		product.Name, product.Description, product.CategoryID, product.Price, product.UpdatedAt, id,
	)
	if err != nil {
		return err
	}

	// Update inventory - use upsert to handle both insert and update cases
	_, err = tx.Exec(
		`INSERT INTO inventory (product_id, quantity, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (product_id)
		DO UPDATE SET quantity = $2, updated_at = $3`,
		id, product.Quantity, product.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Delete existing tags
	_, err = tx.Exec("DELETE FROM product_tags WHERE product_id = $1", id)
	if err != nil {
		return err
	}

	// Insert new tags
	if len(product.Tags) > 0 {
		for _, tag := range product.Tags {
			_, err = tx.Exec(
				"INSERT INTO product_tags (id, product_id, tag, created_at) VALUES ($1, $2, $3, $4)",
				GenerateID(), id, tag, time.Now(),
			)
			if err != nil {
				return err
			}
		}
	}

	// Commit transaction
	return tx.Commit()
}

// DeleteProduct deletes a product from the database
func (r *InventoryRepository) DeleteProduct(id string) error {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete inventory
	_, err = tx.Exec("DELETE FROM inventory WHERE product_id = $1", id)
	if err != nil {
		return err
	}

	// Delete product tags
	_, err = tx.Exec("DELETE FROM product_tags WHERE product_id = $1", id)
	if err != nil {
		return err
	}

	// Delete product
	_, err = tx.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}



// CheckInventory checks if all items are available in inventory
func (r *InventoryRepository) CheckInventory(items []struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}) (bool, []struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Requested   int    `json:"requested"`
	Available   int    `json:"available"`
}, error) {
	var unavailableItems []struct {
		ProductID   string `json:"product_id"`
		ProductName string `json:"product_name"`
		Requested   int    `json:"requested"`
		Available   int    `json:"available"`
	}

	for _, item := range items {
		// Get product with inventory information using JOIN
		var productName string
		var quantity int
		err := r.db.QueryRow(
			`SELECT p.name, COALESCE(i.quantity, 0) as quantity
			FROM products p
			LEFT JOIN inventory i ON p.id = i.product_id
			WHERE p.id = $1`,
			item.ProductID,
		).Scan(&productName, &quantity)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, nil, fmt.Errorf("product with ID %s not found", item.ProductID)
			}
			return false, nil, fmt.Errorf("error getting product: %w", err)
		}

		// Check if quantity is sufficient
		if quantity < item.Quantity {
			unavailableItems = append(unavailableItems, struct {
				ProductID   string `json:"product_id"`
				ProductName string `json:"product_name"`
				Requested   int    `json:"requested"`
				Available   int    `json:"available"`
			}{
				ProductID:   item.ProductID,
				ProductName: productName,
				Requested:   item.Quantity,
				Available:   quantity,
			})
		}
	}

	return len(unavailableItems) == 0, unavailableItems, nil
}

// GetProductsByCategory retrieves products by category ID with their inventory information
func (r *InventoryRepository) GetProductsByCategory(categoryID string, limit int) ([]models.Product, error) {
	// Get products by category with inventory information using JOIN
	rows, err := r.db.Query(
		`SELECT p.id, p.name, p.description, p.category_id, p.price, p.created_at, p.updated_at,
		COALESCE(i.quantity, 0) as quantity
		FROM products p
		LEFT JOIN inventory i ON p.id = i.product_id
		WHERE p.category_id = $1
		LIMIT $2`,
		categoryID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var createdAt, updatedAt time.Time

		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.CategoryID, &product.Price,
			&createdAt, &updatedAt, &product.Quantity)
		if err != nil {
			return nil, err
		}

		product.CreatedAt = createdAt
		product.UpdatedAt = updatedAt

		// Get product tags
		tagRows, err := r.db.Query("SELECT tag FROM product_tags WHERE product_id = $1", product.ID)
		if err == nil {
			defer tagRows.Close()
			var tags []string
			for tagRows.Next() {
				var tag string
				if err := tagRows.Scan(&tag); err == nil {
					tags = append(tags, tag)
				}
			}
			product.Tags = tags
		}

		products = append(products, product)
	}

	return products, nil
}

// GetProductsByTags retrieves products by tags with their inventory information
func (r *InventoryRepository) GetProductsByTags(tags []string, limit int) ([]models.Product, error) {
	// Get products by tags with inventory information using JOIN
	query := `
		SELECT p.id, p.name, p.description, p.category_id, p.price, p.created_at, p.updated_at,
		COALESCE(i.quantity, 0) as quantity
		FROM products p
		JOIN product_tags t ON p.id = t.product_id
		LEFT JOIN inventory i ON p.id = i.product_id
		WHERE t.tag IN (`

	// Create placeholders for tags
	var placeholders []string
	var args []any
	for i, tag := range tags {
		placeholders = append(placeholders, "$"+strconv.Itoa(i+1))
		args = append(args, tag)
	}
	query += strings.Join(placeholders, ",") + ") GROUP BY p.id, i.quantity LIMIT $" + strconv.Itoa(len(tags)+1)
	args = append(args, limit)

	// Execute query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var createdAt, updatedAt time.Time

		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.CategoryID, &product.Price,
			&createdAt, &updatedAt, &product.Quantity)
		if err != nil {
			return nil, err
		}

		product.CreatedAt = createdAt
		product.UpdatedAt = updatedAt

		// Get product tags
		tagRows, err := r.db.Query("SELECT tag FROM product_tags WHERE product_id = $1", product.ID)
		if err == nil {
			defer tagRows.Close()
			var tags []string
			for tagRows.Next() {
				var tag string
				if err := tagRows.Scan(&tag); err == nil {
					tags = append(tags, tag)
				}
			}
			product.Tags = tags
		}

		products = append(products, product)
	}

	return products, nil
}

// GenerateID generates a new UUID
func GenerateID() string {
	return uuid.New().String()
}
