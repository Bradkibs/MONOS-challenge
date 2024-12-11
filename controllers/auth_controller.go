package controllers

import (
	"github.com/Bradkibs/MONOS-challenge/services"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthController struct {
	DB *pgxpool.Pool
}

func (ac *AuthController) RegisterByMail(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := services.RegisterUserByEmail(input.Email, input.Password, input.Role, ac.DB)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token})
}

func (ac *AuthController) RegisterByPhoneNumber(c *fiber.Ctx) error {
	var input struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
		Role        string `json:"role"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := services.RegisterUserByPhoneNumber(input.PhoneNumber, input.Password, input.Role, ac.DB)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token})
}

func (ac *AuthController) LoginByMail(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := services.LoginUser(input.Email, input.Password, ac.DB)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (ac *AuthController) LoginByPhoneNumber(c *fiber.Ctx) error {
	var input struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := services.LoginUser(input.PhoneNumber, input.Password, ac.DB)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (ac *AuthController) ValidateToken(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
	}

	claims, err := services.ParseJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"claims": claims})
}
