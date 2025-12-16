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

	esService, err := services.NewElasticsearchService(&cfg.Elasticsearch)
	if err != nil {
		log.Printf("Warning: Failed to initialize elasticsearch service: %v", err)
		esService = nil
	}

	albumService := services.NewAlbumService(dbService.GetDB())
	commentService := services.NewCommentService(dbService.GetCommentRepo(), dbService.GetUserRepo())
	photoService := services.NewPhotoService(dbService.GetPhotoRepo(), storageService, esService, albumService, dbService.GetCommentRepo())
	systemSettingsService := services.NewSystemSettingsService(dbService.GetSystemSettingsRepo())

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

	setupRoutes(app, authService, storageService, dbService, albumService, photoService, commentService, esService, systemSettingsService)

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

func setupRoutes(app *fiber.App, authService *services.AuthService, storageService *services.StorageService, dbService *services.DatabaseService, albumService *services.AlbumService, photoService *services.PhotoService, commentService *services.CommentService, esService *services.ElasticsearchService, systemSettingsService *services.SystemSettingsService) {
	_ = systemSettingsService
	authHandler := handlers.NewAuthHandler(authService)
	photoHandler := handlers.NewPhotoHandler(storageService, photoService, albumService, commentService, esService)
	albumHandler := handlers.NewAlbumHandler(albumService)
	adminHandler := handlers.NewAdminHandler(authService, dbService)
	photographerHandler := handlers.NewPhotographerHandler(authService)
	searchHandler := handlers.NewSearchHandler(esService, photoService, albumService)

	api := app.Group("/api")

	api.Get("/health", handlers.HealthCheck)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", middleware.AuthRequired(authService), authHandler.Logout)
	auth.Get("/me", middleware.AuthRequired(authService), authHandler.Me)

	admin := api.Group("/admin")
	admin.Post("/photographers", middleware.AdminOnly(authService), adminHandler.CreatePhotographer)
	admin.Get("/photographers", middleware.AdminOnly(authService), adminHandler.ListPhotographers)
	admin.Get("/settings", middleware.AdminOnly(authService), adminHandler.GetSettings)
	admin.Put("/settings/:key", middleware.AdminOnly(authService), adminHandler.UpdateSetting)

	albums := api.Group("/albums")
	albums.Post("/", middleware.AuthRequired(authService), albumHandler.CreateAlbum)
	albums.Get("/", middleware.AuthRequired(authService), albumHandler.GetAlbums)
	albums.Get("/:id", middleware.AuthRequired(authService), albumHandler.GetAlbum)
	albums.Put("/:id", middleware.AuthRequired(authService), albumHandler.UpdateAlbum)
	albums.Delete("/:id", middleware.AuthRequired(authService), albumHandler.DeleteAlbum)
	albums.Post("/:id/users", middleware.AuthRequired(authService), albumHandler.AssignUsers)
	albums.Get("/:id/users", middleware.AuthRequired(authService), albumHandler.GetAlbumUsers)
	albums.Post("/:albumId/photos", middleware.AuthRequired(authService), photoHandler.CreatePhoto)
	albums.Get("/:albumId/photos", middleware.AuthRequired(authService), photoHandler.GetPhotosByAlbum)

	photos := api.Group("/photos")
	photos.Post("/", middleware.PhotographerOnly(authService), photoHandler.UploadPhoto)
	photos.Get("/:id", middleware.AuthRequired(authService), photoHandler.GetPhoto)
	photos.Put("/:id", middleware.AuthRequired(authService), photoHandler.UpdatePhoto)
	photos.Delete("/:id", middleware.AuthRequired(authService), photoHandler.DeletePhoto)
	photos.Get("/:id/download", photoHandler.DownloadPhoto)
	photos.Get("/:id/presigned", photoHandler.GetPresignedURL)
	photos.Put("/:id/state", middleware.AuthRequired(authService), photoHandler.SetPhotoState)
	photos.Put("/:id/stars", middleware.AuthRequired(authService), photoHandler.SetPhotoStars)
	photos.Post("/:id/comments", middleware.AuthRequired(authService), photoHandler.CreateComment)
	photos.Get("/:id/comments", middleware.AuthRequired(authService), photoHandler.GetComments)

	thumbnails := api.Group("/thumbnails")
	thumbnails.Get("/:id", photoHandler.DownloadThumbnail)
	thumbnails.Get("/:id/presigned", photoHandler.GetPresignedThumbnailURL)

	photographer := api.Group("/photographer")
	photographer.Post("/clients", middleware.PhotographerOnly(authService), photographerHandler.CreateOrLinkClient)
	photographer.Get("/clients", middleware.PhotographerOnly(authService), photographerHandler.ListClients)
	photographer.Get("/clients/search", middleware.PhotographerOnly(authService), photographerHandler.SearchClients)

	api.Get("/search", middleware.AuthRequired(authService), searchHandler.Search)
	api.Post("/albums/:albumId/index", middleware.AuthRequired(authService), searchHandler.BulkIndexAlbum)
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
