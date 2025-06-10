package ws

import (
	"encoding/json"
	"log"
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
	c.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err := c.WriteJSON(response); err != nil {
		log.Printf("Error sending error message: %v", err)
	}
}

func getClientCount() int {
	count := 0
	clients.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}
