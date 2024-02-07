package wsRoutes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/iamatila/branch_ass_backend/internals/handlers/ws"
)

func SetupWsRoutes(router fiber.Router) {
	// ws
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	go hub.Run()

	ws := router.Group("/ws")

	// Create room
	ws.Post("/createRoom", wsHandler.CreateRoom)
	// Join room
	// ws.Get("/joinRoom/:roomId/:userId/:username", wsHandler.JoinRoom)
	// ws.Get("/joinRoom/:roomId/:userId/:username", websocket.New(wsHandler.JoinRoom))
	ws.Get("/joinRoom/:roomId/:userId/:username", websocket.New(func(c *websocket.Conn) {
		wsHandler.JoinRoom(c, c.Params("roomId"), c.Params("userId"), c.Params("username"))
	}))
	// Get room
	ws.Get("/getRooms", wsHandler.GetRooms)
	// Get client
	ws.Get("/getClients/:roomId", wsHandler.GetClients)
}
