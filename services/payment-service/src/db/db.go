package db

import (
"database/sql"
"fmt"
"log"
"time"

"github.com/online-order-system/payment-service/config"
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
// Check if payments table exists
var exists bool
err := db.QueryRow(`
    SELECT EXISTS (
        SELECT FROM information_schema.tables
        WHERE table_schema = 'public'
        AND table_name = 'payments'
    )
`).Scan(&exists)
if err != nil {
    log.Printf("Error checking if payments table exists: %v", err)
    return err
}

// If table exists, check if columns need to be altered
if exists {
    // Check and alter payment_method column
    var columnType string
    err = db.QueryRow(`
        SELECT data_type || '(' || character_maximum_length || ')'
        FROM information_schema.columns
        WHERE table_name = 'payments'
        AND column_name = 'payment_method'
    `).Scan(&columnType)

    if err != nil {
        log.Printf("Error checking payment_method column type: %v", err)
    } else {
        log.Printf("Current payment_method column type: %s", columnType)

        // If column type is not VARCHAR(50), alter it
        if columnType != "character varying(50)" {
            log.Printf("Altering payment_method column to VARCHAR(50)")
            _, err = db.Exec(`ALTER TABLE payments ALTER COLUMN payment_method TYPE VARCHAR(50)`)
            if err != nil {
                log.Printf("Error altering payment_method column: %v", err)
                return err
            }
            log.Println("Successfully altered payment_method column to VARCHAR(50)")
        }
    }

    // Check and alter cvv column
    err = db.QueryRow(`
        SELECT data_type || '(' || character_maximum_length || ')'
        FROM information_schema.columns
        WHERE table_name = 'payments'
        AND column_name = 'cvv'
    `).Scan(&columnType)

    if err != nil {
        log.Printf("Error checking cvv column type: %v", err)
    } else {
        log.Printf("Current cvv column type: %s", columnType)

        // If column type is not VARCHAR(100), alter it
        if columnType != "character varying(100)" {
            log.Printf("Altering cvv column to VARCHAR(100)")
            _, err = db.Exec(`ALTER TABLE payments ALTER COLUMN cvv TYPE VARCHAR(100)`)
            if err != nil {
                log.Printf("Error altering cvv column: %v", err)
                return err
            }
            log.Println("Successfully altered cvv column to VARCHAR(100)")
        }
    }

    // Check and alter expiry_month column
    err = db.QueryRow(`
        SELECT data_type || '(' || character_maximum_length || ')'
        FROM information_schema.columns
        WHERE table_name = 'payments'
        AND column_name = 'expiry_month'
    `).Scan(&columnType)

    if err != nil {
        log.Printf("Error checking expiry_month column type: %v", err)
    } else {
        log.Printf("Current expiry_month column type: %s", columnType)

        // If column type is not VARCHAR(10), alter it
        if columnType != "character varying(10)" {
            log.Printf("Altering expiry_month column to VARCHAR(10)")
            _, err = db.Exec(`ALTER TABLE payments ALTER COLUMN expiry_month TYPE VARCHAR(10)`)
            if err != nil {
                log.Printf("Error altering expiry_month column: %v", err)
                return err
            }
            log.Println("Successfully altered expiry_month column to VARCHAR(10)")
        }
    }

    // Check and alter expiry_year column
    err = db.QueryRow(`
        SELECT data_type || '(' || character_maximum_length || ')'
        FROM information_schema.columns
        WHERE table_name = 'payments'
        AND column_name = 'expiry_year'
    `).Scan(&columnType)

    if err != nil {
        log.Printf("Error checking expiry_year column type: %v", err)
    } else {
        log.Printf("Current expiry_year column type: %s", columnType)

        // If column type is not VARCHAR(10), alter it
        if columnType != "character varying(10)" {
            log.Printf("Altering expiry_year column to VARCHAR(10)")
            _, err = db.Exec(`ALTER TABLE payments ALTER COLUMN expiry_year TYPE VARCHAR(10)`)
            if err != nil {
                log.Printf("Error altering expiry_year column: %v", err)
                return err
            }
            log.Println("Successfully altered expiry_year column to VARCHAR(10)")
        }
    }

    // Check if Stripe-specific columns exist and add them if they don't
    var stripePaymentIDExists bool
    err = db.QueryRow(`
        SELECT EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_name = 'payments'
            AND column_name = 'stripe_payment_id'
        )
    `).Scan(&stripePaymentIDExists)

    if err != nil {
        log.Printf("Error checking if stripe_payment_id column exists: %v", err)
    } else if !stripePaymentIDExists {
        // Add Stripe-specific columns
        log.Println("Adding Stripe-specific columns to payments table")
        _, err = db.Exec(`
            ALTER TABLE payments
            ADD COLUMN stripe_payment_id VARCHAR(100),
            ADD COLUMN stripe_client_secret VARCHAR(255),
            ADD COLUMN currency VARCHAR(10),
            ADD COLUMN description TEXT,
            ADD COLUMN customer_email VARCHAR(255),
            ADD COLUMN customer_name VARCHAR(255),
            ADD COLUMN receipt_url VARCHAR(255),
            ADD COLUMN error_message TEXT
        `)
        if err != nil {
            log.Printf("Error adding Stripe-specific columns: %v", err)
            return err
        }
        log.Println("Successfully added Stripe-specific columns to payments table")
    }
} else {
    // Create payments table if it doesn't exist
    _, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    card_number VARCHAR(255),
    expiry_month VARCHAR(10),
    expiry_year VARCHAR(10),
    cvv VARCHAR(100),
    stripe_payment_id VARCHAR(100),
    stripe_client_secret VARCHAR(255),
    currency VARCHAR(10),
    description TEXT,
    customer_email VARCHAR(255),
    customer_name VARCHAR(255),
    receipt_url VARCHAR(255),
    error_message TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
    )
    `)
    if err != nil {
        log.Printf("Error creating payments table: %v", err)
        return err
    }
    log.Println("Payments table created successfully")
}
log.Println("Database tables created or already exist")
return nil
}
