package models

import "time"

type Comment struct {
	ID              int       `json:"id"`
	PhotoID         int       `json:"photoId"`
	UserID          int       `json:"userId"`
	ParentCommentID *int      `json:"parentCommentId,omitempty"`
	Text            string    `json:"text"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
