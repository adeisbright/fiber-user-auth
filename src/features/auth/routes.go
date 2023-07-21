package features

import (
	"github.com/gofiber/fiber/v2"
)

func AuthRoute(route fiber.Router) {
	route.Post("/sign-up", HandleRegistration)
	route.Post("/login", HandleLgoin)
	route.Post("/users/:id", ProtectedProfileHandler)
}
