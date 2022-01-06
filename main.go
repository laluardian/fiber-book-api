package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/laluardian/fiber-book-api/config"
	"github.com/laluardian/fiber-book-api/routes"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	setupRoutes(app)

	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

func setupRoutes(app *fiber.App) {
	// give response when at "/"
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the '/' endpoint 😉",
		})
	})

	// api group
	api := app.Group("/api")

	// give response when at "/api"
	api.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the '/api' endpoint 😉",
		})
	})

	// connect book routes
	routes.BookRoute(api.Group("/books"))
}
