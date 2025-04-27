package models

import (
	"time"
)

type HasPropertyID interface {
	GetPropertyID() string
}

type CreateProperty struct {
	PropertyID string `json:"property_id"`
}

func (m CreateProperty) GetPropertyID() string {
	return m.PropertyID
}

type Property struct {
	ID         uint      `json:"id"`
	PropertyID string    `json:"property_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
