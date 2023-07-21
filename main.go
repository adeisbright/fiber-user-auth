package main

import "github.com/gofiber/fiber/v2"

func serviceHealthHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "The API is running fine",
		"success": true,
	})
}

func main() {
	app := fiber.New()

	app.Get("/", serviceHealthHandler)
	app.Listen("localhost:3000")
}
