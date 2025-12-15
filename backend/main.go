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
	"github.com/suipic/backend/middleware"
	"github.com/suipic/backend/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dbService, err := services.NewDatabaseService(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database service: %v", err)
	}
	defer dbService.Close()

	authService, err := services.NewAuthService(
		dbService,
		cfg.JWT.Secret,
		cfg.JWT.Expiry,
		cfg.Admin.Email,
		cfg.Admin.Password,
		cfg.Admin.Username,
	)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
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

	setupRoutes(app, authService, storageService, dbService)

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

func setupRoutes(app *fiber.App, authService *services.AuthService, storageService *services.StorageService, dbService *services.DatabaseService) {
	authHandler := handlers.NewAuthHandler(authService)
	photoHandler := handlers.NewPhotoHandler(storageService)
	adminHandler := handlers.NewAdminHandler(authService, dbService)
	photographerHandler := handlers.NewPhotographerHandler(authService)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", handlers.HealthCheck)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", middleware.AuthRequired(authService), authHandler.Logout)
	auth.Get("/me", middleware.AuthRequired(authService), authHandler.Me)

	admin := api.Group("/admin")
	admin.Post("/photographers", middleware.AdminOnly(authService), adminHandler.CreatePhotographer)
	admin.Get("/photographers", middleware.AdminOnly(authService), adminHandler.ListPhotographers)

	photos := v1.Group("/photos")
	photos.Post("/", middleware.PhotographerOnly(authService), photoHandler.UploadPhoto)
	photos.Get("/:id", photoHandler.DownloadPhoto)
	photos.Get("/:id/presigned", photoHandler.GetPresignedURL)
	photos.Delete("/:id", middleware.AdminOnly(authService), photoHandler.DeletePhoto)

	thumbnails := v1.Group("/thumbnails")
	thumbnails.Get("/:id", photoHandler.DownloadThumbnail)
	thumbnails.Get("/:id/presigned", photoHandler.GetPresignedThumbnailURL)

	photographer := api.Group("/photographer")
	photographer.Post("/clients", middleware.PhotographerOnly(authService), photographerHandler.CreateOrLinkClient)
	photographer.Get("/clients", middleware.PhotographerOnly(authService), photographerHandler.ListClients)
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
