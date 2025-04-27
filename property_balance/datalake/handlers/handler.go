package handler

import (
	"datalake/models"
	"datalake/natsHandler"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

func ListenForCreateProperty(db *gorm.DB) {
	// Subscribe to the NATS subject
	natsHandler.Subscribe("datalake.property.add", func(msg *nats.Msg) {
		var createProperty models.CreateProperty
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
		var createRecord models.CreateRecord
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

func ListenForGetRecords(db *gorm.DB) {
	// Subscribe to the NATS subject
	natsHandler.Subscribe("datalake.record.get", func(msg *nats.Msg) {
		var params models.GetRecordsParams

		// Unmarshal the incoming message into the GetRecordsParams struct
		if err := json.Unmarshal(msg.Data, &params); err != nil {
			log.Printf("Error unmarshalling params: %v", err)
			errorResponse := map[string]string{"error": "Error unmarshalling parameters"}
			natsHandler.Publish(msg.Reply, errorResponse)
			return
		}

		// Call GetRecords to fetch the records based on the params
		records, err := GetRecords(db, params)
		if err != nil {
			log.Printf("Error fetching records: %v", err)
			errorResponse := map[string]string{"error": err.Error()}
			natsHandler.Publish(msg.Reply, errorResponse)
			return
		}

		// Send the fetched records back to NATS
		natsHandler.Publish(msg.Reply, records)
	})
}
func GetRecords(db *gorm.DB, params models.GetRecordsParams) ([]models.Record, error) {
	// Build the query
	query := db.Model(&models.Record{})
	log.Printf("fucker 1: %v", params)
	// Apply filters based on the parameters
	if params.PropertyID != "" {
		query = query.Where("property_id = ?", params.PropertyID)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.From != "" {
		fromDate, err := time.Parse("2006-01-02", params.From)
		if err != nil {
			return nil, fmt.Errorf("invalid 'from' date format: %v", err)
		}
		query = query.Where("date >= ?", fromDate)
	}
	if params.To != "" {
		toDate, err := time.Parse("2006-01-02", params.To)
		if err != nil {
			return nil, fmt.Errorf("invalid 'to' date format: %v", err)
		}
		query = query.Where("date <= ?", toDate)
	}

	// Apply sorting if specified
	if params.Sort != "" {
		query = query.Order("date " + params.Sort)
	}

	// Implement pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute the query and get the results
	var records []models.Record
	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("error fetching records: %v", err)
	}

	return records, nil
}

// Listen for balance-related requests
func ListenForBalanceRequests(db *gorm.DB) {
	natsHandler.Subscribe("datalake.balance.current", func(msg *nats.Msg) {
		var request models.BalanceParams
		if err := json.Unmarshal(msg.Data, &request); err != nil {
			log.Printf("Error unmarshalling balance request: %v", err)
			return
		}

		log.Printf("Received current balance request for property %s", request.PropertyID)

		// Query the current balance (sum of incomes - sum of expenses)
		var balance float64
		db.Model(&models.Record{}).Where("property_id = ?", request.PropertyID).
			Select("COALESCE(sum(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) - COALESCE(sum(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)").
			Scan(&balance)

		natsHandler.Publish(msg.Reply, models.BalanceResponse{
			PropertyID: request.PropertyID,
			Balance:    balance,
		})
	})
}

func ListenForMonthlyBalanceRequests(db *gorm.DB) {
	natsHandler.Subscribe("datalake.balance.monthly", func(msg *nats.Msg) {
		var params models.MonthlyBalanceParams
		if err := json.Unmarshal(msg.Data, &params); err != nil {
			log.Printf("Error unmarshalling balance request: %v", err)
			return
		}

		log.Printf("Received monthly balance request for property %s, month %s", params.PropertyID, params.YearMonth)

		// Parse the year and month
		year, month, err := parseYearMonth(params.YearMonth)
		if err != nil {
			log.Printf("Error parsing year and month: %v", err)
			return
		}

		// Calculate starting balance up to the start of the requested month
		var startingBalance float64
		db.Model(&models.Record{}).Where("property_id = ? AND date < ?", params.PropertyID, time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)).
			Select("COALESCE(sum(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) - COALESCE(sum(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)").Scan(&startingBalance)

		// Get records for the selected month
		var records []models.Record
		startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0) // Next month's start

		if err := db.Where("property_id = ? AND date >= ? AND date < ?", params.PropertyID, startOfMonth, endOfMonth).Find(&records).Error; err != nil {
			log.Printf("Error fetching records for month %s: %v", params.YearMonth, err)
			return
		}

		// Send the report with the current balance and the records
		natsHandler.Publish(msg.Reply, models.Report{
			PropertyID:      params.PropertyID,
			StartingBalance: startingBalance,
			Records:         records,
		})
	})
}

// Helper function to parse "YYYY-MM" format
func parseYearMonth(yearMonth string) (int, int, error) {
	yearMonthSplit := strings.Split(yearMonth, "-")
	if len(yearMonthSplit) != 2 {
		return 0, 0, fmt.Errorf("invalid year-month format")
	}
	year, err := strconv.Atoi(yearMonthSplit[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year format: %v", err)
	}
	month, err := strconv.Atoi(yearMonthSplit[1])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("invalid month format: %v", err)
	}
	return year, month, nil
}
