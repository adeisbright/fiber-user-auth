package main

import (
	features "github.com/adeisbright/fiber-user-auth/src/features/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func serviceHealthHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "The API is running fine",
		"success": true,
	})
}

func setupRoutes(app *fiber.App) {
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: "2023-07-21",
		TimeZone:   "Africa/Lagos",
	}))

	app.Get("/", serviceHealthHandler)

	api := app.Group("")
	features.AuthRoute(api.Group("/auth"))
}

func main() {
	app := fiber.New()

	setupRoutes(app)

	app.Listen("localhost:3000")
}
