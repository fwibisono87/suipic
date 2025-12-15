package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suipic/backend/models"
	"github.com/suipic/backend/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Email    string          `json:"email"`
	Username string          `json:"username"`
	Password string          `json:"password"`
	Role     models.UserRole `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Email == "" || req.Username == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "email, username, and password are required")
	}

	if req.Role == "" {
		req.Role = models.RoleClient
	}

	if req.Role == models.RoleAdmin {
		return fiber.NewError(fiber.StatusBadRequest, "cannot register as admin")
	}

	user, err := h.authService.Register(req.Email, req.Username, req.Password, req.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "email and password are required")
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	email, _ := c.Locals("user_email").(string)
	username, _ := c.Locals("user_username").(string)
	role, _ := c.Locals("user_role").(models.UserRole)

	user := &models.User{
		ID:       userID,
		Email:    email,
		Username: username,
		Role:     role,
	}

	return c.JSON(user)
}
