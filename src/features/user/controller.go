package user

import (
	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {
	userId := c.Params("id")
	loggedInUserId := c.Locals("userId")

	if userId != loggedInUserId {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "You cannot view this users profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Welcome to your profile",
		"data":    userId,
		"success": true,
	})
}
