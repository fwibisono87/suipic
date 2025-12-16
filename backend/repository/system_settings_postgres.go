package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresSystemSettingsRepository struct {
	db *sql.DB
}

func NewPostgresSystemSettingsRepository(db *sql.DB) *PostgresSystemSettingsRepository {
	return &PostgresSystemSettingsRepository{db: db}
}

func (r *PostgresSystemSettingsRepository) Get(ctx context.Context, key string) (string, error) {
	query := `SELECT value FROM settings WHERE key = $1`
	var value string
	err := r.db.QueryRowContext(ctx, query, key).Scan(&value)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("setting not found: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get setting: %w", err)
	}

	return value, nil
}

func (r *PostgresSystemSettingsRepository) Set(ctx context.Context, key string, value string) error {
	query := `
		INSERT INTO settings (key, value, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (key)
		DO UPDATE SET value = $2, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, key, value)
	if err != nil {
		return fmt.Errorf("failed to set setting: %w", err)
	}

	return nil
}

func (r *PostgresSystemSettingsRepository) GetAll(ctx context.Context) (map[string]string, error) {
	query := `SELECT key, value FROM settings`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all settings: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings[key] = value
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating settings: %w", err)
	}

	return settings, nil
}
