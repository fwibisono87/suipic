package models

import "time"

type AlbumUser struct {
	ID        int       `json:"id"`
	AlbumID   int       `json:"albumId"`
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
