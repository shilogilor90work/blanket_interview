package migrations

import (
	"log"

	"gorm.io/gorm"
)

// Migrate applies the migrations to the database
func Migrate(db *gorm.DB) {
	// Get current version of the schema from a table or config (use version tracking)
	log.Println("Starting migrations")

	// Apply migrations in order
	if err := v1_createTableRecord(db); err != nil {
		log.Fatalf("Migration v1 failed: %v", err)
	}

	log.Println("Migrations completed.")
}
