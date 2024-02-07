package router

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/swagger" // swagger middleware

	"github.com/gofiber/fiber/v2/middleware/logger"
	// _ "github.com/iamatila/branch_ass_backend/docs"

	adminRoutes "github.com/iamatila/branch_ass_backend/internals/routes/admin"
	userRoutes "github.com/iamatila/branch_ass_backend/internals/routes/user"
	wsRoutes "github.com/iamatila/branch_ass_backend/internals/routes/ws"
	// "github.com/iamatila/branch_ass_backend/internals/handlers/ws"
)

func SetupRoutes(app *fiber.App) {
	// app.Get("/swagger/*", swagger.HandlerDefault)
	// Group api calls with param '/api'
	api := app.Group("/api/v1", logger.New())

	// // ws
	// hub := ws.NewHub()
	// wsHandler := ws.NewHandler(hub)
	// go hub.Run()

	// Setup note routes, can use same syntax to add routes for more models
	userRoutes.SetupUserRoutes(api)
	userRoutes.SetupUser2Routes(api)
	adminRoutes.SetupAdminRoutes(api)
	adminRoutes.SetupAdmin2Routes(api)
	wsRoutes.SetupWsRoutes(api)
}
