package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/online-order-system/payment-service/models"
	"github.com/online-order-system/payment-service/utils"
)

// PaymentRepository handles database operations for payments
type PaymentRepository struct {
	db *Database
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *Database) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// CreatePayment creates a new payment in the database
func (r *PaymentRepository) CreatePayment(payment models.Payment) error {
	// Log the payment
	log.Printf("Repository: Creating payment: %+v", payment)

	// Encrypt sensitive payment information if provided
	var encryptedCardNumber, encryptedCVV string
	var err error

	if payment.CardNumber != "" {
		encryptedCardNumber, err = utils.Encrypt(payment.CardNumber)
		if err != nil {
			log.Printf("Error encrypting card number: %v", err)
			return err
		}
	}

	if payment.CVV != "" {
		encryptedCVV, err = utils.Encrypt(payment.CVV)
		if err != nil {
			log.Printf("Error encrypting CVV: %v", err)
			return err
		}
	}

	log.Printf("Executing SQL: INSERT INTO payments with values: id=%s, order_id=%s, amount=%f, status=%s, payment_method=%s",
		payment.ID, payment.OrderID, payment.Amount, payment.Status, payment.PaymentMethod)

	_, err = r.db.Exec(
		`INSERT INTO payments (
			id, order_id, amount, status, payment_method,
			card_number, expiry_month, expiry_year, cvv,
			stripe_payment_id, stripe_client_secret, currency, description,
			customer_email, customer_name, receipt_url, error_message,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		)`,
		payment.ID, payment.OrderID, payment.Amount, payment.Status, payment.PaymentMethod,
		encryptedCardNumber, payment.ExpiryMonth, payment.ExpiryYear, encryptedCVV,
		payment.StripePaymentID, payment.StripeClientSecret, payment.Currency, payment.Description,
		payment.CustomerEmail, payment.CustomerName, payment.ReceiptURL, payment.ErrorMessage,
		payment.CreatedAt, payment.UpdatedAt,
	)

	if err != nil {
		log.Printf("SQL Error: %v", err)
	}
	return err
}

