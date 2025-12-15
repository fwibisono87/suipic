package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Album struct {
	ID               int          `json:"id"`
	Title            string       `json:"title"`
	DateTaken        *time.Time   `json:"dateTaken,omitempty"`
	Description      *string      `json:"description,omitempty"`
	Location         *string      `json:"location,omitempty"`
	CustomFields     CustomFields `json:"customFields,omitempty"`
	ThumbnailPhotoID *int         `json:"thumbnailPhotoId,omitempty"`
	PhotographerID   int          `json:"photographerId"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
}

type CustomFields map[string]interface{}

func (c CustomFields) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *CustomFields) Scan(value interface{}) error {
	if value == nil {
		*c = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, c)
}
