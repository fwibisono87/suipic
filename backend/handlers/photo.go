package handlers

import (
	"io"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suipic/backend/models"
	"github.com/suipic/backend/services"
)

type PhotoHandler struct {
	storageService *services.StorageService
	photoService   *services.PhotoService
	albumService   *services.AlbumService
}

func NewPhotoHandler(storageService *services.StorageService, photoService *services.PhotoService, albumService *services.AlbumService) *PhotoHandler {
	return &PhotoHandler{
		storageService: storageService,
		photoService:   photoService,
		albumService:   albumService,
	}
}

type UpdatePhotoRequest struct {
	Title           *string `json:"title"`
	PickRejectState *string `json:"pickRejectState"`
	Stars           *int    `json:"stars"`
}

func (h *PhotoHandler) CreatePhoto(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	albumID, err := strconv.Atoi(c.Params("albumId"))
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

	if role != models.RoleAdmin && album.PhotographerID != int(userID) {
		return fiber.NewError(fiber.StatusForbidden, "you can only upload photos to your own albums")
	}

	file, err := c.FormFile("photo")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "photo file is required")
	}

	src, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to open file")
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	photo, err := h.photoService.CreatePhoto(
		c.Context(),
		albumID,
		file.Filename,
		src,
		file.Size,
		contentType,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create photo: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(photo)
}

func (h *PhotoHandler) GetPhotosByAlbum(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	albumID, err := strconv.Atoi(c.Params("albumId"))
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

	photos, err := h.photoService.GetPhotosByAlbum(c.Context(), albumID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get photos: "+err.Error())
	}

	return c.JSON(photos)
}

func (h *PhotoHandler) GetPhoto(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	photoID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid photo id")
	}

	photo, err := h.photoService.GetPhotoByID(c.Context(), photoID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get photo: "+err.Error())
	}
	if photo == nil {
		return fiber.NewError(fiber.StatusNotFound, "photo not found")
	}

	if role != models.RoleAdmin {
		canAccess, err := h.albumService.CanUserAccessAlbum(c.Context(), int(userID), photo.AlbumID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if !canAccess {
			return fiber.NewError(fiber.StatusForbidden, "access denied to this photo")
		}
	}

	return c.JSON(photo)
}

func (h *PhotoHandler) UpdatePhoto(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	photoID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid photo id")
	}

	photo, err := h.photoService.GetPhotoByID(c.Context(), photoID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get photo: "+err.Error())
	}
	if photo == nil {
		return fiber.NewError(fiber.StatusNotFound, "photo not found")
	}

	album, err := h.albumService.GetAlbumByID(c.Context(), photo.AlbumID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get album: "+err.Error())
	}
	if album == nil {
		return fiber.NewError(fiber.StatusNotFound, "album not found")
	}

	if role != models.RoleAdmin {
		canAccess, err := h.albumService.CanUserAccessAlbum(c.Context(), int(userID), photo.AlbumID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if !canAccess {
			return fiber.NewError(fiber.StatusForbidden, "access denied to this photo")
		}
	}

	var req UpdatePhotoRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Title != nil {
		photo.Title = req.Title
	}

	if req.PickRejectState != nil {
		switch *req.PickRejectState {
		case "none", "pick", "reject":
			photo.PickRejectState = models.PickRejectState(*req.PickRejectState)
		default:
			return fiber.NewError(fiber.StatusBadRequest, "invalid pickRejectState, must be 'none', 'pick', or 'reject'")
		}
	}

	if req.Stars != nil {
		if *req.Stars < 0 || *req.Stars > 5 {
			return fiber.NewError(fiber.StatusBadRequest, "stars must be between 0 and 5")
		}
		photo.Stars = *req.Stars
	}

	if err := h.photoService.UpdatePhoto(c.Context(), photo); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update photo: "+err.Error())
	}

	return c.JSON(photo)
}

func (h *PhotoHandler) DeletePhoto(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	role, _ := c.Locals("user_role").(models.UserRole)

	photoID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid photo id")
	}

	photo, err := h.photoService.GetPhotoByID(c.Context(), photoID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get photo: "+err.Error())
	}
	if photo == nil {
		return fiber.NewError(fiber.StatusNotFound, "photo not found")
	}

	album, err := h.albumService.GetAlbumByID(c.Context(), photo.AlbumID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get album: "+err.Error())
	}
	if album == nil {
		return fiber.NewError(fiber.StatusNotFound, "album not found")
	}

	if role != models.RoleAdmin && album.PhotographerID != int(userID) {
		return fiber.NewError(fiber.StatusForbidden, "you can only delete photos from your own albums")
	}

	if err := h.photoService.DeletePhoto(c.Context(), photoID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete photo: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "photo deleted successfully",
	})
}

func (h *PhotoHandler) UploadPhoto(c *fiber.Ctx) error {
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "photo file is required",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to open file",
		})
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	result, err := h.storageService.UploadPhoto(
		c.Context(),
		file.Filename,
		src,
		file.Size,
		contentType,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload photo: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *PhotoHandler) DownloadPhoto(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file ID is required",
		})
	}

	object, info, err := h.storageService.DownloadPhoto(c.Context(), fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "photo not found: " + err.Error(),
		})
	}
	defer object.Close()

	c.Set("Content-Type", info.ContentType)
	c.Set("Content-Disposition", "inline; filename=\""+info.Key+"\"")

	data, err := io.ReadAll(object)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read photo",
		})
	}

	return c.Send(data)
}

func (h *PhotoHandler) DownloadThumbnail(c *fiber.Ctx) error {
	thumbnailID := c.Params("id")
	if thumbnailID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "thumbnail ID is required",
		})
	}

	object, info, err := h.storageService.DownloadThumbnail(c.Context(), thumbnailID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "thumbnail not found: " + err.Error(),
		})
	}
	defer object.Close()

	c.Set("Content-Type", info.ContentType)
	c.Set("Content-Disposition", "inline; filename=\""+info.Key+"\"")

	data, err := io.ReadAll(object)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read thumbnail",
		})
	}

	return c.Send(data)
}

func (h *PhotoHandler) GetPresignedURL(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file ID is required",
		})
	}

	expires := 1 * time.Hour
	presignedURL, err := h.storageService.GetPresignedDownloadURL(c.Context(), fileID, expires)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "failed to generate presigned URL: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"url":        presignedURL,
		"expires_in": int(expires.Seconds()),
	})
}

func (h *PhotoHandler) GetPresignedThumbnailURL(c *fiber.Ctx) error {
	thumbnailID := c.Params("id")
	if thumbnailID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "thumbnail ID is required",
		})
	}

	expires := 1 * time.Hour
	presignedURL, err := h.storageService.GetPresignedThumbnailURL(c.Context(), thumbnailID, expires)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "failed to generate presigned URL: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"url":        presignedURL,
		"expires_in": int(expires.Seconds()),
	})
}
