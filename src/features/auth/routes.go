package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoute(route fiber.Router) {
	route.Post("/sign-up", HandleRegistration)
	route.Post("/login", HandleLgoin)
	route.Post("/users/:id", ProtectedProfileHandler)

}

func RegisterRoute(route fiber.Router, db *gorm.DB) {
	h := &Handler{
		DB: db,
	}
	route.Get("/users", h.GetUsers)
	route.Post("/users", h.AddUser)
	route.Post("/login", h.CheckLogin)
}
