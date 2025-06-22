package api

import (
	"elano.fr/src/backend/models"
	"elano.fr/src/backend/storage"
	"elano.fr/src/backend/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterProjectRoutes(app *fiber.App, store storage.ProjectStore, enableDMX bool) {
	r := app.Group("/api/projects")

	r.Get("/", func(c *fiber.Ctx) error {
		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No project found",
			})
		}
		return c.JSON(project)
	})

	r.Post("/", func(c *fiber.Ctx) error {
		var proj models.Project
		if err := c.BodyParser(&proj); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		if proj.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Project name is required",
			})
		}

		if proj.USBInterface == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "USB interface is required",
			})
		}

		isNew := false
		if proj.ID == "" {
			proj.ID = uuid.New().String()
			isNew = true
		}

		if proj.Fixtures == nil {
			proj.Fixtures = []models.Fixture{}
		}
		if proj.Presets == nil {
			proj.Presets = []models.Preset{}
		}
		if proj.Shows == nil {
			proj.Shows = []models.Show{}
		}

		if err := store.Save(&proj); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to save project",
				"details": err.Error(),
			})
		}

		if enableDMX {
			if err := ws.InitializeDMXController(proj.USBInterface); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Failed to initialize DMX controller",
					"details": err.Error(),
				})
			}
		}

		status := fiber.StatusOK
		if isNew {
			status = fiber.StatusCreated
		}
		return c.Status(status).JSON(proj)
	})

	r.Put("/", func(c *fiber.Ctx) error {
		var proj models.Project
		if err := c.BodyParser(&proj); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
		}

		if proj.ID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Project ID is required",
			})
		}

		if proj.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Project name is required",
			})
		}

		if proj.USBInterface == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "USB interface is required",
			})
		}

		currentProject := store.Get()
		if currentProject == nil || currentProject.ID != proj.ID {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Project not found or ID mismatch",
				"id":    proj.ID,
			})
		}

		if proj.Fixtures == nil {
			proj.Fixtures = currentProject.Fixtures
		}
		if proj.Presets == nil {
			proj.Presets = currentProject.Presets
		}
		if proj.Shows == nil {
			proj.Shows = currentProject.Shows
		}

		if err := store.Save(&proj); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to update project",
				"details": err.Error(),
			})
		}

		if enableDMX {
			if err := ws.InitializeDMXController(proj.USBInterface); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Failed to initialize DMX controller",
					"details": err.Error(),
				})
			}
		}

		return c.JSON(proj)
	})

	r.Delete("/", func(c *fiber.Ctx) error {
		project := store.Get()
		if project == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No project to delete",
			})
		}

		if err := store.Delete(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to delete project",
				"details": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
