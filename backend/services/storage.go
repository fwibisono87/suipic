package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/suipic/backend/config"
	_ "golang.org/x/image/webp"
)

type StorageService struct {
	client     *minio.Client
	bucketName string
	config     *config.MinIOConfig
}

type UploadResult struct {
	FileID      string    `json:"file_id"`
	FileName    string    `json:"file_name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	ThumbnailID string    `json:"thumbnail_id,omitempty"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

const (
	thumbnailWidth  = 300
	thumbnailHeight = 300
	thumbnailPrefix = "thumbnails/"
	photosPrefix    = "photos/"
)

func NewStorageService(cfg *config.MinIOConfig) (*StorageService, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	service := &StorageService{
		client:     client,
		bucketName: cfg.Bucket,
		config:     cfg,
	}

	if err := service.InitializeBucket(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize bucket: %w", err)
	}

	return service, nil
}

func (s *StorageService) InitializeBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/photos/*", "arn:aws:s3:::%s/thumbnails/*"]
			}
		]
	}`, s.bucketName, s.bucketName)

	err = s.client.SetBucketPolicy(ctx, s.bucketName, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return nil
}

func (s *StorageService) UploadPhoto(ctx context.Context, fileName string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	fileID := uuid.New().String()
	objectName := fmt.Sprintf("%s%s.webp", photosPrefix, fileID)

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	var webpData []byte
	if isImageContentType(contentType) {
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %w", err)
		}

		var buf bytes.Buffer
		if err := webp.Encode(&buf, img, &webp.Options{Quality: 85}); err != nil {
			return nil, fmt.Errorf("failed to encode image as WebP: %w", err)
		}
		webpData = buf.Bytes()
	} else {
		webpData = data
	}

	_, err = s.client.PutObject(
		ctx,
		s.bucketName,
		objectName,
		bytes.NewReader(webpData),
		int64(len(webpData)),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload photo: %w", err)
	}

	result := &UploadResult{
		FileID:      fileID,
		FileName:    fileName,
		Size:        int64(len(webpData)),
		ContentType: "image/webp",
		UploadedAt:  time.Now(),
	}

	if isImageContentType(contentType) {
		thumbnailID, err := s.generateThumbnail(ctx, fileID, bytes.NewReader(data))
		if err != nil {
			return result, nil
		}
		result.ThumbnailID = thumbnailID
	}

	return result, nil
}

func (s *StorageService) generateThumbnail(ctx context.Context, fileID string, reader io.Reader) (string, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	thumbnail := imaging.Fit(img, thumbnailWidth, thumbnailHeight, imaging.Lanczos)

	var buf bytes.Buffer
	if err := webp.Encode(&buf, thumbnail, &webp.Options{Quality: 85}); err != nil {
		return "", fmt.Errorf("failed to encode thumbnail as WebP: %w", err)
	}

	thumbnailID := fileID
	thumbnailName := fmt.Sprintf("%s%s.webp", thumbnailPrefix, thumbnailID)

	_, err = s.client.PutObject(
		ctx,
		s.bucketName,
		thumbnailName,
		bytes.NewReader(buf.Bytes()),
		int64(buf.Len()),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	return thumbnailID, nil
}

func (s *StorageService) DownloadPhoto(ctx context.Context, fileID string) (io.ReadCloser, *minio.ObjectInfo, error) {
	objects := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix: photosPrefix + fileID,
	})

	var objectName string
	for obj := range objects {
		if obj.Err != nil {
			return nil, nil, fmt.Errorf("error listing objects: %w", obj.Err)
		}
		objectName = obj.Key
		break
	}

	if objectName == "" {
		return nil, nil, fmt.Errorf("photo not found")
	}

	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get photo: %w", err)
	}

	info, err := object.Stat()
	if err != nil {
		object.Close()
		return nil, nil, fmt.Errorf("failed to stat photo: %w", err)
	}

	return object, &info, nil
}

func (s *StorageService) DownloadThumbnail(ctx context.Context, thumbnailID string) (io.ReadCloser, *minio.ObjectInfo, error) {
	objects := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix: thumbnailPrefix + thumbnailID,
	})

	var objectName string
	for obj := range objects {
		if obj.Err != nil {
			return nil, nil, fmt.Errorf("error listing objects: %w", obj.Err)
		}
		objectName = obj.Key
		break
	}

	if objectName == "" {
		return nil, nil, fmt.Errorf("thumbnail not found")
	}

	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get thumbnail: %w", err)
	}

	info, err := object.Stat()
	if err != nil {
		object.Close()
		return nil, nil, fmt.Errorf("failed to stat thumbnail: %w", err)
	}

	return object, &info, nil
}

func (s *StorageService) GetPresignedDownloadURL(ctx context.Context, fileID string, expires time.Duration) (string, error) {
	objects := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix: photosPrefix + fileID,
	})

	var objectName string
	for obj := range objects {
		if obj.Err != nil {
			return "", fmt.Errorf("error listing objects: %w", obj.Err)
		}
		objectName = obj.Key
		break
	}

	if objectName == "" {
		return "", fmt.Errorf("photo not found")
	}

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expires, url.Values{})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

func (s *StorageService) GetPresignedThumbnailURL(ctx context.Context, thumbnailID string, expires time.Duration) (string, error) {
	objects := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix: thumbnailPrefix + thumbnailID,
	})

	var objectName string
	for obj := range objects {
		if obj.Err != nil {
			return "", fmt.Errorf("error listing objects: %w", obj.Err)
		}
		objectName = obj.Key
		break
	}

	if objectName == "" {
		return "", fmt.Errorf("thumbnail not found")
	}

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expires, url.Values{})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

func (s *StorageService) DeletePhoto(ctx context.Context, fileID string) error {
	objects := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix: photosPrefix + fileID,
	})

	var objectName string
	for obj := range objects {
		if obj.Err != nil {
			return fmt.Errorf("error listing objects: %w", obj.Err)
		}
		objectName = obj.Key
		break
	}

	if objectName == "" {
		return fmt.Errorf("photo not found")
	}

	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	thumbnailObjects := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix: thumbnailPrefix + fileID,
	})

	for obj := range thumbnailObjects {
		if obj.Err != nil {
			continue
		}
		s.client.RemoveObject(ctx, s.bucketName, obj.Key, minio.RemoveObjectOptions{})
	}

	return nil
}

func isImageContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}
