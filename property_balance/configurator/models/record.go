package models

import (
	"time"
)

type CreateRecord struct {
	PropertyID string  `gorm:"index;not null" json:"property_id"`
	Amount     float64 `gorm:"not null" json:"amount"`
	Type       string  `gorm:"type:enum('income','expense');not null" json:"type"`
}

type Record struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PropertyID string    `gorm:"index;not null" json:"property_id"`
	Amount     float64   `gorm:"not null" json:"amount"`
	Type       string    `gorm:"type:enum('income','expense');not null" json:"type"`
	Date       time.Time `gorm:"not null;index" json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
