package routes

import (
	"github.com/Bradkibs/MONOS-challenge/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupBusinessRoutes(app *fiber.App, db *pgxpool.Pool) {

	businessController := controllers.BusinessController{DB: db}

	businessGroup := app.Group("/businesses")

	businessGroup.Get("/", businessController.GetAllBusinesses)
	businessGroup.Post("/create", businessController.CreateBusiness)
	businessGroup.Get("/:business_id", businessController.GetBusinessByID)
	businessGroup.Put("/update", businessController.UpdateBusiness)
	businessGroup.Delete("/delete/:business_id", businessController.DeleteBusiness)
	businessGroup.Get("/vendor/:vendor_id", businessController.GetBusinessesByVendorID)
}
