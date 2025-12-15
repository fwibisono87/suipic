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

type PhotographerClient struct {
	ID             int64     `json:"id"`
	PhotographerID int64     `json:"photographer_id"`
	ClientID       int64     `json:"client_id"`
	CreatedAt      time.Time `json:"created_at"`
}
