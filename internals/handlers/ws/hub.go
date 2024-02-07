package ws

import (
	"sync"
	// "github.com/gofiber/fiber/v2"
	// "github.com/gofiber/websocket/v2"
	// "github.com/gofiber/contrib/websocket"
)

type Room struct {
	ID       string             `json:"id"`
	Category string             `json:"category"`
	Name     string             `json:"name"`
	Clients  map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

// func (h *Hub) Run() {
// 	for {
// 		select {
// 		case cl := <-h.Register:
// 			h.mutex.Lock()
// 			if room, ok := h.Rooms[cl.RoomID]; ok {
// 				if _, ok := room.Clients[cl.ID]; !ok {
// 					room.Clients[cl.ID] = cl
// 				}
// 			} else {
// 				h.Rooms[cl.RoomID] = &Room{
// 					ID:      cl.RoomID,
// 					Name:    "Room " + cl.RoomID,
// 					Clients: map[string]*Client{cl.ID: cl},
// 				}
// 			}
// 			h.mutex.Unlock()

// 		case cl := <-h.Unregister:
// 			h.mutex.Lock()
// 			if room, ok := h.Rooms[cl.RoomID]; ok {
// 				if _, ok := room.Clients[cl.ID]; ok {
// 					if len(room.Clients) != 0 {
// 						h.Broadcast <- &Message{
// 							Content:  "user left the chat",
// 							RoomID:   cl.RoomID,
// 							Username: cl.Username,
// 						}
// 					}

// 					delete(room.Clients, cl.ID)
// 					close(cl.Message)
// 				}

// 				if len(room.Clients) == 0 {
// 					delete(h.Rooms, cl.RoomID)
// 				}
// 			}
// 			h.mutex.Unlock()

// 		case m := <-h.Broadcast:
// 			h.mutex.RLock()
// 			if room, ok := h.Rooms[m.RoomID]; ok {
// 				for _, cl := range room.Clients {
// 					cl.Message <- m
// 				}
// 			}
// 			h.mutex.RUnlock()
// 		}
// 	}
// }

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.RoomID]; !ok {
				h.Rooms[client.RoomID] = &Room{
					ID:       client.RoomID,
					Category: client.Category,
					Name:     client.Username, // You might want to replace this with the actual room name
					Clients:  make(map[string]*Client),
				}
			}
			h.Rooms[client.RoomID].Clients[client.ID] = client
		case client := <-h.Unregister:
			if room, ok := h.Rooms[client.RoomID]; ok {
				delete(room.Clients, client.ID)
			}
		case message := <-h.Broadcast:
			if room, ok := h.Rooms[message.RoomID]; ok {
				for _, client := range room.Clients {
					select {
					case client.Message <- message:
					default:
						close(client.Message)
						delete(room.Clients, client.ID)
					}
				}
			}
		}
	}
}

// func UpgradeToWebSocket(handler func(*websocket.Conn)) fiber.Handler {
// 	return websocket.New(handler)
// }
