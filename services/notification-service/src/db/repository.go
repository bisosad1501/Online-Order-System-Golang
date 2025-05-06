package db

import (
"database/sql"
"github.com/google/uuid"
"github.com/online-order-system/notification-service/models"
"time"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
db *Database
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *Database) *NotificationRepository {
return &NotificationRepository{db: db}
}

// GenerateID generates a new UUID
func GenerateID() string {
return uuid.New().String()
}

// CreateNotification creates a new notification in the database
func (r *NotificationRepository) CreateNotification(notification models.Notification) error {
_, err := r.db.Exec(
"INSERT INTO notifications (id, user_id, type, status, subject, content, recipient, created_at, updated_at, sent_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
notification.ID, notification.CustomerID, notification.Type, notification.Status, notification.Subject, notification.Content, notification.Recipient, notification.CreatedAt, notification.UpdatedAt, notification.SentAt,
)
return err
}

// GetNotificationByID retrieves a notification by ID
func (r *NotificationRepository) GetNotificationByID(id string) (models.Notification, error) {
var notification models.Notification
var notificationType, status string
var createdAt, updatedAt time.Time
var sentAt sql.NullTime

err := r.db.QueryRow(
"SELECT id, user_id, type, status, subject, content, recipient, created_at, updated_at, sent_at FROM notifications WHERE id = $1",
id,
).Scan(&notification.ID, &notification.CustomerID, &notificationType, &status, &notification.Subject, &notification.Content, &notification.Recipient, &createdAt, &updatedAt, &sentAt)
if err != nil {
return notification, err
}

notification.Type = models.NotificationType(notificationType)
notification.Status = models.NotificationStatus(status)
notification.CreatedAt = createdAt
notification.UpdatedAt = updatedAt
if sentAt.Valid {
sentTime := sentAt.Time
notification.SentAt = &sentTime
}

return notification, nil
}

// GetNotifications retrieves all notifications
func (r *NotificationRepository) GetNotifications() ([]models.Notification, error) {
rows, err := r.db.Query(
"SELECT id, user_id, type, status, subject, content, recipient, created_at, updated_at, sent_at FROM notifications ORDER BY created_at DESC",
)
if err != nil {
return nil, err
}
defer rows.Close()

var notifications []models.Notification
for rows.Next() {
var notification models.Notification
var notificationType, status string
var createdAt, updatedAt time.Time
var sentAt sql.NullTime

err := rows.Scan(&notification.ID, &notification.CustomerID, &notificationType, &status, &notification.Subject, &notification.Content, &notification.Recipient, &createdAt, &updatedAt, &sentAt)
if err != nil {
return nil, err
}

notification.Type = models.NotificationType(notificationType)
notification.Status = models.NotificationStatus(status)
notification.CreatedAt = createdAt
notification.UpdatedAt = updatedAt
if sentAt.Valid {
sentTime := sentAt.Time
notification.SentAt = &sentTime
}

notifications = append(notifications, notification)
}

return notifications, nil
}

// GetNotificationsByCustomerID retrieves all notifications for a customer
func (r *NotificationRepository) GetNotificationsByCustomerID(customerID string) ([]models.Notification, error) {
rows, err := r.db.Query(
"SELECT id, user_id, type, status, subject, content, recipient, created_at, updated_at, sent_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC",
customerID,
)
if err != nil {
return nil, err
}
defer rows.Close()

var notifications []models.Notification
for rows.Next() {
var notification models.Notification
var notificationType, status string
var createdAt, updatedAt time.Time
var sentAt sql.NullTime

err := rows.Scan(&notification.ID, &notification.CustomerID, &notificationType, &status, &notification.Subject, &notification.Content, &notification.Recipient, &createdAt, &updatedAt, &sentAt)
if err != nil {
return nil, err
}

notification.Type = models.NotificationType(notificationType)
notification.Status = models.NotificationStatus(status)
notification.CreatedAt = createdAt
notification.UpdatedAt = updatedAt
if sentAt.Valid {
sentTime := sentAt.Time
notification.SentAt = &sentTime
}

notifications = append(notifications, notification)
}

return notifications, nil
}

// UpdateNotificationStatus updates the status of a notification
func (r *NotificationRepository) UpdateNotificationStatus(id string, status models.NotificationStatus) error {
now := time.Now()
var sentAt interface{} = nil
if status == models.NotificationStatusSent {
sentAt = now
}

_, err := r.db.Exec(
"UPDATE notifications SET status = $1, updated_at = $2, sent_at = $3 WHERE id = $4",
status, now, sentAt, id,
)
return err
}
