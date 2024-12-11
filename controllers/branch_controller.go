package controllers

import (
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/services"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BranchController struct {
	DB *pgxpool.Pool
}

func (bc *BranchController) AddBranch(c *fiber.Ctx) error {
	var branch models.Branch
	if err := c.BodyParser(&branch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	branch.ID = utils.GenerateUniqueID()
	err := services.AddBranch(&branch, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(branch)
}

func (bc *BranchController) GetBranches(c *fiber.Ctx) error {
	businessID := c.Query("business_id")
	if businessID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "business_id is required"})
	}

	branches, err := services.GetBranchesByBusinessID(businessID, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(branches)
}

func (bc *BranchController) UpdateBranch(c *fiber.Ctx) error {
	var branch models.Branch
	if err := c.BodyParser(&branch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.UpdateBranch(&branch, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(branch)
}

func (bc *BranchController) DeleteBranch(c *fiber.Ctx) error {
	branchID := c.Query("branch_id")
	businessID := c.Query("business_id")
	if branchID == "" || businessID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "branch_id and business_id are required"})
	}

	err := services.DeleteBranch(branchID, businessID, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Branch deleted successfully"})
}

func (bc *BranchController) UpdateBranchesForSubscription(c *fiber.Ctx) error {
	var req struct {
		SubscriptionID string   `json:"subscription_id"`
		BranchChange   int      `json:"branch_change"`
		BranchNames    []string `json:"branch_names"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.UpdateBranchesForSubscription(req.SubscriptionID, req.BranchChange, req.BranchNames, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Branches updated successfully"})
}
