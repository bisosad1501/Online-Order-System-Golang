package db

import (
"database/sql"
"fmt"
"log"
"time"

"github.com/online-order-system/inventory-service/config"
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
// Create products table
_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS products (
id VARCHAR(36) PRIMARY KEY,
name VARCHAR(100) NOT NULL,
description TEXT NOT NULL,
category_id VARCHAR(36),
price DECIMAL(10, 2) NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL
)
`)
if err != nil {
return err
}

// Create inventory table
_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS inventory (
product_id VARCHAR(36) PRIMARY KEY REFERENCES products(id),
quantity INTEGER NOT NULL,
updated_at TIMESTAMP NOT NULL
)
`)
if err != nil {
return err
}

// Create product_tags table
_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS product_tags (
id VARCHAR(36) PRIMARY KEY,
product_id VARCHAR(36) REFERENCES products(id),
tag VARCHAR(50) NOT NULL,
created_at TIMESTAMP NOT NULL
)
`)
if err != nil {
return err
}

// Create product_similarities table
_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS product_similarities (
id VARCHAR(36) PRIMARY KEY,
product_id VARCHAR(36) REFERENCES products(id),
similar_product_id VARCHAR(36) REFERENCES products(id),
similarity_score DECIMAL(5, 4) NOT NULL,
created_at TIMESTAMP NOT NULL
)
`)
if err != nil {
return err
}

log.Println("Database tables created or already exist")
return nil
}
