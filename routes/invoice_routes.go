package routes

import (
	"github.com/Bradkibs/MONOS-challenge/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupInvoiceRoutes(app *fiber.App, db *pgxpool.Pool) {

	invoiceController := controllers.InvoiceController{DB: db}

	invoiceGroup := app.Group("/invoices")

	invoiceGroup.Get("/", invoiceController.GetAllInvoices)
	invoiceGroup.Post("/create", invoiceController.AddInvoice)
	invoiceGroup.Get("/:invoice_id", invoiceController.GetInvoiceByID)
	invoiceGroup.Put("/update", invoiceController.UpdateInvoice)
	invoiceGroup.Delete("/delete/:invoice_id", invoiceController.DeleteInvoice)
	invoiceGroup.Post("/generate/:payment_id/:user_id", invoiceController.GenerateInvoiceForPayment)
}
