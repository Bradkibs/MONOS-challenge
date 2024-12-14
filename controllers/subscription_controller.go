package controllers

import (
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/services"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

type SubscriptionController struct {
	DB *pgxpool.Pool
}

func NewSubscriptionController(db *pgxpool.Pool) *SubscriptionController {
	return &SubscriptionController{DB: db}
}

func (sc *SubscriptionController) CreateSubscription(c *fiber.Ctx) error {
	var req models.Subscription
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Set default values
	req.ID = utils.GenerateUniqueID()
	req.Status = "active"
	req.StartDate = time.Now()

	if req.EndDate == nil {
		endDate := req.StartDate.AddDate(0, 1, 0)
		req.EndDate = &endDate
	}

	err := services.CreateSubscription(&req, sc.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(req)
}

func (sc *SubscriptionController) GetSubscription(c *fiber.Ctx) error {
	subscriptionID := c.Params("subscription_id")
	_, err := uuid.Parse(subscriptionID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subscription ID"})
	}

	subscription, err := services.GetSubscription(subscriptionID, sc.DB)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(subscription)
}

func (sc *SubscriptionController) UpdateSubscription(c *fiber.Ctx) error {
	subscriptionID := c.Params("subscription_id")
	_, err := uuid.Parse(subscriptionID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subscription ID"})
	}

	var req models.Subscription
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	req.ID = uuid.MustParse(subscriptionID)
	err = services.UpdateSubscription(&req, sc.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(req)
}

func (sc *SubscriptionController) CancelSubscription(c *fiber.Ctx) error {
	subscriptionID := c.Params("subscription_id")
	_, err := uuid.Parse(subscriptionID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subscription ID"})
	}

	err = services.CancelSubscription(subscriptionID, sc.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Subscription canceled successfully"})
}

func (sc *SubscriptionController) DowngradeSubscription(c *fiber.Ctx) error {
	subscriptionID := c.Params("subscription_id")
	_, err := uuid.Parse(subscriptionID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subscription ID"})
	}

	var request struct {
		NewTier string `json:"new_tier"`
	}

	if err := c.BodyParser(&request); err != nil || request.NewTier == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "New tier is required"})
	}

	err = services.DowngradeSubscription(subscriptionID, request.NewTier, sc.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Subscription downgraded successfully"})
}

func (sc *SubscriptionController) DeleteSubscription(c *fiber.Ctx) error {
	subscriptionID := c.Params("subscription_id")
	_, err := uuid.Parse(subscriptionID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid subscription ID"})
	}

	err = services.DeleteSubscription(subscriptionID, sc.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Subscription deleted successfully"})
}
