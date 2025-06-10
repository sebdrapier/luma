package api

import (
	"elano.fr/src/backend/models"
	"elano.fr/src/backend/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterPresetRoutes(app *fiber.App, store storage.ProjectStore) {
	api := app.Group("/api/presets")

	api.Get("/", func(c *fiber.Ctx) error {
		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}
		return c.JSON(project.Presets)
	})

	api.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Preset ID is required",
			})
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		for _, p := range project.Presets {
			if p.ID == id {
				return c.JSON(p)
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Preset not found",
			"id":    id,
		})
	})

	api.Post("/", func(c *fiber.Ctx) error {
		var input models.Preset
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		if input.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Preset name is required",
			})
		}

		if len(input.Channels) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "At least one channel value is required",
			})
		}

		for i, ch := range input.Channels {
			if ch.DMXAddress < 1 || ch.DMXAddress > 512 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":         "DMX address must be between 1 and 512",
					"channel_index": i,
					"dmx_address":   ch.DMXAddress,
				})
			}

			if ch.Value > 255 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":         "Channel value must be between 0 and 255",
					"channel_index": i,
					"value":         ch.Value,
				})
			}
		}

		input.ID = uuid.NewString()

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		project.Presets = append(project.Presets, input)

		if err := store.Save(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to save preset",
				"details": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(input)
	})

	api.Put("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Preset ID is required",
			})
		}

		var input models.Preset
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		if input.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Preset name is required",
			})
		}

		if len(input.Channels) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "At least one channel value is required",
			})
		}

		for i, ch := range input.Channels {
			if ch.DMXAddress < 1 || ch.DMXAddress > 512 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":         "DMX address must be between 1 and 512",
					"channel_index": i,
					"dmx_address":   ch.DMXAddress,
				})
			}
			if ch.Value > 255 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":         "Channel value must be between 0 and 255",
					"channel_index": i,
					"value":         ch.Value,
				})
			}
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		found := false
		for i, p := range project.Presets {
			if p.ID == id {
				input.ID = id
				project.Presets[i] = input
				found = true
				break
			}
		}

		if !found {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Preset not found",
				"id":    id,
			})
		}

		if err := store.Save(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to update preset",
				"details": err.Error(),
			})
		}

		return c.JSON(input)
	})

	api.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Preset ID is required",
			})
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		presets := make([]models.Preset, 0, len(project.Presets))
		found := false

		for _, p := range project.Presets {
			if p.ID != id {
				presets = append(presets, p)
			} else {
				found = true
			}
		}

		if !found {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Preset not found",
				"id":    id,
			})
		}

		project.Presets = presets

		if err := store.Save(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to delete preset",
				"details": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
