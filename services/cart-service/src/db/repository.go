package db

import (
"database/sql"
"errors"
"time"

"github.com/online-order-system/cart-service/models"
)

// CartRepository handles database operations for carts
type CartRepository struct {
db *Database
}

// NewCartRepository creates a new cart repository
func NewCartRepository(db *Database) *CartRepository {
return &CartRepository{db: db}
}

// CreateCart creates a new cart in the database
func (r *CartRepository) CreateCart(cart models.Cart) error {
_, err := r.db.Exec(
`INSERT INTO carts (id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4)`,
cart.ID, cart.CustomerID, cart.CreatedAt, cart.UpdatedAt,
)
return err
}

// GetCartByID retrieves a cart by ID
func (r *CartRepository) GetCartByID(id string) (models.Cart, error) {
var cart models.Cart
err := r.db.QueryRow(
`SELECT id, user_id, created_at, updated_at
FROM carts WHERE id = $1`,
id,
).Scan(
&cart.ID, &cart.CustomerID, &cart.CreatedAt, &cart.UpdatedAt,
)
if err == sql.ErrNoRows {
return models.Cart{}, errors.New("cart not found")
}
if err != nil {
return models.Cart{}, err
}

// Get cart items
cart.Items, err = r.GetCartItems(id)
if err != nil {
return models.Cart{}, err
}

return cart, nil
}

// GetCartByUserID retrieves a cart by user ID
func (r *CartRepository) GetCartByUserID(userID string) (models.Cart, error) {
var cart models.Cart
err := r.db.QueryRow(
`SELECT id, user_id, created_at, updated_at
FROM carts WHERE user_id = $1`,
userID,
).Scan(
&cart.ID, &cart.CustomerID, &cart.CreatedAt, &cart.UpdatedAt,
)
if err == sql.ErrNoRows {
return models.Cart{}, errors.New("cart not found")
}
if err != nil {
return models.Cart{}, err
}

// Get cart items
cart.Items, err = r.GetCartItems(cart.ID)
if err != nil {
return models.Cart{}, err
}

return cart, nil
}

// GetCartItems retrieves all items for a cart
func (r *CartRepository) GetCartItems(cartID string) ([]models.CartItem, error) {
rows, err := r.db.Query(
`SELECT id, cart_id, product_id, quantity, price, created_at, updated_at
FROM cart_items WHERE cart_id = $1 ORDER BY created_at ASC`,
cartID,
)
if err != nil {
return nil, err
}
defer rows.Close()

var items []models.CartItem
for rows.Next() {
var item models.CartItem
err := rows.Scan(
&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.Price,
&item.CreatedAt, &item.UpdatedAt,
)
if err != nil {
return nil, err
}
items = append(items, item)
}
return items, nil
}

// AddCartItem adds an item to a cart
func (r *CartRepository) AddCartItem(item models.CartItem) error {
// Check if item already exists
var existingItemID string
err := r.db.QueryRow(
`SELECT id FROM cart_items WHERE cart_id = $1 AND product_id = $2`,
item.CartID, item.ProductID,
).Scan(&existingItemID)

if err == nil {
// Item exists, update quantity
_, err = r.db.Exec(
`UPDATE cart_items SET quantity = quantity + $1, price = $2, updated_at = $3
WHERE id = $4`,
item.Quantity, item.Price, item.UpdatedAt, existingItemID,
)
return err
} else if err != sql.ErrNoRows {
// Unexpected error
return err
}

// Item doesn't exist, insert new item
_, err = r.db.Exec(
`INSERT INTO cart_items (id, cart_id, product_id, quantity, price, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)`,
item.ID, item.CartID, item.ProductID, item.Quantity, item.Price,
item.CreatedAt, item.UpdatedAt,
)
return err
}

// UpdateCartItem updates an item in a cart
func (r *CartRepository) UpdateCartItem(cartID string, itemID string, quantity int) error {
now := time.Now()
_, err := r.db.Exec(
`UPDATE cart_items SET quantity = $1, updated_at = $2
WHERE id = $3 AND cart_id = $4`,
quantity, now, itemID, cartID,
)
return err
}

// RemoveCartItem removes an item from a cart
func (r *CartRepository) RemoveCartItem(cartID string, itemID string) error {
_, err := r.db.Exec(
`DELETE FROM cart_items WHERE id = $1 AND cart_id = $2`,
itemID, cartID,
)
return err
}

// DeleteCart deletes a cart and all its items
func (r *CartRepository) DeleteCart(id string) error {
// Start a transaction
tx, err := r.db.Begin()
if err != nil {
return err
}
defer tx.Rollback()

// Delete cart items
_, err = tx.Exec("DELETE FROM cart_items WHERE cart_id = $1", id)
if err != nil {
return err
}

// Delete cart
_, err = tx.Exec("DELETE FROM carts WHERE id = $1", id)
if err != nil {
return err
}

// Commit transaction
return tx.Commit()
}

// DeleteCartByUserID deletes a cart by user ID
func (r *CartRepository) DeleteCartByUserID(userID string) error {
// Get cart ID
var cartID string
err := r.db.QueryRow(
`SELECT id FROM carts WHERE user_id = $1`,
userID,
).Scan(&cartID)
if err == sql.ErrNoRows {
return nil // No cart to delete
}
if err != nil {
return err
}

// Delete cart
return r.DeleteCart(cartID)
}

// UpdateCartUpdatedAt updates the updated_at field of a cart
func (r *CartRepository) UpdateCartUpdatedAt(cartID string) error {
now := time.Now()
_, err := r.db.Exec(
`UPDATE carts SET updated_at = $1 WHERE id = $2`,
now, cartID,
)
return err
}
