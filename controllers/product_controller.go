package controllers

import (
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/services"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductController struct {
	DB *pgxpool.Pool
}

func NewProductController(db *pgxpool.Pool) *ProductController {
	return &ProductController{DB: db}
}

func (pc *ProductController) AddProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Generate a UUID if none provided
	if product.ID == uuid.Nil {
		product.ID = utils.GenerateUniqueID()
	}

	if err := services.AddProduct(product, pc.DB); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Product added successfully"})
}

func (pc *ProductController) GetProducts(c *fiber.Ctx) error {
	businessIDStr := c.Query("business_id")
	if businessIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing business_id query parameter"})
	}

	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid business_id"})
	}

	products, err := services.GetProductsByBusinessID(businessID.String(), pc.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(products)
}

func (pc *ProductController) UpdateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if product.ID == uuid.Nil || product.BusinessID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product ID and Business ID are required"})
	}

	if err := services.UpdateProduct(product, pc.DB); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

func (pc *ProductController) DeleteProduct(c *fiber.Ctx) error {
	productIDStr := c.Query("product_id")
	businessIDStr := c.Query("business_id")

	if productIDStr == "" || businessIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing product_id or business_id query parameters"})
	}

	productID, err1 := uuid.Parse(productIDStr)
	businessID, err2 := uuid.Parse(businessIDStr)

	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product_id or business_id"})
	}

	if err := services.DeleteProduct(productID.String(), businessID.String(), pc.DB); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}
