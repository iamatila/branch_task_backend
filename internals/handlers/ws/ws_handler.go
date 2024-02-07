package ws

import (
	// "fmt"
	// "math/rand"

	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type CreateRoomReq struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Name     string `json:"name"`
}

func (h *Handler) CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// rand.Seed(time.Now().UnixNano())

	// num := rand.Intn(100) * 100
	// num += rand.Intn(100) * 10
	// // num += rand.Intn(100) * 1

	// var NewRoomID = fmt.Sprintf("%s%s%d", req.ID, "_", num)
	// var NewNameID = fmt.Sprintf("%s%s%s%s%d", req.ID, "_", req.Category, "_", num)

	h.hub.Rooms[req.ID] = &Room{
		// ID: NewRoomID,
		ID:       req.ID,
		Category: req.Category,
		// Name:     NewNameID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	// return c.JSON(h.hub.Rooms[req.ID])
	return c.JSON(req)
}

// func (h *Handler) JoinRoom(c *fiber.Ctx) error {
// 	// Your logic for handling WebSocket connection
// 	roomID := c.Params("roomId")
// 	clientID := c.Params("userId")
// 	username := c.Params("username")

// 	// Upgrade the WebSocket connection
// 	websocket.New(func(conn *websocket.Conn) {
// 		cl := &Client{
// 			Conn:     conn,
// 			Message:  make(chan *Message, 10),
// 			ID:       clientID,
// 			RoomID:   roomID,
// 			Username: username,
// 		}

// 		// Register a new client through the register channel
// 		h.hub.Register <- cl

// 		m := &Message{
// 			Content:  username + " is Ready to chat",
// 			RoomID:   roomID,
// 			Username: username,
// 		}

// 		// Broadcast that message
// 		h.hub.Broadcast <- m

// 		go cl.writeMessage()
// 		go cl.readMessage(h.hub)
// 	})

// 	return nil
// }

func (h *Handler) JoinRoom(conn *websocket.Conn, roomId, userId, username string) {
	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       userId,
		RoomID:   roomId,
		Username: username,
	}

	// Register a new client through the register channel
	h.hub.Register <- cl
	fmt.Println("cl message 1")
	fmt.Println(cl)
	fmt.Println("cl message 2")

	m := &Message{
		Content: "A new user has joined the room",
		// Content:  username + " is Ready to chat",
		RoomID:   roomId,
		Username: username,
	}

	// Broadcast that message
	h.hub.Broadcast <- m
	fmt.Println("m message 1")
	fmt.Println(m)
	fmt.Println("m message 2")

	go cl.writeMessage()
	fmt.Println("cl write message 1")
	fmt.Println(cl)
	fmt.Println("cl write message 2")
	go cl.readMessage(h.hub)
	fmt.Println("cl read message 1")
	fmt.Println(cl)
	fmt.Println("cl read message 2")

	// Keep the connection open
	for {
		// You can replace this with more useful logic
		// For example, you could listen for incoming messages here
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

type RoomRes struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Name     string `json:"name"`
}

func (h *Handler) GetRooms(c *fiber.Ctx) error {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:       r.ID,
			Category: r.Category,
			Name:     r.Name,
		})
	}

	return c.JSON(rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *fiber.Ctx) error {
	var clients []ClientRes
	roomId := c.Params("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		return c.JSON(clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	return c.JSON(clients)
}
