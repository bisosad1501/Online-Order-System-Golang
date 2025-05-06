package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/online-order-system/user-service/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *Database
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user models.Customer) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, email, first_name, last_name, phone, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID, user.Email, user.FirstName, user.LastName, user.Phone, user.Address,
		user.CreatedAt, user.UpdatedAt,
	)
	return err
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id string) (models.Customer, error) {
	var user models.Customer
	err := r.db.QueryRow(
		`SELECT id, email, first_name, last_name, phone, address, created_at, updated_at
		FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Customer{}, errors.New("user not found")
	}
	return user, err
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (models.Customer, error) {
	var user models.Customer
	err := r.db.QueryRow(
		`SELECT id, email, first_name, last_name, phone, address, created_at, updated_at
		FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Customer{}, errors.New("user not found")
	}
	return user, err
}

// GetUsers retrieves all users
func (r *UserRepository) GetUsers() ([]models.Customer, error) {
	rows, err := r.db.Query(
		`SELECT id, email, first_name, last_name, phone, address, created_at, updated_at
		FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.Customer
	for rows.Next() {
		var user models.Customer
		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName,
			&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUser updates a user in the database
func (r *UserRepository) UpdateUser(user models.Customer) error {
	_, err := r.db.Exec(
		`UPDATE users SET email = $1, first_name = $2, last_name = $3, phone = $4, address = $5, updated_at = $6
		WHERE id = $7`,
		user.Email, user.FirstName, user.LastName, user.Phone, user.Address,
		user.UpdatedAt, user.ID,
	)
	return err
}

// DeleteUser deletes a user from the database
func (r *UserRepository) DeleteUser(id string) error {
	// First delete from user_orders
	_, err := r.db.Exec("DELETE FROM user_orders WHERE user_id = $1", id)
	if err != nil {
		return err
	}

	// Then delete from users
	_, err = r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// GetUserOrders retrieves all orders for a user
func (r *UserRepository) GetUserOrders(userID string) ([]models.CustomerOrder, error) {
	rows, err := r.db.Query(
		`SELECT user_id, order_id, order_date, order_status
		FROM user_orders WHERE user_id = $1 ORDER BY order_date DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.CustomerOrder
	for rows.Next() {
		var order models.CustomerOrder
		err := rows.Scan(
			&order.CustomerID, &order.OrderID, &order.OrderDate, &order.OrderStatus,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// AddUserOrder adds an order to a user
func (r *UserRepository) AddUserOrder(userID string, orderID string, orderStatus string) error {
	_, err := r.db.Exec(
		`INSERT INTO user_orders (user_id, order_id, order_date, order_status)
		VALUES ($1, $2, $3, $4)`,
		userID, orderID, time.Now(), orderStatus,
	)
	return err
}

// UpdateUserOrderStatus updates the status of a user's order
func (r *UserRepository) UpdateUserOrderStatus(userID string, orderID string, orderStatus string) error {
	_, err := r.db.Exec(
		`UPDATE user_orders SET order_status = $1 WHERE user_id = $2 AND order_id = $3`,
		orderStatus, userID, orderID,
	)
	return err
}
