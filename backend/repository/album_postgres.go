package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/suipic/backend/models"
)

type PostgresAlbumRepository struct {
	db *sql.DB
}

func NewPostgresAlbumRepository(db *sql.DB) *PostgresAlbumRepository {
	return &PostgresAlbumRepository{db: db}
}

func (r *PostgresAlbumRepository) Create(ctx context.Context, album *models.Album) error {
	query := `
		INSERT INTO albums (title, date_taken, description, location, custom_fields, thumbnail_photo_id, photographer_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		album.Title,
		album.DateTaken,
		album.Description,
		album.Location,
		album.CustomFields,
		album.ThumbnailPhotoID,
		album.PhotographerID,
	).Scan(&album.ID, &album.CreatedAt, &album.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create album: %w", err)
	}

	return nil
}

func (r *PostgresAlbumRepository) GetByID(ctx context.Context, id int) (*models.Album, error) {
	query := `
		SELECT id, title, date_taken, description, location, custom_fields, thumbnail_photo_id, photographer_id, created_at, updated_at
		FROM albums
		WHERE id = $1
	`
	album := &models.Album{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&album.ID,
		&album.Title,
		&album.DateTaken,
		&album.Description,
		&album.Location,
		&album.CustomFields,
		&album.ThumbnailPhotoID,
		&album.PhotographerID,
		&album.CreatedAt,
		&album.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get album by id: %w", err)
	}

	return album, nil
}

func (r *PostgresAlbumRepository) Update(ctx context.Context, album *models.Album) error {
	query := `
		UPDATE albums
		SET title = $1, date_taken = $2, description = $3, location = $4, custom_fields = $5, thumbnail_photo_id = $6, photographer_id = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		album.Title,
		album.DateTaken,
		album.Description,
		album.Location,
		album.CustomFields,
		album.ThumbnailPhotoID,
		album.PhotographerID,
		album.ID,
	).Scan(&album.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("album not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update album: %w", err)
	}

	return nil
}

func (r *PostgresAlbumRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM albums WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete album: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("album not found")
	}

	return nil
}

func (r *PostgresAlbumRepository) List(ctx context.Context, limit, offset int) ([]*models.Album, error) {
	query := `
		SELECT id, title, date_taken, description, location, custom_fields, thumbnail_photo_id, photographer_id, created_at, updated_at
		FROM albums
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list albums: %w", err)
	}
	defer rows.Close()

	var albums []*models.Album
	for rows.Next() {
		album := &models.Album{}
		err := rows.Scan(
			&album.ID,
			&album.Title,
			&album.DateTaken,
			&album.Description,
			&album.Location,
			&album.CustomFields,
			&album.ThumbnailPhotoID,
			&album.PhotographerID,
			&album.CreatedAt,
			&album.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, album)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating albums: %w", err)
	}

	return albums, nil
}

func (r *PostgresAlbumRepository) GetByPhotographer(ctx context.Context, photographerID int) ([]*models.Album, error) {
	query := `
		SELECT id, title, date_taken, description, location, custom_fields, thumbnail_photo_id, photographer_id, created_at, updated_at
		FROM albums
		WHERE photographer_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, photographerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get albums by photographer: %w", err)
	}
	defer rows.Close()

	var albums []*models.Album
	for rows.Next() {
		album := &models.Album{}
		err := rows.Scan(
			&album.ID,
			&album.Title,
			&album.DateTaken,
			&album.Description,
			&album.Location,
			&album.CustomFields,
			&album.ThumbnailPhotoID,
			&album.PhotographerID,
			&album.CreatedAt,
			&album.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, album)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating albums: %w", err)
	}

	return albums, nil
}
