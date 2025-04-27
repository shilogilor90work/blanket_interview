package models

import (
	"time"
)

type Record struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PropertyID string    `gorm:"index;not null" json:"property_id"`
	Amount     float64   `gorm:"not null" json:"amount"`
	Type       string    `json:"type"`
	Date       time.Time `gorm:"not null;index" json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Property   Property  `gorm:"foreignKey:PropertyID;references:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

type CreateRecord struct {
	PropertyID string  `gorm:"index;not null" json:"property_id"`
	Amount     float64 `gorm:"not null" json:"amount"`
	Type       string  `gorm:"type:enum('income','expense');not null" json:"type"`
}

type GetRecordsParams struct {
	PropertyID string `json:"property_id"`
	Type       string `json:"type"`
	From       string `json:"from"`
	To         string `json:"to"`
	Sort       string `json:"sort"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}
