package models

import "time"

type UserRole string

const (
	RoleAdmin        UserRole = "admin"
	RolePhotographer UserRole = "photographer"
	RoleClient       UserRole = "client"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	FriendlyName string    `json:"friendlyName"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
