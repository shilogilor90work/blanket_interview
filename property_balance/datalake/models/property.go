package models

import (
	"time"
)

type Property struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PropertyID string    `gorm:"uniqueIndex;not null" json:"property_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateProperty struct {
	PropertyID string `gorm:"index;not null" json:"property_id"`
}
