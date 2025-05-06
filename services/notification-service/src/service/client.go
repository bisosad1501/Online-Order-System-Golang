package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/online-order-system/notification-service/config"
)

// OrderClient is a client for the order service
type OrderClient struct {
	config *config.Config
	client *http.Client
}

// NewOrderClient creates a new order client
func NewOrderClient(cfg *config.Config) *OrderClient {
	return &OrderClient{
		config: cfg,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// OrderResponse represents a response from the order service
type OrderResponse struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customer_id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetOrderByID retrieves an order by ID
func (c *OrderClient) GetOrderByID(orderID string) (*OrderResponse, error) {
	url := fmt.Sprintf("%s/orders/%s", c.config.OrderServiceURL, orderID)
	
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get order, status code: %d", resp.StatusCode)
	}

	var order OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	return &order, nil
}
