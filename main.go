package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adeisbright/fiber-user-auth/src/common"
	"github.com/adeisbright/fiber-user-auth/src/features/auth"
	"github.com/adeisbright/fiber-user-auth/src/features/user"
	"github.com/adeisbright/fiber-user-auth/src/loaders"
	"github.com/adeisbright/fiber-user-auth/src/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

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

	api := app.Group("")
	auth.AuthRoute(api.Group("/auth"), DB)
	app.Use(middleware.ValidateToken)

	user.UserRoute(api.Group("/users"))
}

//Database Setup

func GetDB() *gorm.DB {
	return DB
}

func setupDB() {

	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUserName := os.Getenv("DB_USERNAME")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUserName, dbPassword, dbHost, dbName)

	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{})
	if err != nil {
		common.CustomLog(err.Error(), "Error")
		panic("Failed to connect to database")
	}
	DB = db

	db.AutoMigrate(&user.User{})

}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		common.CustomLog(err.Error(), "Error")
		log.Fatal("Error loading .env file")
	}
	setupDB()

	_, error := loaders.ConnectToRedis().Ping().Result()

	if error != nil {
		common.CustomLog(error.Error(), "Error")
		fmt.Println("Issues with connecting to redis", err)
	}

	app := fiber.New()

	setupRoutes(app)

	log.Fatal(app.Listen("localhost:3000"))
}
