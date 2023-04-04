package handlers

import (
	"errors"
	"template/database"
	"template/models"
	"template/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   errors.New("User validation failed"),
			"details": err.Error(),
		})
	}

	hashedPassword, err := utils.HashingPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Password hashing failed",
		})
	}
	user.Password = hashedPassword

	if err := database.DB.Db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User registration failed",
		})
	}

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var foundUser models.User
	database.DB.Db.Where("email = ?", user.Email).First(&foundUser)

	if foundUser.ID == 0 || !utils.CheckPasswordHash(user.Password, foundUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	token, err := utils.GenerateJWT(foundUser.ID) // Call GenerateJWT function
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating JWT token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
func WhoAmI(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var user models.User
	err := database.DB.Db.First(&user, userID).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user data",
		})
	}

	return c.JSON(user)
}
