package db

import (
"database/sql"
"fmt"
"log"
"time"

"github.com/online-order-system/notification-service/config"
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
// Create notifications table
_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS notifications (
id VARCHAR(36) PRIMARY KEY,
user_id VARCHAR(36) NOT NULL,
type VARCHAR(20) NOT NULL,
status VARCHAR(20) NOT NULL,
content TEXT NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
sent_at TIMESTAMP
)
`)
if err != nil {
return err
}

log.Println("Database tables created or already exist")
return nil
}
