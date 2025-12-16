package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suipic/backend/models"
	"github.com/suipic/backend/services"
)

type AlbumHandler struct {
	albumService *services.AlbumService
}

func NewAlbumHandler(albumService *services.AlbumService) *AlbumHandler {
	return &AlbumHandler{
		albumService: albumService,
	}
}

type CreateAlbumRequest struct {
	Title            string                 `json:"title"`
	DateTaken        *string                `json:"dateTaken"`
	Description      *string                `json:"description"`
	Location         *string                `json:"location"`
	CustomFields     map[string]interface{} `json:"customFields"`
	ThumbnailPhotoID *int                   `json:"thumbnailPhotoId"`
}

type UpdateAlbumRequest struct {
	Title            string                 `json:"title"`
	DateTaken        *string                `json:"dateTaken"`
	Description      *string                `json:"description"`
	Location         *string                `json:"location"`
	CustomFields     map[string]interface{} `json:"customFields"`
	ThumbnailPhotoID *int                   `json:"thumbnailPhotoId"`
}

type AssignUsersRequest struct {
	UserIDs []int `json:"userIds"`
}

func (h *AlbumHandler) CreateAlbum(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, ok := c.Locals("user_role").(models.UserRole)
	if !ok || (role != models.RolePhotographer && role != models.RoleAdmin) {
		return fiber.NewError(fiber.StatusForbidden, "only photographers can create albums")
	}

	var req CreateAlbumRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Title == "" {
		return fiber.NewError(fiber.StatusBadRequest, "title is required")
	}

	album := &models.Album{
		Title:            req.Title,
		Description:      req.Description,
		Location:         req.Location,
		CustomFields:     req.CustomFields,
		ThumbnailPhotoID: req.ThumbnailPhotoID,
		PhotographerID:   int(userID),
	}

	if req.DateTaken != nil && *req.DateTaken != "" {
		dateTaken, err := parseDateTime(*req.DateTaken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid dateTaken format")
		}
		album.DateTaken = &dateTaken
	}

	if err := h.albumService.CreateAlbum(c.Context(), album); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create album: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(album)
}

func (h *AlbumHandler) GetAlbums(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	photographerIDStr := c.Query("photographerId")
	var albums []*models.Album
	var err error

	if photographerIDStr != "" {
		photographerID, parseErr := strconv.Atoi(photographerIDStr)
		if parseErr != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid photographerId")
		}

		if role == models.RolePhotographer && photographerID != int(userID) {
			return fiber.NewError(fiber.StatusForbidden, "photographers can only view their own albums")
		}

		albums, err = h.albumService.ListAlbums(c.Context(), &photographerID, nil)
	} else if role == models.RolePhotographer {
		photographerID := int(userID)
		albums, err = h.albumService.ListAlbums(c.Context(), &photographerID, nil)
	} else if role == models.RoleClient {
		userIDInt := int(userID)
		albums, err = h.albumService.ListAlbums(c.Context(), nil, &userIDInt)
	} else {
		albums, err = h.albumService.ListAlbums(c.Context(), nil, nil)
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to list albums: "+err.Error())
	}

	return c.JSON(albums)
}

func (h *AlbumHandler) GetAlbum(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	albumID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid album id")
	}

	album, err := h.albumService.GetAlbumByID(c.Context(), albumID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get album: "+err.Error())
	}
	if album == nil {
		return fiber.NewError(fiber.StatusNotFound, "album not found")
	}

	if role != models.RoleAdmin {
		canAccess, err := h.albumService.CanUserAccessAlbum(c.Context(), int(userID), albumID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if !canAccess {
			return fiber.NewError(fiber.StatusForbidden, "access denied to this album")
		}
	}

	return c.JSON(album)
}

func (h *AlbumHandler) UpdateAlbum(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	albumID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid album id")
	}

	existingAlbum, err := h.albumService.GetAlbumByID(c.Context(), albumID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get album: "+err.Error())
	}
	if existingAlbum == nil {
		return fiber.NewError(fiber.StatusNotFound, "album not found")
	}

	if role != models.RoleAdmin && existingAlbum.PhotographerID != int(userID) {
		return fiber.NewError(fiber.StatusForbidden, "you can only update your own albums")
	}

	var req UpdateAlbumRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Title == "" {
		return fiber.NewError(fiber.StatusBadRequest, "title is required")
	}

	existingAlbum.Title = req.Title
	existingAlbum.Description = req.Description
	existingAlbum.Location = req.Location
	existingAlbum.CustomFields = req.CustomFields
	existingAlbum.ThumbnailPhotoID = req.ThumbnailPhotoID

	if req.DateTaken != nil && *req.DateTaken != "" {
		dateTaken, err := parseDateTime(*req.DateTaken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid dateTaken format")
		}
		existingAlbum.DateTaken = &dateTaken
	} else {
		existingAlbum.DateTaken = nil
	}

	if err := h.albumService.UpdateAlbum(c.Context(), existingAlbum); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update album: "+err.Error())
	}

	return c.JSON(existingAlbum)
}

func (h *AlbumHandler) DeleteAlbum(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	albumID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid album id")
	}

	existingAlbum, err := h.albumService.GetAlbumByID(c.Context(), albumID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get album: "+err.Error())
	}
	if existingAlbum == nil {
		return fiber.NewError(fiber.StatusNotFound, "album not found")
	}

	if role != models.RoleAdmin && existingAlbum.PhotographerID != int(userID) {
		return fiber.NewError(fiber.StatusForbidden, "you can only delete your own albums")
	}

	if err := h.albumService.DeleteAlbum(c.Context(), albumID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete album: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "album deleted successfully",
	})
}

func (h *AlbumHandler) AssignUsers(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	albumID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid album id")
	}

	existingAlbum, err := h.albumService.GetAlbumByID(c.Context(), albumID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get album: "+err.Error())
	}
	if existingAlbum == nil {
		return fiber.NewError(fiber.StatusNotFound, "album not found")
	}

	if role != models.RoleAdmin && existingAlbum.PhotographerID != int(userID) {
		return fiber.NewError(fiber.StatusForbidden, "you can only assign users to your own albums")
	}

	var req AssignUsersRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if len(req.UserIDs) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "userIds is required")
	}

	if err := h.albumService.AssignUsersToAlbum(c.Context(), albumID, req.UserIDs); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to assign users: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "users assigned successfully",
	})
}

func (h *AlbumHandler) GetAlbumUsers(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	albumID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid album id")
	}

	existingAlbum, err := h.albumService.GetAlbumByID(c.Context(), albumID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get album: "+err.Error())
	}
	if existingAlbum == nil {
		return fiber.NewError(fiber.StatusNotFound, "album not found")
	}

	if role != models.RoleAdmin && existingAlbum.PhotographerID != int(userID) {
		return fiber.NewError(fiber.StatusForbidden, "you can only view users for your own albums")
	}

	albumUsers, err := h.albumService.GetAlbumUsers(c.Context(), albumID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get album users: "+err.Error())
	}

	return c.JSON(albumUsers)
}

func parseDateTime(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid date format")
}
