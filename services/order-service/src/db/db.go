package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/online-order-system/order-service/config"
)

// Database represents a database connection
type Database struct {
	*sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config) (*Database, error) {
	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to PostgreSQL database")

	return &Database{db}, nil
}

// CreateTables creates the necessary tables if they don't exist
func (db *Database) CreateTables() error {
	// Create orders table
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS orders (
id VARCHAR(36) PRIMARY KEY,
user_id VARCHAR(36) NOT NULL,
status VARCHAR(20) NOT NULL,
total_amount DECIMAL(10, 2) NOT NULL,
shipping_address TEXT NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
inventory_locked BOOLEAN DEFAULT FALSE,
payment_processed BOOLEAN DEFAULT FALSE,
shipping_scheduled BOOLEAN DEFAULT FALSE,
failure_reason TEXT,
tracking_number VARCHAR(100)
)
`)
	if err != nil {
		return err
	}

	// Add new columns if they don't exist (for backward compatibility)
	_, _ = db.Exec(`ALTER TABLE orders ADD COLUMN IF NOT EXISTS inventory_locked BOOLEAN DEFAULT FALSE`)
	_, _ = db.Exec(`ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_processed BOOLEAN DEFAULT FALSE`)
	_, _ = db.Exec(`ALTER TABLE orders ADD COLUMN IF NOT EXISTS shipping_scheduled BOOLEAN DEFAULT FALSE`)
	_, _ = db.Exec(`ALTER TABLE orders ADD COLUMN IF NOT EXISTS failure_reason TEXT`)

	// Create order_items table
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS order_items (
id VARCHAR(36) PRIMARY KEY,
order_id VARCHAR(36) NOT NULL REFERENCES orders(id),
product_id VARCHAR(36) NOT NULL,
quantity INTEGER NOT NULL,
price DECIMAL(10, 2) NOT NULL
)
`)
	if err != nil {
		return err
	}

	// Create audit_logs table
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS audit_logs (
id VARCHAR(36) PRIMARY KEY,
service_name VARCHAR(50) NOT NULL,
action VARCHAR(50) NOT NULL,
user_id VARCHAR(36) NOT NULL,
timestamp TIMESTAMP NOT NULL,
details TEXT
)
`)
	if err != nil {
		return err
	}

	log.Println("Database tables created or already exist")
	return nil
}
