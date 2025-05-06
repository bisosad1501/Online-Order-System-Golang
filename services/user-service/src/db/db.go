package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/online-order-system/user-service/config"
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
	// Create users table
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(100) NOT NULL UNIQUE,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		phone VARCHAR(20),
		address TEXT,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	)
	`)
	if err != nil {
		return err
	}

	// Create user_orders table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS user_orders (
		user_id VARCHAR(36) REFERENCES users(id),
		order_id VARCHAR(36) NOT NULL,
		order_date TIMESTAMP NOT NULL,
		order_status VARCHAR(20) NOT NULL,
		PRIMARY KEY (user_id, order_id)
	)
	`)
	if err != nil {
		return err
	}

	log.Println("Database tables created or already exist")
	return nil
}
