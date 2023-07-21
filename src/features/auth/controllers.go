package auth

import "github.com/gofiber/fiber/v2"

func HandleRegistration(c *fiber.Ctx) error {
	return c.SendString("Registration: Not Implemented Yet")
}

func HandleLgoin(c *fiber.Ctx) error {
	return c.SendString("Login: Not Implemented Yet")
}

func AuthenticateRequest(c *fiber.Ctx) error {
	return c.SendString("Authentication: Not Implemented Yet")
}

func ProtectedProfileHandler(c *fiber.Ctx) error {
	return c.SendString("Profile: Not Implemented Yet")
}
