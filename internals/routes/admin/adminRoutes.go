package adminRoutes

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	// jwtware "github.com/gofiber/jwt/v3"

	jwtware "github.com/gofiber/contrib/jwt"
	_ "github.com/golang-jwt/jwt/v5"
	adminHandler "github.com/iamatila/branch_ass_backend/internals/handlers/admin"
	"github.com/joho/godotenv"
)

func SetupAdminRoutes(router fiber.Router) {
	admin := router.Group("/admins")
	// Create a Admin
	admin.Post("/signup", adminHandler.CreateAdmin)
	// Login a Admin
	admin.Post("/login", adminHandler.AdminLogin)
}

func SetupAdmin2Routes(router fiber.Router) {
	err := godotenv.Load("env")
	// err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	secret := os.Getenv("BRANCH_SECRET")

	admin := router.Group("/admin")

	// JWT Middleware
	// admin.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: jwtware.SigningKey{key: []byte(secret)},
	// }))

	// JWT Middleware
	admin.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	}))

	// Read all Admin
	admin.Get("/all", adminHandler.GetAllAdmins)
	// Read one Admin
	admin.Get("/one/:adminid", adminHandler.GetOneOfTheAdmins)
}
