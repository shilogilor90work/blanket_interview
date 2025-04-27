package models

import (
	"time"
)

type CreateProperty struct {
	PropertyID string `json:"property_id"`
}

type Property struct {
	ID         uint      `json:"id"`
	PropertyID string    `json:"property_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
