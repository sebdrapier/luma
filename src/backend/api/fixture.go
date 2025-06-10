package api

import (
	"elano.fr/src/backend/models"
	"elano.fr/src/backend/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterFixtureRoutes(app *fiber.App, store storage.ProjectStore) {
	r := app.Group("/api/fixtures")

	r.Get("/", func(c *fiber.Ctx) error {
		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}
		return c.JSON(project.Fixtures)
	})

	r.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Fixture ID is required",
			})
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		for _, f := range project.Fixtures {
			if f.ID == id {
				return c.JSON(f)
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Fixture not found",
			"id":    id,
		})
	})

	r.Post("/", func(c *fiber.Ctx) error {
		var fixture models.Fixture
		if err := c.BodyParser(&fixture); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		if fixture.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Fixture name is required",
			})
		}
		if fixture.Type == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Fixture type is required",
			})
		}
		if len(fixture.Channels) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "At least one channel is required",
			})
		}

		for i, ch := range fixture.Channels {
			if ch.Name == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":         "Channel name is required",
					"channel_index": i,
				})
			}
			if ch.Min < 0 || ch.Max > 255 || ch.Min > ch.Max {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Invalid channel range (must be 0-255, min <= max)",
					"channel": ch.Name,
				})
			}
			if ch.ChannelAddress < 1 || ch.ChannelAddress > 512 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Channel address must be between 1 and 512",
					"channel": ch.Name,
				})
			}
		}

		fixture.ID = uuid.New().String()

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		project.Fixtures = append(project.Fixtures, fixture)

		if err := store.Save(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to save fixture",
				"details": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fixture)
	})

	r.Put("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Fixture ID is required",
			})
		}

		var update models.Fixture
		if err := c.BodyParser(&update); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		if update.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Fixture name is required",
			})
		}
		if update.Type == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Fixture type is required",
			})
		}
		if len(update.Channels) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "At least one channel is required",
			})
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		found := false
		for i, f := range project.Fixtures {
			if f.ID == id {
				update.ID = id
				project.Fixtures[i] = update
				found = true
				break
			}
		}

		if !found {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Fixture not found",
				"id":    id,
			})
		}

		if err := store.Save(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to update fixture",
				"details": err.Error(),
			})
		}

		return c.JSON(update)
	})

	r.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Fixture ID is required",
			})
		}

		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load project",
			})
		}

		fixtures := make([]models.Fixture, 0, len(project.Fixtures))
		found := false

		for _, f := range project.Fixtures {
			if f.ID != id {
				fixtures = append(fixtures, f)
			} else {
				found = true
			}
		}

		if !found {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Fixture not found",
				"id":    id,
			})
		}

		project.Fixtures = fixtures

		if err := store.Save(project); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to delete fixture",
				"details": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
