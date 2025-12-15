package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/suipic/backend/config"
)

var DB *sql.DB

func Connect(cfg *config.DatabaseConfig) error {
	connStr := cfg.ConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
