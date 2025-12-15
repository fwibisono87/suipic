package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/suipic/backend/models"
)

type PostgresAlbumUserRepository struct {
	db *sql.DB
}

func NewPostgresAlbumUserRepository(db *sql.DB) *PostgresAlbumUserRepository {
	return &PostgresAlbumUserRepository{db: db}
}

func (r *PostgresAlbumUserRepository) Create(ctx context.Context, albumUser *models.AlbumUser) error {
	query := `
		INSERT INTO album_users (album_id, user_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		albumUser.AlbumID,
		albumUser.UserID,
	).Scan(&albumUser.ID, &albumUser.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create album user: %w", err)
	}

	return nil
}

func (r *PostgresAlbumUserRepository) GetByID(ctx context.Context, id int) (*models.AlbumUser, error) {
	query := `
		SELECT id, album_id, user_id, created_at
		FROM album_users
		WHERE id = $1
	`
	albumUser := &models.AlbumUser{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&albumUser.ID,
		&albumUser.AlbumID,
		&albumUser.UserID,
		&albumUser.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get album user by id: %w", err)
	}

	return albumUser, nil
}

func (r *PostgresAlbumUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM album_users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete album user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("album user not found")
	}

	return nil
}

func (r *PostgresAlbumUserRepository) DeleteByAlbumAndUser(ctx context.Context, albumID, userID int) error {
	query := `DELETE FROM album_users WHERE album_id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, albumID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete album user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("album user not found")
	}

	return nil
}

func (r *PostgresAlbumUserRepository) List(ctx context.Context, limit, offset int) ([]*models.AlbumUser, error) {
	query := `
		SELECT id, album_id, user_id, created_at
		FROM album_users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list album users: %w", err)
	}
	defer rows.Close()

	var albumUsers []*models.AlbumUser
	for rows.Next() {
		albumUser := &models.AlbumUser{}
		err := rows.Scan(
			&albumUser.ID,
			&albumUser.AlbumID,
			&albumUser.UserID,
			&albumUser.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album user: %w", err)
		}
		albumUsers = append(albumUsers, albumUser)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating album users: %w", err)
	}

	return albumUsers, nil
}

func (r *PostgresAlbumUserRepository) GetByAlbum(ctx context.Context, albumID int) ([]*models.AlbumUser, error) {
	query := `
		SELECT id, album_id, user_id, created_at
		FROM album_users
		WHERE album_id = $1
		ORDER BY created_at
	`
	rows, err := r.db.QueryContext(ctx, query, albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get album users by album: %w", err)
	}
	defer rows.Close()

	var albumUsers []*models.AlbumUser
	for rows.Next() {
		albumUser := &models.AlbumUser{}
		err := rows.Scan(
			&albumUser.ID,
			&albumUser.AlbumID,
			&albumUser.UserID,
			&albumUser.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album user: %w", err)
		}
		albumUsers = append(albumUsers, albumUser)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating album users: %w", err)
	}

	return albumUsers, nil
}

func (r *PostgresAlbumUserRepository) GetByUser(ctx context.Context, userID int) ([]*models.AlbumUser, error) {
	query := `
		SELECT id, album_id, user_id, created_at
		FROM album_users
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get album users by user: %w", err)
	}
	defer rows.Close()

	var albumUsers []*models.AlbumUser
	for rows.Next() {
		albumUser := &models.AlbumUser{}
		err := rows.Scan(
			&albumUser.ID,
			&albumUser.AlbumID,
			&albumUser.UserID,
			&albumUser.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album user: %w", err)
		}
		albumUsers = append(albumUsers, albumUser)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating album users: %w", err)
	}

	return albumUsers, nil
}
