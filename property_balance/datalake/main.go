package main

import (
	handler "datalake/handlers"
	"datalake/migrations"
	"datalake/natsHandler"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const retryLimit = 30 * time.Second

var db *gorm.DB
var err error

func main() {
	// Get NATS URL and Postgres connection string from environment
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	// Initialize NATS connection
	natsHandler.InitNATS(natsURL)
	log.Println("Connected to NATS at", natsURL)

	// Initialize PostgreSQL connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=datalake port=5432 sslmode=disable"
	}
	startTime := time.Now()

	for {
		// Try to open the database connection
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to PostgreSQL")
			break
		}

		// Check if the retry time limit has been exceeded
		if time.Since(startTime) > retryLimit {
			log.Fatalf("Failed to connect to PostgreSQL within the time limit: %v", err)
		}

		// Log the error and retry after a delay
		log.Printf("Error connecting to PostgreSQL: %v. Retrying...", err)
		time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
	}

	// Run migrations
	migrations.Migrate(db)

	// Subscribe to NATS subjects
	go handler.ListenForCreateProperty(db)
	go handler.ListenForCreateRecord(db)
	go handler.ListenForGetRecords(db)
	go handler.ListenForBalanceRequests(db)
	go handler.ListenForMonthlyBalanceRequests(db)

	// Block forever
	select {}
}
