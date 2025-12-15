package services

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/suipic/backend/config"
	"github.com/suipic/backend/models"
	"github.com/suipic/backend/repository"
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

	return service, nil
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

func (s *DatabaseService) CreatePhotographerClient(photographerID, clientID int64) (*models.PhotographerClient, error) {
	pc := &models.PhotographerClient{}
	query := `
		INSERT INTO photographer_clients (photographer_id, client_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, photographer_id, client_id, created_at
	`
	err := s.db.QueryRow(query, photographerID, clientID).Scan(
		&pc.ID, &pc.PhotographerID, &pc.ClientID, &pc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return pc, nil
}

func (s *DatabaseService) GetPhotographerClient(photographerID, clientID int64) (*models.PhotographerClient, error) {
	pc := &models.PhotographerClient{}
	query := `
		SELECT id, photographer_id, client_id, created_at
		FROM photographer_clients
		WHERE photographer_id = $1 AND client_id = $2
	`
	err := s.db.QueryRow(query, photographerID, clientID).Scan(
		&pc.ID, &pc.PhotographerID, &pc.ClientID, &pc.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return pc, nil
}

func (s *DatabaseService) GetClientsByPhotographer(photographerID int64) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.username, u.password_hash, u.role, u.created_at, u.updated_at
		FROM users u
		INNER JOIN photographer_clients pc ON u.id = pc.client_id
		WHERE pc.photographer_id = $1
		ORDER BY u.username
	`
	rows, err := s.db.Query(query, photographerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.PasswordHash,
			&user.Role, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		clients = append(clients, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (s *DatabaseService) Close() error {
	return s.db.Close()
}

func (s *DatabaseService) GetDB() *sql.DB {
	return s.db
}

func (s *DatabaseService) GetPhotoRepo() repository.PhotoRepository {
	return repository.NewPostgresPhotoRepository(s.db)
}

func (s *DatabaseService) GetCommentRepo() repository.CommentRepository {
	return repository.NewPostgresCommentRepository(s.db)
}
