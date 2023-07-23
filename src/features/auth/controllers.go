package auth

import (
	"fmt"
	"time"

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

func (h Handler) GetUsers(c *fiber.Ctx) error {
	var users []user.User

	if result := h.DB.Find(&users); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&users)
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

	user.Username = body.Username
	user.Email = body.Email
	user.Password = body.Password

	if result := h.DB.Create(&user); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&user)
}

func (h Handler) CheckLogin(c *fiber.Ctx) error {
	var validUser user.User
	body := UserSchema{}
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"success": false,
		})
	}

	validUser.Username = body.Username
	validUser.Password = body.Password

	var foundUser user.User
	if err := h.DB.Where("username = ?", validUser.Username).First(&foundUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if foundUser.Password != validUser.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := GenerateJWTToken(foundUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate JWT token"})
	}

	expirationTime := time.Now().Add(time.Hour * 24).Unix()
	redisKey := "user:" + fmt.Sprint(foundUser.ID)
	err = loaders.ConnectToRedis().Set(redisKey, validUser.Username, time.Duration(expirationTime)).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"token":   "",
		})
	}

	return c.JSON(fiber.Map{"token": token})
}

func ValidateToken(c *fiber.Ctx) error {
	// Extract the token from the Authorization header
	authHeader := c.Get("Authorization")
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	tokenString := authHeader[7:]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID := uint(claims["user_id"].(float64))

	c.Locals("userId", userID)
	return c.Next()

}

func GetUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	loggedInUserId := c.Locals("userId")
	if userId != loggedInUserId {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "You cannot view this users profile",
		})
	}

	redisKey := "user:" + userId
	data := loaders.ConnectToRedis().Get(redisKey)
	fmt.Println(data)

	if data == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"token":   "",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Welcome to your profile",
		"data":    userId,
		"success": true,
	})
}
