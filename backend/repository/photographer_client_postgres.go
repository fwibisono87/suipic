package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/suipic/backend/models"
)

type PostgresPhotographerClientRepository struct {
	db *sql.DB
}

func NewPostgresPhotographerClientRepository(db *sql.DB) *PostgresPhotographerClientRepository {
	return &PostgresPhotographerClientRepository{db: db}
}

func (r *PostgresPhotographerClientRepository) Create(ctx context.Context, photographerClient *models.PhotographerClient) error {
	query := `
		INSERT INTO photographer_clients (photographer_id, client_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		photographerClient.PhotographerID,
		photographerClient.ClientID,
	).Scan(&photographerClient.ID, &photographerClient.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create photographer client relationship: %w", err)
	}

	return nil
}

func (r *PostgresPhotographerClientRepository) GetByPhotographerAndClient(ctx context.Context, photographerID, clientID int64) (*models.PhotographerClient, error) {
	query := `
		SELECT id, photographer_id, client_id, created_at
		FROM photographer_clients
		WHERE photographer_id = $1 AND client_id = $2
	`
	photographerClient := &models.PhotographerClient{}
	err := r.db.QueryRowContext(ctx, query, photographerID, clientID).Scan(
		&photographerClient.ID,
		&photographerClient.PhotographerID,
		&photographerClient.ClientID,
		&photographerClient.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get photographer client relationship: %w", err)
	}

	return photographerClient, nil
}

func (r *PostgresPhotographerClientRepository) GetClientsByPhotographer(ctx context.Context, photographerID int64) ([]*models.User, error) {
	query := `
		SELECT u.id, u.username, u.password_hash, u.email, u.friendly_name, u.role, u.created_at, u.updated_at
		FROM users u
		INNER JOIN photographer_clients pc ON u.id = pc.client_id
		WHERE pc.photographer_id = $1
		ORDER BY u.username
	`
	rows, err := r.db.QueryContext(ctx, query, photographerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients by photographer: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.PasswordHash,
			&user.Email,
			&user.FriendlyName,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *PostgresPhotographerClientRepository) Delete(ctx context.Context, photographerID, clientID int64) error {
	query := `DELETE FROM photographer_clients WHERE photographer_id = $1 AND client_id = $2`
	result, err := r.db.ExecContext(ctx, query, photographerID, clientID)
	if err != nil {
		return fmt.Errorf("failed to delete photographer client relationship: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("photographer client relationship not found")
	}

	return nil
}
