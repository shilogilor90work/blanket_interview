package migrations

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Property struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PropertyID string    `gorm:"uniqueIndex;not null" json:"property_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Record struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PropertyID string    `gorm:"index;not null" json:"property_id"`
	Amount     float64   `gorm:"not null" json:"amount"`
	Type       string    `gorm:"type:transaction_type;not null" json:"type"`
	Date       time.Time `gorm:"not null;index" json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Property   Property  `gorm:"foreignKey:PropertyID;references:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

// Version 1: Create records table
func v1_createTableRecord(db *gorm.DB) error {
	err := db.Exec(`
		CREATE TYPE transaction_type AS ENUM ('income', 'expense');
	`).Error
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Property{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Record{})
	if err != nil {
		return err
	}
	log.Println("Applied migration v1: Create Property and Record tables")
	return nil
}
