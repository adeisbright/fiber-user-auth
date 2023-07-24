package auth

import (
	"fmt"
	"time"

	"github.com/adeisbright/fiber-user-auth/src/common"
	"github.com/adeisbright/fiber-user-auth/src/features/user"
	"github.com/adeisbright/fiber-user-auth/src/loaders"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var jwtSecret = []byte("example")

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

type Handler struct {
	DB *gorm.DB
}

func (h Handler) AddUser(c *fiber.Ctx) error {
	var user user.User
	body := UserSchema{}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"success": false,
		})
	}

	hashedPassword, error := common.HashPassword(body.Password)
	if error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	fmt.Println(hashedPassword)
	user.Username = body.Username
	user.Email = body.Email
	user.Password = hashedPassword

	if result := h.DB.Create(&user); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&user)
}

func (h Handler) HandleLogin(c *fiber.Ctx) error {
	var validUser user.User
	body := UserSchema{}
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"success": false,
		})
	}

	validUser.Email = body.Email
	validUser.Password = body.Password

	var foundUser user.User
	if err := h.DB.Where("email = ?", validUser.Email).First(&foundUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Invalid Login Credentials. Try again",
			"success": false,
		})
	}

	if !common.CheckPasswordHash(body.Password, foundUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials. Try again",
			"success": false,
		})
	}

	token, err := GenerateJWTToken(foundUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   token,
		"success": true,
		"data":    foundUser,
	})
}

func (h Handler) HandleLogout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenString := authHeader[7:]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Bad Token",
			"success": false,
		})
	}

	blacklistedTokenKey := "token:blacklist:" + tokenString
	err = loaders.ConnectToRedis().Set(blacklistedTokenKey, tokenString, 0).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
		"success": true,
	})
}
