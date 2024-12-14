package controllers

import (
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/services"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BusinessController struct {
	DB *pgxpool.Pool
}

func (bc *BusinessController) GetAllBusinesses(c *fiber.Ctx) error {
	businesses, err := services.GetAllBusinesses(bc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(businesses)
}

func (bc *BusinessController) CreateBusiness(c *fiber.Ctx) error {
	var business models.Business
	if err := c.BodyParser(&business); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	business.ID = utils.GenerateUniqueID()
	err := services.CreateBusiness(&business, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(business)
}

func (bc *BusinessController) GetBusinessByID(c *fiber.Ctx) error {
	businessID := c.Params("business_id")
	id, err := uuid.Parse(businessID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid business ID"})
	}

	business, err := services.GetBusinessByID(id, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Business not found"})
	}

	return c.Status(fiber.StatusOK).JSON(business)
}

func (bc *BusinessController) UpdateBusiness(c *fiber.Ctx) error {
	var business models.Business
	if err := c.BodyParser(&business); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.UpdateBusiness(&business, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(business)
}

func (bc *BusinessController) DeleteBusiness(c *fiber.Ctx) error {
	businessID := c.Params("business_id")
	id, err := uuid.Parse(businessID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid business ID"})
	}

	err = services.DeleteBusiness(id, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Business not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Business deleted successfully"})
}

func (bc *BusinessController) GetBusinessesByVendorID(c *fiber.Ctx) error {
	vendorID := c.Params("vendor_id")
	id, err := uuid.Parse(vendorID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid vendor ID"})
	}

	businesses, err := services.GetBusinessesByVendorID(id, bc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(businesses)
}
