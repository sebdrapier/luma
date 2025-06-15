package api

import (
	"elano.fr/src/backend/dmx"
	"github.com/gofiber/fiber/v2"
)

func RegisterUSBRoutes(app *fiber.App) {
	r := app.Group("/api/usb")
	r.Get("/interfaces", GetUSBInterfaces)
}

func GetUSBInterfaces(c *fiber.Ctx) error {
	ports, err := dmx.ListDMXPorts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(ports)
}
