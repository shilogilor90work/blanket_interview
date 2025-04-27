package models

import (
	"time"
)

type CreateRecord struct {
	PropertyID string  `json:"property_id"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type"`
}

type Record struct {
	ID         uint      `json:"id"`
	PropertyID string    `json:"property_id"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	Date       time.Time `json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
