package user

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {

	loggedInUserId := c.Locals("userId").(uint)

	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid parameter value, must be an integer",
			"success": false,
		})
	}
	unsignedUserId := uint(userId)

	if unsignedUserId != loggedInUserId {
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
