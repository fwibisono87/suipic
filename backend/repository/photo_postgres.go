package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/suipic/backend/models"
)

type PostgresPhotoRepository struct {
	db *sql.DB
}

func NewPostgresPhotoRepository(db *sql.DB) *PostgresPhotoRepository {
	return &PostgresPhotoRepository{db: db}
}

func (r *PostgresPhotoRepository) Create(ctx context.Context, photo *models.Photo) error {
	query := `
		INSERT INTO photos (album_id, filename, title, date_time, exif_data, pick_reject_state, stars, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		photo.AlbumID,
		photo.Filename,
		photo.Title,
		photo.DateTime,
		photo.ExifData,
		photo.PickRejectState,
		photo.Stars,
	).Scan(&photo.ID, &photo.CreatedAt, &photo.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create photo: %w", err)
	}

	return nil
}

func (r *PostgresPhotoRepository) GetByID(ctx context.Context, id int) (*models.Photo, error) {
	query := `
		SELECT id, album_id, filename, title, date_time, exif_data, pick_reject_state, stars, created_at, updated_at
		FROM photos
		WHERE id = $1
	`
	photo := &models.Photo{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&photo.ID,
		&photo.AlbumID,
		&photo.Filename,
		&photo.Title,
		&photo.DateTime,
		&photo.ExifData,
		&photo.PickRejectState,
		&photo.Stars,
		&photo.CreatedAt,
		&photo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get photo by id: %w", err)
	}

	return photo, nil
}

func (r *PostgresPhotoRepository) Update(ctx context.Context, photo *models.Photo) error {
	query := `
		UPDATE photos
		SET album_id = $1, filename = $2, title = $3, date_time = $4, exif_data = $5, pick_reject_state = $6, stars = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		photo.AlbumID,
		photo.Filename,
		photo.Title,
		photo.DateTime,
		photo.ExifData,
		photo.PickRejectState,
		photo.Stars,
		photo.ID,
	).Scan(&photo.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("photo not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update photo: %w", err)
	}

	return nil
}

func (r *PostgresPhotoRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM photos WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("photo not found")
	}

	return nil
}

func (r *PostgresPhotoRepository) List(ctx context.Context, limit, offset int) ([]*models.Photo, error) {
	query := `
		SELECT id, album_id, filename, title, date_time, exif_data, pick_reject_state, stars, created_at, updated_at
		FROM photos
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list photos: %w", err)
	}
	defer rows.Close()

	var photos []*models.Photo
	for rows.Next() {
		photo := &models.Photo{}
		err := rows.Scan(
			&photo.ID,
			&photo.AlbumID,
			&photo.Filename,
			&photo.Title,
			&photo.DateTime,
			&photo.ExifData,
			&photo.PickRejectState,
			&photo.Stars,
			&photo.CreatedAt,
			&photo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating photos: %w", err)
	}

	return photos, nil
}

func (r *PostgresPhotoRepository) GetByAlbum(ctx context.Context, albumID int) ([]*models.Photo, error) {
	query := `
		SELECT id, album_id, filename, title, date_time, exif_data, pick_reject_state, stars, created_at, updated_at
		FROM photos
		WHERE album_id = $1
		ORDER BY date_time DESC NULLS LAST, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos by album: %w", err)
	}
	defer rows.Close()

	var photos []*models.Photo
	for rows.Next() {
		photo := &models.Photo{}
		err := rows.Scan(
			&photo.ID,
			&photo.AlbumID,
			&photo.Filename,
			&photo.Title,
			&photo.DateTime,
			&photo.ExifData,
			&photo.PickRejectState,
			&photo.Stars,
			&photo.CreatedAt,
			&photo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating photos: %w", err)
	}

	return photos, nil
}
