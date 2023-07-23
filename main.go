package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adeisbright/fiber-user-auth/src/config"
	"github.com/adeisbright/fiber-user-auth/src/features/auth"
	"github.com/adeisbright/fiber-user-auth/src/features/user"
	"github.com/adeisbright/fiber-user-auth/src/loaders"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

//Checks that the Service is running fine
func serviceHealthHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "The API is running fine",
		"success": true,
	})
}

//Application Route Setup
func setupRoutes(app *fiber.App) {
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: "2023-07-21",
		TimeZone:   "Africa/Lagos",
	}))

	app.Get("/", serviceHealthHandler)
	app.Post("/cities", AddCity)
	api := app.Group("")
	auth.AuthRoute(api.Group("/auth"))
	auth.RegisterRoute(api.Group("/test"), DB)
	app.Use(auth.ValidateToken)
	app.Get("/users/:id", auth.GetUser)
}

//Database Setup

func GetDB() *gorm.DB {
	return DB
}

type City struct {
	gorm.Model

	ID          uint   `gorm:"primarykey"`
	Title       string `gorm:"index"`
	Description string
}

func setupDB() {

	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUserName := os.Getenv("DB_USERNAME")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUserName, dbPassword, dbHost, dbName)

	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	DB = db
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&City{})

}

func AddCity(c *fiber.Ctx) error {

	type CitySchema struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	body := CitySchema{}

	var city City

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}
	city.Title = body.Title
	city.Description = body.Description

	err = DB.Create(&city).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Could not add city")
	}

	return c.JSON(city)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	password := "secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	setupDB()

	response, error := loaders.ConnectToRedis().Ping().Result()
	if error != nil {
		fmt.Println("Issues with connecting to redis", err)
	}
	fmt.Println(response)
	hostName := config.AppConfig.DBHost
	fmt.Println(hostName, "where is the name")
	app := fiber.New()

	setupRoutes(app)

	log.Fatal(app.Listen("localhost:3000"))
}
