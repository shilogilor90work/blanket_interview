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
