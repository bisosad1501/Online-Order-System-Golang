package db

import (
"database/sql"
"fmt"
"log"
"time"

"github.com/online-order-system/shipping-service/config"
_ "github.com/lib/pq"
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
// Create shipments table
_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS shipments (
id VARCHAR(36) PRIMARY KEY,
order_id VARCHAR(36) NOT NULL,
status VARCHAR(20) NOT NULL,
tracking_number VARCHAR(100),
shipping_address TEXT NOT NULL,
carrier VARCHAR(100),
estimated_delivery TIMESTAMP,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
customer_id VARCHAR(36)
)
`)
if err != nil {
return err
}

log.Println("Database tables created or already exist")
return nil
}
