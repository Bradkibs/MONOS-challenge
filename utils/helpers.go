package utils

import (
	"errors"
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

// StripeService interface defines the methods for Stripe payment processing
type StripeService interface {
	Charge(amount float64, currency, description string) (string, error)
}

// MpesaService interface defines the methods for M-Pesa payment processing
type MpesaService interface {
	ProcessExpressPayment(amount float64, phoneNumber, shortcode string) (string, error)
}

// MockStripeService is a mock implementation of StripeService
type MockStripeService struct{}

func (s *MockStripeService) Charge(amount float64, currency, description string) (string, error) {
	// Simulate successful payment processing
	if amount <= 0 {
		return "", errors.New("invalid amount")
	}
	return fmt.Sprintf("mock_stripe_charge_id_%f", amount), nil
}

// MockMpesaService is a mock implementation of MpesaService
type MockMpesaService struct{}

func (m *MockMpesaService) ProcessExpressPayment(amount float64, phoneNumber, shortcode string) (string, error) {
	// Simulate successful M-Pesa payment processing
	if amount <= 0 {
		return "", errors.New("invalid amount")
	}
	if phoneNumber == "" || shortcode == "" {
		return "", errors.New("invalid phone number or shortcode")
	}
	return fmt.Sprintf("mock_mpesa_transaction_id_%f", amount), nil
}

// NewMockStripeService creates a new instance of MockStripeService
func NewMockStripeService() StripeService {
	return &MockStripeService{}
}

// NewMockMpesaService creates a new instance of MockMpesaService
func NewMockMpesaService() MpesaService {
	return &MockMpesaService{}
}
