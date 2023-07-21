package auth

import (
	"fmt"
	"time"

	"github.com/adeisbright/fiber-user-auth/src/features/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var jwtSecret = []byte("example")
var db *gorm.DB

func GenerateJWTToken(userID uint) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	return token.SignedString(jwtSecret)
}

type UserSchema struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleRegistration(c *fiber.Ctx) error {
	var validUser user.User

	body := UserSchema{}

	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"success": false,
		})
	}

	validUser.Email = body.Email
	validUser.Username = body.Username
	validUser.Password = body.Password

	err = db.Create(user.User{
		Email:    "Adebayo",
		Username: "Someto",
		Password: "123456",
	}).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	fmt.Println("Everything fine here")
	return c.JSON(validUser)
}

func HandleLgoin(c *fiber.Ctx) error {
	var validUser user.User
	err := c.BodyParser(&validUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"success": false,
		})
	}

	var foundUser user.User
	if err := db.Where("username = ?", validUser.Username).First(&foundUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if foundUser.Password != validUser.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := GenerateJWTToken(foundUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate JWT token"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func AuthenticateRequest(c *fiber.Ctx) error {
	return c.SendString("Authentication: Not Implemented Yet")
}

func ProtectedProfileHandler(c *fiber.Ctx) error {
	return c.SendString("Profile: Not Implemented Yet")
}
