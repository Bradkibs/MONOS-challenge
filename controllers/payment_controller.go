package controllers

import (
	"fmt"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/services"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentController struct {
	DB *pgxpool.Pool
}

func (pc *PaymentController) AddPayment(c *fiber.Ctx) error {
	var payment models.Payment
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	payment.ID = utils.GenerateUniqueID()
	err := services.AddPayment(&payment, pc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(payment)
}

func (pc *PaymentController) GetPaymentByID(c *fiber.Ctx) error {
	paymentID := c.Params("id")
	id, err := uuid.Parse(paymentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payment ID"})
	}

	payment, err := services.GetPaymentByID(id, pc.DB)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payment not found"})
	}

	return c.Status(fiber.StatusOK).JSON(payment)
}

func (pc *PaymentController) UpdatePayment(c *fiber.Ctx) error {
	var payment models.Payment
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.UpdatePayment(&payment, pc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(payment)
}

func (pc *PaymentController) DeletePayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")
	id, err := uuid.Parse(paymentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payment ID"})
	}

	err = services.DeletePayment(id, pc.DB)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payment not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment deleted successfully"})
}

func (pc *PaymentController) GetAllPayments(c *fiber.Ctx) error {
	payments, err := services.GetAllPayments(pc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(payments)
}
func (pc *PaymentController) ProcessPayment(c *fiber.Ctx) error {

	type PaymentRequest struct {
		Amount        float64 `json:"amount" validate:"required,gt=0"`
		PaymentMethod string  `json:"payment_method" validate:"required"`
		PhoneNumber   string  `json:"phone_number,omitempty"`
		Shortcode     string  `json:"shortcode,omitempty"`
		Description   string  `json:"description,omitempty"`
	}

	var paymentReq PaymentRequest
	if err := c.BodyParser(&paymentReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if paymentReq.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Amount must be greater than zero"})
	}

	if paymentReq.PaymentMethod == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payment method is required"})
	}

	payment := &models.Payment{
		Amount: paymentReq.Amount,
	}

	stripeService := utils.NewMockStripeService()
	mpesaService := utils.NewMockMpesaService()

	// Process the payment using the service
	err := services.ProcessPayment(payment, pc.DB, paymentReq.PaymentMethod, stripeService, mpesaService)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Optionally, send an email or perform other post-processing
	emailErr := utils.SendEmail("user@example.com", "Payment Processed", fmt.Sprintf("Your payment of %.2f has been successfully processed.", payment.Amount))
	if emailErr != nil {
		fmt.Printf("Failed to send email: %v\n", emailErr)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment processed successfully"})
}
