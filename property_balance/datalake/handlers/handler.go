package handler

import (
	"datalake/models"
	"datalake/natsHandler"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type CreateRecord struct {
	PropertyID string  `gorm:"index;not null" json:"property_id"`
	Amount     float64 `gorm:"not null" json:"amount"`
	Type       string  `gorm:"type:enum('income','expense');not null" json:"type"`
}

type CreateProperty struct {
	PropertyID string `gorm:"index;not null" json:"property_id"`
}

func ListenForCreateProperty(db *gorm.DB) {
	// Subscribe to the NATS subject
	natsHandler.Subscribe("datalake.property.add", func(msg *nats.Msg) {
		var createProperty CreateProperty
		var property models.Property

		// Unmarshal the incoming message into the CreatePropery struct
		if err := json.Unmarshal(msg.Data, &createProperty); err != nil {
			log.Printf("Error unmarshalling createProperty: %v", err)
			// Send error response back to NATS
			errorResponse := map[string]string{"error": "Error unmarshalling createProperty"}
			natsHandler.Publish(msg.Reply, errorResponse)
			return
		}

		// Create the actual property
		property = models.Property{
			PropertyID: createProperty.PropertyID,
		}

		// Add the property to the database
		if err := db.Create(&property).Error; err != nil {
			log.Printf("Error adding property to database: %v", err)
			// Send error response back to NATS
			errorResponse := map[string]string{"error": err.Error()}
			natsHandler.Publish(msg.Reply, errorResponse)
			return
		}

		log.Printf("Record added: %v", property)

		// Send success response back to NATS with the created record
		natsHandler.Publish(msg.Reply, property)
	})
}

func ListenForCreateRecord(db *gorm.DB) {
	// Subscribe to the NATS subject
	natsHandler.Subscribe("datalake.record.add", func(msg *nats.Msg) {
		var createRecord CreateRecord
		var record models.Record

		// Unmarshal the incoming message into the CreateRecord struct
		if err := json.Unmarshal(msg.Data, &createRecord); err != nil {
			log.Printf("Error unmarshalling createRecord: %v", err)
			// Send error response back to NATS
			errorResponse := map[string]string{"error": "Error unmarshalling createRecord"}
			natsHandler.Publish(msg.Reply, errorResponse)
			return
		}

		// Create the actual record
		record = models.Record{
			PropertyID: createRecord.PropertyID,
			Amount:     createRecord.Amount,
			Type:       createRecord.Type,
			Date:       time.Now(), // Set current time for Date (or parse it if needed)
		}

		// Add the record to the database
		if err := db.Create(&record).Error; err != nil {
			log.Printf("Error adding record to database: %v", err)
			// Send error response back to NATS
			errorResponse := map[string]string{"error": err.Error()}
			natsHandler.Publish(msg.Reply, errorResponse)
			return
		}

		log.Printf("Record added: %v", record)

		// Send success response back to NATS with the created record
		natsHandler.Publish(msg.Reply, record)
	})
}

// Listen for balance-related requests
func ListenForBalanceRequests(db *gorm.DB) {
	natsHandler.Subscribe("datalake.balance.current", func(msg *nats.Msg) {
		var request map[string]string
		if err := json.Unmarshal(msg.Data, &request); err != nil {
			log.Printf("Error unmarshalling balance request: %v", err)
			return
		}

		propertyID := request["property_id"]
		log.Printf("Received current balance request for property %s", propertyID)

		// Query the current balance (sum of incomes - sum of expenses)
		var income, expense float64
		db.Model(&models.Record{}).Where("property_id = ? AND type = ?", propertyID, "income").Select("sum(amount)").Scan(&income)
		db.Model(&models.Record{}).Where("property_id = ? AND type = ?", propertyID, "expense").Select("sum(amount)").Scan(&expense)

		balance := income - expense
		natsHandler.Publish("datalake.balance.current", map[string]interface{}{
			"property_id": propertyID,
			"balance":     balance,
		})
	})
}
