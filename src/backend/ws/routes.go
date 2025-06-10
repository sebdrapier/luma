package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupWebSocketRoutes(app *fiber.App) {
	app.Get("/ws/control", websocket.New(handleWebSocket))
	go handleBroadcast()
}
