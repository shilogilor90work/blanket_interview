package main

import (
	"log"
	"os"

	"configurator/handler"
	"configurator/natsHandler"
)

func main() {
	// Get NATS URL from environment or use default
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	// Initialize NATS connection
	natsHandler.InitNATS(natsURL)
	log.Println("Connected to NATS at", natsURL)

	// Subscribe to NATS subjects
	go handler.ListenForCreateProperty()
	go handler.ListenForCreateRecord()

	// Block forever (to keep the application running)
	select {}
}
