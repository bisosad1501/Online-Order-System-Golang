package service

import (
"log"

"github.com/online-order-system/notification-service/config"
"github.com/online-order-system/notification-service/interfaces"
)

// EmailSender handles sending emails
type EmailSender struct {
config *config.Config
}

// Ensure EmailSender implements EmailSender interface
var _ interfaces.EmailSender = (*EmailSender)(nil)

// NewEmailSender creates a new email sender
func NewEmailSender(cfg *config.Config) *EmailSender {
return &EmailSender{
config: cfg,
}
}

// SendEmail sends an email
func (s *EmailSender) SendEmail(to, subject, body string) error {
// In a real system, we would use a library like gomail to send emails
// For demo purposes, we'll just log the email
log.Printf("Sending email to %s with subject '%s': %s", to, subject, body)
log.Printf("SMTP settings: %s:%s, %s, %s", s.config.SMTPHost, s.config.SMTPPort, s.config.SMTPUsername, s.config.SMTPFrom)

// Simulate success (in a real system, this would be the result of the actual email sending)
return nil
}
