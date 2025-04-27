package handler

import (
	"configurator/models"
	"configurator/natsHandler"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func ListenForCreateProperty() {
	natsHandler.Subscribe("configurator.property.add", func(msg *nats.Msg) {
		var createProperty models.CreateProperty
		if err := json.Unmarshal(msg.Data, &createProperty); err != nil {
			log.Printf("Error unmarshalling property: %v", err)
			return
		}

		// Process createProperty (for now, just log)
		log.Printf("Received createProperty: %v", createProperty)

		// Send a request to Datalake service and wait for a response
		response, err := natsHandler.Request("datalake.property.add", createProperty, 10*time.Second)
		if err != nil {
			log.Printf("Error requesting property from Datalake: %v", err)
			// Respond back to the sender with an error message
			if err := msg.Respond([]byte("Error processing property in Datalake")); err != nil {
				log.Printf("Error responding to sender: %v", err)
			}
			return
		}

		// Log the response received from the Datalake service
		log.Printf("Received response from Datalake: %s", response.Data)

		// Respond back to the sender with the response from Datalake
		if err := msg.Respond(response.Data); err != nil {
			log.Printf("Error responding to sender: %v", err)
		}
	})
}

func ListenForCreateRecord() {
	natsHandler.Subscribe("configurator.record.add", func(msg *nats.Msg) {
		var createRecord models.CreateRecord
		if err := json.Unmarshal(msg.Data, &createRecord); err != nil {
			log.Printf("Error unmarshalling record: %v", err)
			return
		}

		// Process createRecord (for now, just log)
		log.Printf("Received createRecord: %v", createRecord)

		// Send a request to Datalake service and wait for a response
		response, err := natsHandler.Request("datalake.record.add", createRecord, 10*time.Second)
		if err != nil {
			log.Printf("Error requesting record from Datalake: %v", err)
			// Respond back to the sender with an error message
			if err := msg.Respond([]byte("Error processing record in Datalake")); err != nil {
				log.Printf("Error responding to sender: %v", err)
			}
			return
		}

		// Log the response received from the Datalake service
		log.Printf("Received response from Datalake: %s", response.Data)

		// Respond back to the sender with the response from Datalake
		if err := msg.Respond(response.Data); err != nil {
			log.Printf("Error responding to sender: %v", err)
		}
	})
}

// Listen for balance-related requests
func ListenForBalanceRequests() {
	natsHandler.Subscribe("configurator.balance.current", func(msg *nats.Msg) {
		var request map[string]string
		if err := json.Unmarshal(msg.Data, &request); err != nil {
			log.Printf("Error unmarshalling balance request: %v", err)
			return
		}

		propertyID := request["property_id"]
		log.Printf("Received current balance request for property %s", propertyID)

		// Query the current balance (simulate for now)
		balance := 500.0 // Example balance
		// Publish the balance to Datalake using natsHandler.Publish
		if err := natsHandler.Publish("datalake.balance.current", map[string]interface{}{
			"property_id": propertyID,
			"balance":     balance,
		}); err != nil {
			log.Printf("Error publishing balance to Datalake: %v", err)
		}
	})

	// Similar subscription logic for "property.balance.monthly" can go here
}
