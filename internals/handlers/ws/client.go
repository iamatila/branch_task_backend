package ws

import (
	"fmt"
	"log"

	// "github.com/gofiber/websocket/v2"
	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Category string `json:"category"`
	Username string `json:"username"`
}

type Message struct {
	Content string `json:"content"`
	RoomID  string `json:"roomId"`
	// Category string `json:"category"`
	Username string `json:"username"`
}

// func (c *Client) writeMessage() {
// 	defer func() {
// 		if c.Conn != nil {
// 			c.Conn.Close()
// 		}
// 	}()

// 	for {
// 		message, ok := <-c.Message
// 		if !ok {
// 			return
// 		}
// 		if c.Conn != nil {
// 			if err := c.Conn.WriteJSON(message); err != nil {
// 				log.Printf("error writing JSON message: %v", err)
// 				break
// 			}
// 		} else {
// 			log.Println("writeMessage: Connection is nil")
// 			break
// 		}
// 	}
// }

// func (c *Client) readMessage(hub *Hub) {
// 	defer func() {
// 		if c.Conn != nil {
// 			hub.Unregister <- c
// 			c.Conn.Close()
// 		}
// 	}()

// 	for {
// 		msg := new(Message)
// 		if c.Conn != nil {
// 			if err := c.Conn.ReadJSON(msg); err != nil {
// 				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
// 					break
// 				}
// 				log.Printf("error reading JSON message: %v", err)
// 				break
// 			}
// 			hub.Broadcast <- msg
// 		} else {
// 			log.Println("readMessage: Connection is nil")
// 			break
// 		}
// 	}
// }

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		fmt.Println("write message")
		fmt.Println(message)
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		fmt.Println("read msg")
		fmt.Println(msg)

		hub.Broadcast <- msg
	}
}
