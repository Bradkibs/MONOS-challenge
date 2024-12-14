package routes

import (
	"github.com/Bradkibs/MONOS-challenge/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupAuthRoutes(app *fiber.App, db *pgxpool.Pool) {

	authController := controllers.AuthController{DB: db}

	authGroup := app.Group("/auth")

	authGroup.Post("/register/email", authController.RegisterByMail)
	authGroup.Post("/register/phone", authController.RegisterByPhoneNumber)
	authGroup.Post("/login/email", authController.LoginByMail)
	authGroup.Post("/login/phone", authController.LoginByPhoneNumber)
	authGroup.Get("/validate", authController.ValidateToken)
}
