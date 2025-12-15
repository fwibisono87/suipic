package services

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/suipic/backend/config"
	"github.com/suipic/backend/models"
)

type DatabaseService struct {
	db *sql.DB
}

func NewDatabaseService(cfg *config.DatabaseConfig) (*DatabaseService, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	service := &DatabaseService{db: db}

	if err := service.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return service, nil
}

func (s *DatabaseService) initSchema() error {
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'client',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

		CREATE TABLE IF NOT EXISTS albums (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			date_taken TIMESTAMP WITH TIME ZONE,
			description TEXT,
			location VARCHAR(500),
			custom_fields JSONB,
			thumbnail_photo_id INTEGER,
			photographer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_albums_photographer_id ON albums(photographer_id);
		CREATE INDEX IF NOT EXISTS idx_albums_date_taken ON albums(date_taken);

		CREATE TABLE IF NOT EXISTS album_permissions (
			id SERIAL PRIMARY KEY,
			album_id INTEGER NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			can_view BOOLEAN NOT NULL DEFAULT true,
			can_edit BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(album_id, user_id)
		);

		CREATE INDEX IF NOT EXISTS idx_album_permissions_album_id ON album_permissions(album_id);
		CREATE INDEX IF NOT EXISTS idx_album_permissions_user_id ON album_permissions(user_id);
	`

	_, err := s.db.Exec(schema)
	return err
}

func (s *DatabaseService) CreateUser(email, username, passwordHash string, role models.UserRole) (*models.User, error) {
	user := &models.User{}
	query := `
		INSERT INTO users (email, username, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, username, password_hash, role, created_at, updated_at
	`
	err := s.db.QueryRow(query, email, username, passwordHash, role).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *DatabaseService) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, role, created_at, updated_at
		FROM users WHERE email = $1
	`
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *DatabaseService) GetUserByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, role, created_at, updated_at
		FROM users WHERE id = $1
	`
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *DatabaseService) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, role, created_at, updated_at
		FROM users WHERE username = $1
	`
	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *DatabaseService) GetUsersByRole(role models.UserRole) ([]*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, role, created_at, updated_at
		FROM users WHERE role = $1
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.PasswordHash,
			&user.Role, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *DatabaseService) GetAlbumByID(id int64) (*models.Album, error) {
	album := &models.Album{}
	query := `
		SELECT id, title, date_taken, description, location, custom_fields, 
		       thumbnail_photo_id, photographer_id, created_at, updated_at
		FROM albums WHERE id = $1
	`
	err := s.db.QueryRow(query, id).Scan(
		&album.ID, &album.Title, &album.DateTaken, &album.Description, &album.Location,
		&album.CustomFields, &album.ThumbnailPhotoID, &album.PhotographerID,
		&album.CreatedAt, &album.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (s *DatabaseService) GetAlbumPermission(albumID, userID int64) (*models.AlbumPermission, error) {
	permission := &models.AlbumPermission{}
	query := `
		SELECT id, album_id, user_id, can_view, can_edit, created_at
		FROM album_permissions WHERE album_id = $1 AND user_id = $2
	`
	err := s.db.QueryRow(query, albumID, userID).Scan(
		&permission.ID, &permission.AlbumID, &permission.UserID,
		&permission.CanView, &permission.CanEdit, &permission.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (s *DatabaseService) Close() error {
	return s.db.Close()
}
