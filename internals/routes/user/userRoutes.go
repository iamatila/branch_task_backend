package userRoutes

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	// jwtware "github.com/gofiber/jwt/v3"
	jwtware "github.com/gofiber/contrib/jwt"
	_ "github.com/golang-jwt/jwt/v5"
	userHandler "github.com/iamatila/branch_ass_backend/internals/handlers/user"
	"github.com/joho/godotenv"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/users")
	// Create a User
	user.Post("/signup", userHandler.CreateUser)
	// Login a User
	user.Post("/login", userHandler.Login)
}

func SetupUser2Routes(router fiber.Router) {
	err := godotenv.Load("env")
	// err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	secret := os.Getenv("BRANCH_SECRET")

	user := router.Group("/user")

	// JWT Middleware
	// user.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: jwtware.SigningKey{key: []byte(secret)},
	// }))

	// JWT Middleware
	user.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	}))

	// Read all Users
	user.Get("/all", userHandler.GetAllUsers)
	// Read one User
	user.Get("/one/:userid", userHandler.GetOneUser)
}
