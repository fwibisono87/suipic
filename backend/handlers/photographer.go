package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suipic/backend/models"
	"github.com/suipic/backend/services"
)

type PhotographerHandler struct {
	authService *services.AuthService
}

func NewPhotographerHandler(authService *services.AuthService) *PhotographerHandler {
	return &PhotographerHandler{
		authService: authService,
	}
}

type CreateClientRequest struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type ClientResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

func (h *PhotographerHandler) CreateOrLinkClient(c *fiber.Ctx) error {
	photographerID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	var req CreateClientRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Username == "" {
		return fiber.NewError(fiber.StatusBadRequest, "username is required")
	}

	existingUser, err := h.authService.GetUserByUsername(req.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to check username")
	}

	var client *models.User
	if existingUser != nil {
		if existingUser.Role != models.RoleClient {
			return fiber.NewError(fiber.StatusBadRequest, "user is not a client")
		}
		client = existingUser
	} else {
		if req.Email == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "email and password are required for new client")
		}

		client, err = h.authService.Register(req.Email, req.Username, req.Password, models.RoleClient)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	existingRelation, err := h.authService.GetPhotographerClient(photographerID, client.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to check existing relationship")
	}

	if existingRelation != nil {
		return fiber.NewError(fiber.StatusConflict, "client already linked to photographer")
	}

	_, err = h.authService.CreatePhotographerClient(photographerID, client.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to link client")
	}

	return c.Status(fiber.StatusCreated).JSON(ClientResponse{
		ID:        client.ID,
		Username:  client.Username,
		Email:     client.Email,
		Role:      string(client.Role),
		CreatedAt: client.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *PhotographerHandler) ListClients(c *fiber.Ctx) error {
	photographerID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	clients, err := h.authService.GetClientsByPhotographer(photographerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get clients")
	}

	response := make([]ClientResponse, 0, len(clients))
	for _, client := range clients {
		response = append(response, ClientResponse{
			ID:        client.ID,
			Username:  client.Username,
			Email:     client.Email,
			Role:      string(client.Role),
			CreatedAt: client.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return c.JSON(response)
}
