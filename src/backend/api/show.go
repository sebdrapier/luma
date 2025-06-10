package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"elano.fr/src/backend/models"
	"elano.fr/src/backend/storage"
)

func RegisterShowRoutes(app *fiber.App, store storage.ProjectStore) {
	api := app.Group("/api/shows")

	api.Get("/", func(c *fiber.Ctx) error {
		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}
		return c.JSON(project.Shows)
	})

	api.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Show ID is required",
			})
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		for _, s := range project.Shows {
			if s.ID == id {
				return c.JSON(s)
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Show not found",
			"id":    id,
		})
	})

	api.Post("/", func(c *fiber.Ctx) error {
		var input models.Show
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		if input.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Show name is required",
			})
		}

		if len(input.Steps) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "At least one step is required",
			})
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		presetMap := make(map[string]bool)
		for _, p := range project.Presets {
			presetMap[p.ID] = true
		}

		for i, step := range input.Steps {
			if step.PresetID == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":      "Preset ID is required for all steps",
					"step_index": i,
				})
			}
			if !presetMap[step.PresetID] {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":      "Preset not found",
					"preset_id":  step.PresetID,
					"step_index": i,
				})
			}
			if step.DelayMS < 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":      "Delay must be non-negative",
					"step_index": i,
					"delay_ms":   step.DelayMS,
				})
			}
			if step.FadeMS < 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":      "Fade time must be non-negative",
					"step_index": i,
					"fade_ms":    step.FadeMS,
				})
			}
		}

		input.ID = uuid.NewString()
		project.Shows = append(project.Shows, input)

		if err := store.Save(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to save show",
				"details": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(input)
	})

	api.Put("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Show ID is required",
			})
		}

		var input models.Show
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		if input.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Show name is required",
			})
		}

		if len(input.Steps) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "At least one step is required",
			})
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		presetMap := make(map[string]bool)
		for _, p := range project.Presets {
			presetMap[p.ID] = true
		}

		for i, step := range input.Steps {
			if step.PresetID == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":      "Preset ID is required for all steps",
					"step_index": i,
				})
			}
			if !presetMap[step.PresetID] {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":      "Preset not found",
					"preset_id":  step.PresetID,
					"step_index": i,
				})
			}
			if step.DelayMS < 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":      "Delay must be non-negative",
					"step_index": i,
					"delay_ms":   step.DelayMS,
				})
			}
			if step.FadeMS < 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":      "Fade time must be non-negative",
					"step_index": i,
					"fade_ms":    step.FadeMS,
				})
			}
		}

		found := false
		for i, s := range project.Shows {
			if s.ID == id {
				input.ID = id
				project.Shows[i] = input
				found = true
				break
			}
		}

		if !found {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Show not found",
				"id":    id,
			})
		}

		if err := store.Save(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to update show",
				"details": err.Error(),
			})
		}

		return c.JSON(input)
	})

	api.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Show ID is required",
			})
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		shows := make([]models.Show, 0, len(project.Shows))
		found := false

		for _, s := range project.Shows {
			if s.ID != id {
				shows = append(shows, s)
			} else {
				found = true
			}
		}

		if !found {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Show not found",
				"id":    id,
			})
		}

		project.Shows = shows

		if err := store.Save(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to delete show",
				"details": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
