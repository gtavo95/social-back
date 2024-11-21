package main

import (
	"log"
	"social/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	app := fiber.New(fiber.Config{DisablePreParseMultipartForm: true, StreamRequestBody: true})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	app.Post("/text", routes.SocialPostText)

	app.Listen(":8080")
}
