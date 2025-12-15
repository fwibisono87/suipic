package handlers

import (
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suipic/backend/services"
)

type PhotoHandler struct {
	storageService *services.StorageService
}

func NewPhotoHandler(storageService *services.StorageService) *PhotoHandler {
	return &PhotoHandler{
		storageService: storageService,
	}
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

func (h *PhotoHandler) DeletePhoto(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file ID is required",
		})
	}

	err := h.storageService.DeletePhoto(c.Context(), fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "failed to delete photo: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "photo deleted successfully",
	})
}