// GetPaymentByID retrieves a payment by ID
func (r *PaymentRepository) GetPaymentByID(id string) (models.Payment, error) {
	var payment models.Payment
	var status string
	var createdAt, updatedAt time.Time
	var encryptedCardNumber, encryptedCVV string
	var stripePaymentID, stripeClientSecret, currency, description sql.NullString
	var customerEmail, customerName, receiptURL, errorMessage sql.NullString

	// Get payment
	err := r.db.QueryRow(
		`SELECT
			id, order_id, amount, status, payment_method,
			card_number, expiry_month, expiry_year, cvv,
			stripe_payment_id, stripe_client_secret, currency, description,
			customer_email, customer_name, receipt_url, error_message,
			created_at, updated_at
		FROM payments WHERE id = $1`,
		id,
	).Scan(&payment.ID, &payment.OrderID, &payment.Amount, &status, &payment.PaymentMethod,
		&encryptedCardNumber, &payment.ExpiryMonth, &payment.ExpiryYear, &encryptedCVV,
		&stripePaymentID, &stripeClientSecret, &currency, &description,
		&customerEmail, &customerName, &receiptURL, &errorMessage,
		&createdAt, &updatedAt)
	if err != nil {
		return payment, err
	}

	// Decrypt sensitive payment information
	var cardNumber string
	if encryptedCardNumber == "" {
		cardNumber = ""
	} else {
		decrypted, err := utils.Decrypt(encryptedCardNumber)
		if err != nil {
			log.Printf("Error decrypting card number: %v", err)
			// Continue with masked card number
			if len(encryptedCardNumber) > 4 {
				cardNumber = "************" + encryptedCardNumber[len(encryptedCardNumber)-4:]
			} else {
				cardNumber = "************"
			}
		} else {
			cardNumber = decrypted
		}
	}
	payment.CardNumber = cardNumber

	var cvv string
	if encryptedCVV == "" {
		cvv = ""
	} else {
		decrypted, err := utils.Decrypt(encryptedCVV)
		if err != nil {
			log.Printf("Error decrypting CVV: %v", err)
			// Continue with masked CVV
			cvv = "***"
		} else {
			cvv = decrypted
		}
	}
	payment.CVV = cvv

	// Set Stripe fields if they exist
	if stripePaymentID.Valid {
		payment.StripePaymentID = stripePaymentID.String
	}
	if stripeClientSecret.Valid {
		payment.StripeClientSecret = stripeClientSecret.String
	}
	if currency.Valid {
		payment.Currency = currency.String
	}
	if description.Valid {
		payment.Description = description.String
	}
	if customerEmail.Valid {
		payment.CustomerEmail = customerEmail.String
	}
	if customerName.Valid {
		payment.CustomerName = customerName.String
	}
	if receiptURL.Valid {
		payment.ReceiptURL = receiptURL.String
	}
	if errorMessage.Valid {
		payment.ErrorMessage = errorMessage.String
	}

	payment.Status = models.PaymentStatus(status)
	payment.CreatedAt = createdAt
	payment.UpdatedAt = updatedAt

	return payment, nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (r *PaymentRepository) GetPaymentByOrderID(orderID string) (models.Payment, error) {
	var payment models.Payment
	var status string
	var createdAt, updatedAt time.Time
	var encryptedCardNumber, encryptedCVV string
	var stripePaymentID, stripeClientSecret, currency, description sql.NullString
	var customerEmail, customerName, receiptURL, errorMessage sql.NullString

	// Get payment
	err := r.db.QueryRow(
		`SELECT
			id, order_id, amount, status, payment_method,
			card_number, expiry_month, expiry_year, cvv,
			stripe_payment_id, stripe_client_secret, currency, description,
			customer_email, customer_name, receipt_url, error_message,
			created_at, updated_at
		FROM payments WHERE order_id = $1`,
		orderID,
	).Scan(&payment.ID, &payment.OrderID, &payment.Amount, &status, &payment.PaymentMethod,
		&encryptedCardNumber, &payment.ExpiryMonth, &payment.ExpiryYear, &encryptedCVV,
		&stripePaymentID, &stripeClientSecret, &currency, &description,
		&customerEmail, &customerName, &receiptURL, &errorMessage,
		&createdAt, &updatedAt)
	if err != nil {
		return payment, err
	}

	// Decrypt sensitive payment information
	var cardNumber string
	if encryptedCardNumber == "" {
		cardNumber = ""
	} else {
		decrypted, err := utils.Decrypt(encryptedCardNumber)
		if err != nil {
			log.Printf("Error decrypting card number: %v", err)
			// Continue with masked card number
			if len(encryptedCardNumber) > 4 {
				cardNumber = "************" + encryptedCardNumber[len(encryptedCardNumber)-4:]
			} else {
				cardNumber = "************"
			}
		} else {
			cardNumber = decrypted
		}
	}
	payment.CardNumber = cardNumber

	var cvv string
	if encryptedCVV == "" {
		cvv = ""
	} else {
		decrypted, err := utils.Decrypt(encryptedCVV)
		if err != nil {
			log.Printf("Error decrypting CVV: %v", err)
			// Continue with masked CVV
			cvv = "***"
		} else {
			cvv = decrypted
		}
	}
	payment.CVV = cvv

	// Set Stripe fields if they exist
	if stripePaymentID.Valid {
		payment.StripePaymentID = stripePaymentID.String
	}
	if stripeClientSecret.Valid {
		payment.StripeClientSecret = stripeClientSecret.String
	}
	if currency.Valid {
		payment.Currency = currency.String
	}
	if description.Valid {
		payment.Description = description.String
	}
	if customerEmail.Valid {
		payment.CustomerEmail = customerEmail.String
	}
	if customerName.Valid {
		payment.CustomerName = customerName.String
	}
	if receiptURL.Valid {
		payment.ReceiptURL = receiptURL.String
	}
	if errorMessage.Valid {
		payment.ErrorMessage = errorMessage.String
	}

	payment.Status = models.PaymentStatus(status)
	payment.CreatedAt = createdAt
	payment.UpdatedAt = updatedAt

	return payment, nil
}

// GetPayments retrieves all payments
func (r *PaymentRepository) GetPayments() ([]models.Payment, error) {
	// Get all payments
	rows, err := r.db.Query(
		"SELECT id, order_id, amount, status, payment_method, card_number, expiry_month, expiry_year, cvv, created_at, updated_at FROM payments",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		var status string
		var createdAt, updatedAt time.Time
		var encryptedCardNumber, encryptedCVV string

		err := rows.Scan(&payment.ID, &payment.OrderID, &payment.Amount, &status, &payment.PaymentMethod,
			&encryptedCardNumber, &payment.ExpiryMonth, &payment.ExpiryYear, &encryptedCVV,
			&createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		// Decrypt sensitive payment information
		var cardNumber string
		if encryptedCardNumber == "" {
			cardNumber = ""
		} else {
			decrypted, err := utils.Decrypt(encryptedCardNumber)
			if err != nil {
				log.Printf("Error decrypting card number: %v", err)
				// Continue with masked card number
				if len(encryptedCardNumber) > 4 {
					cardNumber = "************" + encryptedCardNumber[len(encryptedCardNumber)-4:]
				} else {
					cardNumber = "************"
				}
			} else {
				cardNumber = decrypted
			}
		}
		payment.CardNumber = cardNumber

		var cvv string
		if encryptedCVV == "" {
			cvv = ""
		} else {
			decrypted, err := utils.Decrypt(encryptedCVV)
			if err != nil {
				log.Printf("Error decrypting CVV: %v", err)
				// Continue with masked CVV
				cvv = "***"
			} else {
				cvv = decrypted
			}
		}
		payment.CVV = cvv

		payment.Status = models.PaymentStatus(status)
		payment.CreatedAt = createdAt
		payment.UpdatedAt = updatedAt

		payments = append(payments, payment)
	}

	return payments, nil
}

// UpdatePaymentStatus updates the status of a payment
func (r *PaymentRepository) UpdatePaymentStatus(id string, status models.PaymentStatus) error {
	_, err := r.db.Exec(
		"UPDATE payments SET status = $1, updated_at = $2 WHERE id = $3",
		status, time.Now(), id,
	)
	return err
}

// UpdatePaymentStripeInfo updates the Stripe-specific information for a payment
func (r *PaymentRepository) UpdatePaymentStripeInfo(id, stripePaymentID, stripeClientSecret string) error {
	_, err := r.db.Exec(
		"UPDATE payments SET stripe_payment_id = $1, stripe_client_secret = $2, updated_at = $3 WHERE id = $4",
		stripePaymentID, stripeClientSecret, time.Now(), id,
	)
	return err
}

// UpdatePaymentReceiptURL updates the receipt URL for a payment
func (r *PaymentRepository) UpdatePaymentReceiptURL(id, receiptURL string) error {
	_, err := r.db.Exec(
		"UPDATE payments SET receipt_url = $1, updated_at = $2 WHERE id = $3",
		receiptURL, time.Now(), id,
	)
	return err
}

// UpdatePaymentStatusWithError updates the status and error message of a payment
func (r *PaymentRepository) UpdatePaymentStatusWithError(id string, status models.PaymentStatus, errorMessage string) error {
	_, err := r.db.Exec(
		"UPDATE payments SET status = $1, error_message = $2, updated_at = $3 WHERE id = $4",
		status, errorMessage, time.Now(), id,
	)
	return err
}

// GenerateID generates a new UUID
func GenerateID() string {
	return uuid.New().String()
}
