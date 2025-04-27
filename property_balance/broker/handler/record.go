package handler

import (
	"broker/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AddRecord(w http.ResponseWriter, r *http.Request) {
	handleCreateRequest[models.CreateRecord, models.Record](w, r, "configurator.record.add")
}

func mapGetRecordsParams(r *http.Request) (models.GetRecordsParams, error) {
	// Extract path variables (e.g., property_id)
	vars := mux.Vars(r)
	propertyID, exists := vars["property_id"]
	if !exists {
		return models.GetRecordsParams{}, fmt.Errorf("missing property_id in URL")
	}

	// Manually map query parameters to struct fields
	queryParams := r.URL.Query()

	params := models.GetRecordsParams{
		PropertyID: propertyID,
		Type:       queryParams.Get("type"),
		From:       queryParams.Get("from"),
		To:         queryParams.Get("to"),
		Sort:       queryParams.Get("sort"),
	}

	// Handle page and limit (convert to integers if necessary)
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		page = 1 // Default value
	}
	params.Page = page

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		limit = 10 // Default value
	}
	params.Limit = limit

	return params, nil
}

func GetRecords(w http.ResponseWriter, r *http.Request) {
	// Use the handleGetRequest function with a custom query param mapper
	handleGetRequest[models.GetRecordsParams, interface{}](w, r, "get_records_subject", mapGetRecordsParams, nil)
}
