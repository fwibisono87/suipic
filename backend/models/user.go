package models

import "time"

type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RolePhotographer UserRole = "photographer"
	RoleClient       UserRole = "client"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	FriendlyName string    `json:"friendlyName"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type LegacyAlbum struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     int64     `json:"owner_id"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AlbumPermission struct {
	ID        int64     `json:"id"`
	AlbumID   int64     `json:"album_id"`
	UserID    int64     `json:"user_id"`
	CanView   bool      `json:"can_view"`
	CanEdit   bool      `json:"can_edit"`
	CreatedAt time.Time `json:"created_at"`
}

type PhotographerClient struct {
	ID             int64     `json:"id"`
	PhotographerID int64     `json:"photographer_id"`
	ClientID       int64     `json:"client_id"`
	CreatedAt      time.Time `json:"created_at"`
}
