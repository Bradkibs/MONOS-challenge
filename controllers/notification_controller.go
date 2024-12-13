package controllers

import (
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationController struct {
	Pool *pgxpool.Pool
}

func (c *NotificationController) CreateNotification(ctx *fiber.Ctx) error {
	var notification models.Notification
	if err := ctx.BodyParser(&notification); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := services.CreateNotification(c.Pool, &notification); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create notification",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(notification)
}

func (c *NotificationController) GetNotificationByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	notification, err := services.GetNotificationByID(c.Pool, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notification not found",
		})
	}

	return ctx.JSON(notification)
}

func (c *NotificationController) GetNotificationsByUserID(ctx *fiber.Ctx) error {
	userIdStr := ctx.Params("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	notification, err := services.GetNotificationByUserID(c.Pool, userId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notifications not found",
		})
	}

	return ctx.JSON(notification)
}

func (c *NotificationController) UpdateNotification(ctx *fiber.Ctx) error {
	var notification models.Notification
	if err := ctx.BodyParser(&notification); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := services.UpdateNotification(c.Pool, &notification); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update notification",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(notification)
}

func (c *NotificationController) DeleteNotification(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := services.DeleteNotification(c.Pool, id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete notification",
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *NotificationController) SendReminderNotifications(ctx *fiber.Ctx) error {
	if err := services.SendReminderNotification(c.Pool); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send reminders",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reminders sent successfully",
	})
}
