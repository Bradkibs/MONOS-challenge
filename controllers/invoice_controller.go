package controllers

import (
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/services"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvoiceController struct {
	DB *pgxpool.Pool
}

func (ic *InvoiceController) AddInvoice(c *fiber.Ctx) error {
	var invoice models.Invoice
	if err := c.BodyParser(&invoice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	invoice.ID = utils.GenerateUniqueID()
	err := services.AddInvoice(&invoice, ic.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(invoice)
}

func (ic *InvoiceController) GetInvoiceByID(c *fiber.Ctx) error {
	invoiceID := c.Params("id")
	id, err := uuid.Parse(invoiceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid invoice ID"})
	}

	invoice, err := services.GetInvoiceByID(id, ic.DB)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invoice not found"})
	}

	return c.Status(fiber.StatusOK).JSON(invoice)
}

func (ic *InvoiceController) UpdateInvoice(c *fiber.Ctx) error {
	var invoice models.Invoice
	if err := c.BodyParser(&invoice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.UpdateInvoice(&invoice, ic.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(invoice)
}

func (ic *InvoiceController) DeleteInvoice(c *fiber.Ctx) error {
	invoiceID := c.Params("id")
	id, err := uuid.Parse(invoiceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid invoice ID"})
	}

	err = services.DeleteInvoice(id, ic.DB)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invoice not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Invoice deleted successfully"})
}

func (ic *InvoiceController) GenerateInvoiceForPayment(c *fiber.Ctx) error {
	paymentID := c.Params("payment_id")
	userID := c.Params("user_id")
	paymentUUID, err := uuid.Parse(paymentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payment ID"})
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	invoice, err := services.GenerateInvoiceForPayment(paymentUUID, userUUID, ic.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(invoice)
}

func (ic *InvoiceController) GetAllInvoices(c *fiber.Ctx) error {
	invoices, err := services.GetAllInvoices(ic.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(invoices)
}
