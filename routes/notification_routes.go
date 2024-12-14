package routes

import (
	"github.com/Bradkibs/MONOS-challenge/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupNotificationRoutes(app *fiber.App, db *pgxpool.Pool) {
	notificationController := controllers.NotificationController{Pool: db}

	notificationGroup := app.Group("/notifications")

	notificationGroup.Post("/", notificationController.CreateNotification)
	notificationGroup.Get("/:notification_id", notificationController.GetNotificationByID)
	notificationGroup.Get("/user/:user_id", notificationController.GetNotificationsByUserID)
	notificationGroup.Put("/", notificationController.UpdateNotification)
	notificationGroup.Delete("/:notification_id", notificationController.DeleteNotification)
	notificationGroup.Post("/reminders", notificationController.SendReminderNotifications)
}
