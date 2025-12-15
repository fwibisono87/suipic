package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/suipic/backend/models"
	"github.com/suipic/backend/repository"
)

type PhotoService struct {
	photoRepo      repository.PhotoRepository
	storageService *StorageService
}

func NewPhotoService(photoRepo repository.PhotoRepository, storageService *StorageService) *PhotoService {
	return &PhotoService{
		photoRepo:      photoRepo,
		storageService: storageService,
	}
}

func (s *PhotoService) CreatePhoto(ctx context.Context, albumID int, fileName string, fileReader io.Reader, fileSize int64, contentType string) (*models.Photo, error) {
	data, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	exifData := s.extractEXIF(bytes.NewReader(data))

	uploadResult, err := s.storageService.UploadPhoto(ctx, fileName, bytes.NewReader(data), fileSize, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload photo: %w", err)
	}

	photo := &models.Photo{
		AlbumID:         albumID,
		Filename:        uploadResult.FileID,
		ExifData:        exifData,
		PickRejectState: models.PickRejectNone,
		Stars:           0,
	}

	if dateTime := extractDateTime(exifData); dateTime != nil {
		photo.DateTime = dateTime
	}

	if err := s.photoRepo.Create(ctx, photo); err != nil {
		s.storageService.DeletePhoto(ctx, uploadResult.FileID)
		return nil, fmt.Errorf("failed to create photo record: %w", err)
	}

	return photo, nil
}

func (s *PhotoService) GetPhotoByID(ctx context.Context, id int) (*models.Photo, error) {
	return s.photoRepo.GetByID(ctx, id)
}

func (s *PhotoService) GetPhotosByAlbum(ctx context.Context, albumID int) ([]*models.Photo, error) {
	return s.photoRepo.GetByAlbum(ctx, albumID)
}

func (s *PhotoService) UpdatePhoto(ctx context.Context, photo *models.Photo) error {
	return s.photoRepo.Update(ctx, photo)
}

func (s *PhotoService) DeletePhoto(ctx context.Context, id int) error {
	photo, err := s.photoRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get photo: %w", err)
	}
	if photo == nil {
		return fmt.Errorf("photo not found")
	}

	if err := s.storageService.DeletePhoto(ctx, photo.Filename); err != nil {
		return fmt.Errorf("failed to delete from storage: %w", err)
	}

	if err := s.photoRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete photo record: %w", err)
	}

	return nil
}

func (s *PhotoService) extractEXIF(reader io.Reader) models.ExifData {
	exifData := make(models.ExifData)

	x, err := exif.Decode(reader)
	if err != nil {
		return exifData
	}

	exifFields := []struct {
		name exif.FieldName
		key  string
	}{
		{exif.Make, "Make"},
		{exif.Model, "Model"},
		{exif.LensModel, "LensModel"},
		{exif.FocalLength, "FocalLength"},
		{exif.FNumber, "FNumber"},
		{exif.ExposureTime, "ExposureTime"},
		{exif.ISOSpeedRatings, "ISO"},
		{exif.DateTimeOriginal, "DateTimeOriginal"},
		{exif.ImageWidth, "ImageWidth"},
		{exif.ImageLength, "ImageHeight"},
		{exif.Orientation, "Orientation"},
		{exif.Software, "Software"},
		{exif.Artist, "Artist"},
		{exif.Copyright, "Copyright"},
		{exif.ExposureProgram, "ExposureProgram"},
		{exif.MeteringMode, "MeteringMode"},
		{exif.Flash, "Flash"},
		{exif.WhiteBalance, "WhiteBalance"},
	}

	for _, field := range exifFields {
		if tag, err := x.Get(field.name); err == nil {
			if str, err := tag.StringVal(); err == nil {
				exifData[field.key] = str
			} else if intVal, err := tag.Int(0); err == nil {
				exifData[field.key] = intVal
			} else if ratVal, err := tag.Rat(0); err == nil {
				floatVal, _ := ratVal.Float64()
				exifData[field.key] = floatVal
			}
		}
	}

	if lat, lon, err := x.LatLong(); err == nil {
		exifData["Latitude"] = lat
		exifData["Longitude"] = lon
	}

	return exifData
}

func extractDateTime(exifData models.ExifData) *time.Time {
	if dateStr, ok := exifData["DateTimeOriginal"].(string); ok {
		formats := []string{
			"2006:01:02 15:04:05",
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, dateStr); err == nil {
				return &t
			}
		}
	}
	return nil
}
