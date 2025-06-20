package ws

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

func handleBroadcast() {
	for msg := range broadcast {
		var conns []*websocket.Conn
		clients.Range(func(key, _ interface{}) bool {
			if c, ok := key.(*websocket.Conn); ok {
				conns = append(conns, c)
			}
			return true
		})
		for _, c := range conns {
			go func(c *websocket.Conn) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic in broadcast: %v", r)
						clients.Delete(c)
					}
				}()
				if err := writeJSON(c, msg); err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					clients.Delete(c)
				}
			}(c)
		}
	}
}
