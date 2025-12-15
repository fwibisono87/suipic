package handlers

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/suipic/backend/models"
	"github.com/suipic/backend/services"
)

type AdminHandler struct {
	authService *services.AuthService
	dbService   *services.DatabaseService
}

func NewAdminHandler(authService *services.AuthService, dbService *services.DatabaseService) *AdminHandler {
	return &AdminHandler{
		authService: authService,
		dbService:   dbService,
	}
}

type CreatePhotographerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type CreatePhotographerResponse struct {
	User     *models.User `json:"user"`
	Password string       `json:"password"`
}

func (h *AdminHandler) CreatePhotographer(c *fiber.Ctx) error {
	var req CreatePhotographerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Email == "" || req.Username == "" {
		return fiber.NewError(fiber.StatusBadRequest, "email and username are required")
	}

	password, err := generateRandomPassword(16)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate password")
	}

	user, err := h.authService.Register(req.Email, req.Username, password, models.RolePhotographer)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(CreatePhotographerResponse{
		User:     user,
		Password: password,
	})
}

func (h *AdminHandler) ListPhotographers(c *fiber.Ctx) error {
	photographers, err := h.dbService.GetUsersByRole(models.RolePhotographer)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch photographers")
	}

	return c.JSON(fiber.Map{
		"photographers": photographers,
	})
}

func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
