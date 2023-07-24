package auth

import (
	"github.com/adeisbright/fiber-user-auth/src/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoute(route fiber.Router, db *gorm.DB) {
	h := &Handler{
		DB: db,
	}

	route.Post("/sign-up", middleware.ValidateCreateUser, h.AddUser)
	route.Post("/login", h.HandleLogin)
	route.Get("/logout", middleware.ValidateToken, h.HandleLogout)

}
