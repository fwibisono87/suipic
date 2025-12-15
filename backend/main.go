package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/suipic/backend/config"
	"github.com/suipic/backend/handlers"
	"github.com/suipic/backend/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	storageService, err := services.NewStorageService(&cfg.MinIO)
	if err != nil {
		log.Fatalf("Failed to initialize storage service: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "Suipic API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: joinStrings(cfg.CORS.Origins, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	setupRoutes(app, storageService)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		log.Printf("Server starting on %s (Environment: %s)", addr, cfg.Server.Env)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func setupRoutes(app *fiber.App, storageService *services.StorageService) {
	photoHandler := handlers.NewPhotoHandler(storageService)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", handlers.HealthCheck)

	photos := v1.Group("/photos")
	photos.Post("/", photoHandler.UploadPhoto)
	photos.Get("/:id", photoHandler.DownloadPhoto)
	photos.Get("/:id/presigned", photoHandler.GetPresignedURL)
	photos.Delete("/:id", photoHandler.DeletePhoto)

	thumbnails := v1.Group("/thumbnails")
	thumbnails.Get("/:id", photoHandler.DownloadThumbnail)
	thumbnails.Get("/:id/presigned", photoHandler.GetPresignedThumbnailURL)
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
