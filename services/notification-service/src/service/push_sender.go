package service

import (
"log"

"github.com/online-order-system/notification-service/config"
"github.com/online-order-system/notification-service/interfaces"
)

// PushSender handles sending push notifications
type PushSender struct {
config *config.Config
}

// Ensure PushSender implements PushSender interface
var _ interfaces.PushSender = (*PushSender)(nil)

// NewPushSender creates a new push notification sender
func NewPushSender(cfg *config.Config) *PushSender {
return &PushSender{
config: cfg,
}
}

// SendPush sends a push notification
func (s *PushSender) SendPush(to, title, message string) error {
// In a real system, we would use a library like firebase-admin-go to send push notifications
// For demo purposes, we'll just log the push notification
log.Printf("Sending push notification to %s with title '%s': %s", to, title, message)
log.Printf("Push settings: API Key: %s", s.config.PushAPIKey)

// Simulate success (in a real system, this would be the result of the actual push notification sending)
return nil
}
