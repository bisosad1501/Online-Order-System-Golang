package db

import (
"database/sql"
"time"

"github.com/google/uuid"
"github.com/online-order-system/shipping-service/models"
)

// ShippingRepository handles database operations for shipments
type ShippingRepository struct {
db *Database
}

// NewShippingRepository creates a new shipping repository
func NewShippingRepository(db *Database) *ShippingRepository {
return &ShippingRepository{db: db}
}

// CreateShipment creates a new shipment in the database
func (r *ShippingRepository) CreateShipment(shipment models.Shipment) error {
_, err := r.db.Exec(
"INSERT INTO shipments (id, order_id, status, tracking_number, shipping_address, carrier, estimated_delivery, created_at, updated_at, customer_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
shipment.ID, shipment.OrderID, shipment.Status, shipment.TrackingNumber, shipment.ShippingAddress, shipment.Carrier, shipment.EstimatedDelivery, shipment.CreatedAt, shipment.UpdatedAt, shipment.CustomerID,
)
return err
}

// GetShipmentByID retrieves a shipment by ID
func (r *ShippingRepository) GetShipmentByID(id string) (models.Shipment, error) {
var shipment models.Shipment
var status string
var createdAt, updatedAt time.Time
var estimatedDeliveryNull sql.NullTime

// Get shipment
err := r.db.QueryRow(
"SELECT id, order_id, status, tracking_number, shipping_address, carrier, estimated_delivery, created_at, updated_at, customer_id FROM shipments WHERE id = $1",
id,
).Scan(&shipment.ID, &shipment.OrderID, &status, &shipment.TrackingNumber, &shipment.ShippingAddress, &shipment.Carrier, &estimatedDeliveryNull, &createdAt, &updatedAt, &shipment.CustomerID)
if err != nil {
return shipment, err
}

shipment.Status = models.ShipmentStatus(status)
shipment.CreatedAt = createdAt
shipment.UpdatedAt = updatedAt
if estimatedDeliveryNull.Valid {
shipment.EstimatedDelivery = estimatedDeliveryNull.Time
}

return shipment, nil
}

// GetShipmentByOrderID retrieves a shipment by order ID
func (r *ShippingRepository) GetShipmentByOrderID(orderID string) (models.Shipment, error) {
var shipment models.Shipment
var status string
var createdAt, updatedAt time.Time
var estimatedDeliveryNull sql.NullTime

// Get shipment
err := r.db.QueryRow(
"SELECT id, order_id, status, tracking_number, shipping_address, carrier, estimated_delivery, created_at, updated_at, customer_id FROM shipments WHERE order_id = $1",
orderID,
).Scan(&shipment.ID, &shipment.OrderID, &status, &shipment.TrackingNumber, &shipment.ShippingAddress, &shipment.Carrier, &estimatedDeliveryNull, &createdAt, &updatedAt, &shipment.CustomerID)
if err != nil {
return shipment, err
}

shipment.Status = models.ShipmentStatus(status)
shipment.CreatedAt = createdAt
shipment.UpdatedAt = updatedAt
if estimatedDeliveryNull.Valid {
shipment.EstimatedDelivery = estimatedDeliveryNull.Time
}

return shipment, nil
}

// GetShipments retrieves all shipments
func (r *ShippingRepository) GetShipments() ([]models.Shipment, error) {
// Get all shipments
rows, err := r.db.Query(
"SELECT id, order_id, status, tracking_number, shipping_address, carrier, estimated_delivery, created_at, updated_at, customer_id FROM shipments",
)
if err != nil {
return nil, err
}
defer rows.Close()

var shipments []models.Shipment
for rows.Next() {
var shipment models.Shipment
var status string
var createdAt, updatedAt time.Time
var estimatedDeliveryNull sql.NullTime

err := rows.Scan(&shipment.ID, &shipment.OrderID, &status, &shipment.TrackingNumber, &shipment.ShippingAddress, &shipment.Carrier, &estimatedDeliveryNull, &createdAt, &updatedAt, &shipment.CustomerID)
if err != nil {
return nil, err
}

shipment.Status = models.ShipmentStatus(status)
shipment.CreatedAt = createdAt
shipment.UpdatedAt = updatedAt
if estimatedDeliveryNull.Valid {
shipment.EstimatedDelivery = estimatedDeliveryNull.Time
}

shipments = append(shipments, shipment)
}

return shipments, nil
}

// UpdateShipmentStatus updates the status of a shipment
func (r *ShippingRepository) UpdateShipmentStatus(id string, status models.ShipmentStatus) error {
_, err := r.db.Exec(
"UPDATE shipments SET status = $1, updated_at = $2 WHERE id = $3",
status, time.Now(), id,
)
return err
}

// UpdateTrackingNumber updates the tracking number of a shipment
func (r *ShippingRepository) UpdateTrackingNumber(id string, trackingNumber string) error {
_, err := r.db.Exec(
"UPDATE shipments SET tracking_number = $1, updated_at = $2 WHERE id = $3",
trackingNumber, time.Now(), id,
)
return err
}

// GenerateID generates a new UUID
func GenerateID() string {
return uuid.New().String()
}
