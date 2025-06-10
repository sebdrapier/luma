package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elano.fr/src/backend/api"
	"elano.fr/src/backend/storage"
	"elano.fr/src/backend/utils"
	"elano.fr/src/backend/ws"
	"elano.fr/src/frontend"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config := utils.LoadConfig()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err := os.MkdirAll(".data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	store, err := storage.NewYAMLStore(config.DataFilePath)
	if err != nil {
		log.Printf("Warning: failed to load project: %v", err)
		log.Println("Creating default project...")
		store, err = storage.NewYAMLStoreWithDefault(config.DataFilePath)
		if err != nil {
			log.Fatalf("Failed to create default project: %v", err)
		}
	}

	ws.SetProjectStore(store)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return ctx.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
		DisableStartupMessage: true,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           120 * time.Second,
		BodyLimit:             4 * 1024 * 1024, // 4MB
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e interface{}) {
			log.Printf("Panic recovered: %v", e)
		},
	}))

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     utils.GetEnv("CORS_ORIGINS", "*"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: false,
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		dmxStatus := "disabled"
		if config.EnableDMX {
			dmxStatus = "enabled"
		}

		return c.JSON(fiber.Map{
			"status":     "healthy",
			"timestamp":  time.Now().Unix(),
			"dmx_status": dmxStatus,
			"version":    "1.0.0",
		})
	})

	api.RegisterUSBRoutes(app)
	api.RegisterFixtureRoutes(app, store)
	api.RegisterPresetRoutes(app, store)
	api.RegisterShowRoutes(app, store)
	api.RegisterProjectRoutes(app, store)

	if config.EnableDMX {
		log.Printf("Initializing DMX controller on %s...", config.DMXPort)
		if err := ws.InitializeDMXController(config.DMXPort); err != nil {
			log.Printf("Warning: Failed to initialize DMX controller: %v", err)
			log.Println("DMX features will be disabled")
		} else {
			log.Println("DMX controller initialized successfully")
		}
	} else {
		log.Println("DMX controller disabled (set ENABLE_DMX=true to enable)")
	}

	ws.SetupWebSocketRoutes(app)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(frontend.DistFS),
		Index:        "index.html",
		NotFoundFile: "index.html", // pour le mode SPA
	}))

	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := ws.CloseDMXController(); err != nil {
			log.Printf("Error closing DMX controller: %v", err)
		}

		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}

		log.Println("Server shutdown complete")
	}()

	log.Printf("🚀 Server starting on http://localhost%s", config.ServerPort)
	log.Printf("📁 Data file: %s", config.DataFilePath)
	log.Printf("🎛️  DMX: %s", func() string {
		if config.EnableDMX {
			return "enabled on " + config.DMXPort
		}
		return "disabled"
	}())

	if err := app.Listen(config.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
