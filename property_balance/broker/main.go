package main

import (
	"broker/natsHandler"
	"broker/router"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get NATS URL from env or use default
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	// Initialize NATS connection
	err := natsHandler.InitNATS(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	log.Println("Connected to NATS at", natsURL)

	// Set up router
	r := router.SetupRouter()

	// Start HTTP server in a goroutine to allow for simultaneous NATS handling
	go func() {
		log.Println("Broker service listening on :80")
		if err := http.ListenAndServe(":80", r); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	// Now, continue handling NATS messages
	select {} // Blocks forever, allowing NATS and HTTP server to run concurrently
}
