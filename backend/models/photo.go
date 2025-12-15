package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type PickRejectState string

const (
	PickRejectNone   PickRejectState = "none"
	PickRejectPick   PickRejectState = "pick"
	PickRejectReject PickRejectState = "reject"
)

type Photo struct {
	ID              int             `json:"id"`
	AlbumID         int             `json:"albumId"`
	Filename        string          `json:"filename"`
	Title           *string         `json:"title,omitempty"`
	DateTime        *time.Time      `json:"dateTime,omitempty"`
	ExifData        ExifData        `json:"exifData,omitempty"`
	PickRejectState PickRejectState `json:"pickRejectState"`
	Stars           int             `json:"stars"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type ExifData map[string]interface{}

func (e ExifData) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	return json.Marshal(e)
}

func (e *ExifData) Scan(value interface{}) error {
	if value == nil {
		*e = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, e)
}
