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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	Token string `json:"token"`
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

	if req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "password is required")
	}

	if req.Username == "" && req.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "username or email is required")
	}

	user, token, err := h.authService.LoginWithUsernameOrEmail(req.Username, req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "token is required")
	}

	claims, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
	}

	user, err := h.authService.GetUserByID(claims.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve user")
	}

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "user not found")
	}

	newToken, err := h.authService.GenerateToken(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(AuthResponse{
		User:  user,
		Token: newToken,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "logged out successfully",
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
