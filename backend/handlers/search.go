package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suipic/backend/services"
)

type SearchHandler struct {
	esService    *services.ElasticsearchService
	photoService *services.PhotoService
	albumService *services.AlbumService
}

func NewSearchHandler(esService *services.ElasticsearchService, photoService *services.PhotoService, albumService *services.AlbumService) *SearchHandler {
	return &SearchHandler{
		esService:    esService,
		photoService: photoService,
		albumService: albumService,
	}
}

func (h *SearchHandler) Search(c *fiber.Ctx) error {
	if h.esService == nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "search service is not available")
	}

	query := c.Query("q", "")

	filter := &services.SearchFilter{
		Query:  query,
		Limit:  50,
		Offset: 0,
	}

	if albumIDStr := c.Query("album"); albumIDStr != "" {
		albumID, err := strconv.Atoi(albumIDStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid album id")
		}
		filter.AlbumID = &albumID
	}

	if dateFromStr := c.Query("dateFrom"); dateFromStr != "" {
		dateFrom, err := time.Parse(time.RFC3339, dateFromStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid dateFrom format, use RFC3339")
		}
		filter.DateFrom = &dateFrom
	}

	if dateToStr := c.Query("dateTo"); dateToStr != "" {
		dateTo, err := time.Parse(time.RFC3339, dateToStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid dateTo format, use RFC3339")
		}
		filter.DateTo = &dateTo
	}

	if minStarsStr := c.Query("minStars"); minStarsStr != "" {
		minStars, err := strconv.Atoi(minStarsStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid minStars")
		}
		if minStars < 0 || minStars > 5 {
			return fiber.NewError(fiber.StatusBadRequest, "minStars must be between 0 and 5")
		}
		filter.MinStars = &minStars
	}

	if maxStarsStr := c.Query("maxStars"); maxStarsStr != "" {
		maxStars, err := strconv.Atoi(maxStarsStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid maxStars")
		}
		if maxStars < 0 || maxStars > 5 {
			return fiber.NewError(fiber.StatusBadRequest, "maxStars must be between 0 and 5")
		}
		filter.MaxStars = &maxStars
	}

	if state := c.Query("state"); state != "" {
		if state != "none" && state != "pick" && state != "reject" {
			return fiber.NewError(fiber.StatusBadRequest, "state must be 'none', 'pick', or 'reject'")
		}
		filter.State = &state
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid limit")
		}
		if limit > 0 && limit <= 1000 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid offset")
		}
		if offset >= 0 {
			filter.Offset = offset
		}
	}

	result, err := h.esService.Search(c.Context(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "search failed: "+err.Error())
	}

	return c.JSON(result)
}

func (h *SearchHandler) BulkIndexAlbum(c *fiber.Ctx) error {
	if h.esService == nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "search service is not available")
	}

	albumID, err := strconv.Atoi(c.Params("albumId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid album id")
	}

	if err := h.photoService.BulkIndexPhotosByAlbum(c.Context(), albumID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to bulk index album: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "album photos indexed successfully",
	})
}
