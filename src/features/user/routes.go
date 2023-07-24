package user

import "github.com/gofiber/fiber/v2"

func UserRoute(route fiber.Router) {

	route.Get("/:id", GetProfile)
}
