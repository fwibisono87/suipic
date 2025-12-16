package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suipic/backend/services"
)

type SettingsHandler struct {
	systemSettingsService *services.SystemSettingsService
}

func NewSettingsHandler(systemSettingsService *services.SystemSettingsService) *SettingsHandler {
	return &SettingsHandler{
		systemSettingsService: systemSettingsService,
	}
}

type PublicSettingsResponse struct {
	ImageProtectionEnabled bool `json:"image_protection_enabled"`
}

func (h *SettingsHandler) GetPublicSettings(c *fiber.Ctx) error {
	imageProtectionEnabled, err := h.systemSettingsService.GetImageProtectionEnabled(c.Context())
	if err != nil {
		imageProtectionEnabled = false
	}

	return c.JSON(PublicSettingsResponse{
		ImageProtectionEnabled: imageProtectionEnabled,
	})
}
