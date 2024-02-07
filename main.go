package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iamatila/branch_ass_backend/database"
	"github.com/iamatila/branch_ass_backend/router"
	"github.com/joho/godotenv"
)

func main() {
	// Start a new fiber app
	app := fiber.New()

	// Initialize default config
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Connect to the Database
	database.ConnectDB()

	// Setup the router
	router.SetupRoutes(app)

	err := godotenv.Load("env")
	// err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")

	// Listen on PORT 3000
	// app.Listen(port)
	app.Listen(`0.0.0.0:` + port)
}
