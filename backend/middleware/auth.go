package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/suipic/backend/models"
	"github.com/suipic/backend/services"
)

func authenticate(c *fiber.Ctx, authService *services.AuthService) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header format")
	}

	token := parts[1]
	claims, err := authService.ValidateToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("user_email", claims.Email)
	c.Locals("user_username", claims.Username)
	c.Locals("user_role", claims.Role)

	return nil
}

func AuthRequired(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := authenticate(c, authService); err != nil {
			return err
		}
		return c.Next()
	}
}

func AdminOnly(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := authenticate(c, authService); err != nil {
			return err
		}

		role, ok := c.Locals("user_role").(models.UserRole)
		if !ok || role != models.RoleAdmin {
			return fiber.NewError(fiber.StatusForbidden, "admin access required")
		}

		return c.Next()
	}
}

func PhotographerOnly(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := authenticate(c, authService); err != nil {
			return err
		}

		role, ok := c.Locals("user_role").(models.UserRole)
		if !ok || (role != models.RolePhotographer && role != models.RoleAdmin) {
			return fiber.NewError(fiber.StatusForbidden, "photographer access required")
		}

		return c.Next()
	}
}
