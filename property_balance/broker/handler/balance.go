package handler

import (
	"broker/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Define a function to map the query parameters to the expected input for NATS
func balanceQueryParamMapper(r *http.Request) (models.BalanceParams, error) {
	vars := mux.Vars(r)
	propertyID, exists := vars["property_id"]
	if !exists {
		return models.BalanceParams{}, fmt.Errorf("missing property_id in URL")
	}
	return models.BalanceParams{
		PropertyID: propertyID,
	}, nil
}

// Refactor GetCurrentBalance to use handleGetRequest
func GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	handleGetRequest[models.BalanceParams, models.BalanceResponse](
		w, r,
		"datalake.balance.current",
		balanceQueryParamMapper,
		nil,
	)
}

// Define the mapper function for the monthly balance query parameters
func monthlyBalanceQueryParamMapper(r *http.Request) (models.MonthlyBalanceParams, error) {
	vars := mux.Vars(r)
	propertyID, exists := vars["property_id"]
	if !exists {
		return models.MonthlyBalanceParams{}, fmt.Errorf("missing property_id in URL")
	}
	queryParams := r.URL.Query()
	month := queryParams.Get("year_month")

	if month == "" {
		// Default to current month if no month is provided
		month = time.Now().Format("2006-01") // Format as YYYY-MM
	}

	return models.MonthlyBalanceParams{
		PropertyID: propertyID,
		YearMonth:  month,
	}, nil
}

// Define a function to process the report and format it as a string
func processReport(report models.Report) string {
	balance := report.StartingBalance
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Starting cash = $%.2f\n", balance))

	// Loop over the records and calculate the balance
	for i, record := range report.Records {
		balanceChange := record.Amount
		if record.Type == "expense" {
			balance -= balanceChange
		} else {
			balance += balanceChange
		}

		result.WriteString(fmt.Sprintf("Record %d => type=%s, amount=$%.2f, $%.2f\n",
			i+1, record.Type, record.Amount, balance))
	}

	result.WriteString(fmt.Sprintf("Ending cash = $%.2f\n", balance))

	return result.String()
}

// Handle the result of the monthly balance query
func GetMonthlyBalance(w http.ResponseWriter, r *http.Request) {
	// Custom response handler to process the result before sending
	customResponseHandler := func(report models.Report) []byte {
		// Process and format the report response into the desired string
		formattedReport := processReport(report)
		// Return the formatted response to the client
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(formattedReport))
		return []byte(formattedReport)
	}

	// Call handleGetRequest with the custom response handler
	handleGetRequest[models.MonthlyBalanceParams, models.Report](
		w, r,
		"datalake.balance.monthly",
		monthlyBalanceQueryParamMapper,
		customResponseHandler,
	)
}
