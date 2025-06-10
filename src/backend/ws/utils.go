package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func mustMarshal(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}

func sendError(c *websocket.Conn, errType, message, details string) {
	response := Message{
		Type:    "error",
		Payload: mustMarshal(ErrorResponse{Error: message, Details: details}),
	}
	if err := writeJSON(c, response); err != nil {
		log.Printf("Error sending error message: %v", err)
	}
}
func getConnMutex(c *websocket.Conn) *sync.Mutex {
	if m, ok := clients.Load(c); ok {
		if mu, ok := m.(*sync.Mutex); ok {
			return mu
		}
	}
	return nil
}

func writeJSON(c *websocket.Conn, v interface{}) error {
	mu := getConnMutex(c)
	if mu != nil {
		mu.Lock()
		defer mu.Unlock()
	}
	c.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return c.WriteJSON(v)
}

func writeMessage(c *websocket.Conn, messageType int, data []byte) error {
	mu := getConnMutex(c)
	if mu != nil {
		mu.Lock()
		defer mu.Unlock()
	}
	c.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return c.WriteMessage(messageType, data)
}

func getClientCount() int {
	count := 0
	clients.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}
