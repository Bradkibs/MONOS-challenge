package utils

import (
	"fmt"
	"github.com/google/uuid"
	"net/smtp"
)

// GenerateUniqueID generates a new unique UUID string.
func GenerateUniqueID() uuid.UUID {
	return uuid.New()
}

func SendEmail(to string, subject string, message string) error {
	// Define your SMTP server details
	smtpHost := "smtp.example.com"
	smtpPort := "587"
	smtpUser := "your-email@example.com"
	smtpPass := "your-password"

	// Sender and recipient
	from := smtpUser
	recipients := []string{to}

	// Build the email body
	body := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, message)

	// Authenticate with the SMTP server
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, []byte(body))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
