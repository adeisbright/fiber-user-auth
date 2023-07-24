package middleware

import (
	"github.com/adeisbright/fiber-user-auth/src/loaders"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtSecret = []byte("example")

func ValidateToken(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing Authorization Header",
			"success": false,
		})
	}

	tokenString := authHeader[7:]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Authorization Token",
			"success": false,
		})
	}

	blacklistedTokenKey := "token:blacklist:" + tokenString
	data, _ := loaders.ConnectToRedis().Get(blacklistedTokenKey).Result()
	if len(data) > 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Session Expired. Please Login Again",
			"success": false,
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Authorization Token",
			"success": false,
		})
	}

	userID := uint(claims["user_id"].(float64))

	c.Locals("userId", userID)
	return c.Next()

}
