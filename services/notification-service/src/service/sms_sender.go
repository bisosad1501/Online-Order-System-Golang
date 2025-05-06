package service

import (
"log"

"github.com/online-order-system/notification-service/config"
"github.com/online-order-system/notification-service/interfaces"
)

// SMSSender handles sending SMS messages
type SMSSender struct {
config *config.Config
}

// Ensure SMSSender implements SMSSender interface
var _ interfaces.SMSSender = (*SMSSender)(nil)

// NewSMSSender creates a new SMS sender
func NewSMSSender(cfg *config.Config) *SMSSender {
return &SMSSender{
config: cfg,
}
}

// SendSMS sends an SMS message
func (s *SMSSender) SendSMS(to, message string) error {
// In a real system, we would use a library like twilio-go to send SMS messages
// For demo purposes, we'll just log the SMS
log.Printf("Sending SMS to %s: %s", to, message)
log.Printf("SMS settings: API Key: %s, From: %s", s.config.SMSAPIKey, s.config.SMSFromNumber)

// Simulate success (in a real system, this would be the result of the actual SMS sending)
return nil
}
