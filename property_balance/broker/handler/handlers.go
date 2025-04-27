package handler

import (
	"broker/models"
	"broker/natsHandler"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func handleCreateRequest[T any, R any](w http.ResponseWriter, r *http.Request, subject string) {
	var payload T
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Send request to NATS and wait for response
	resp, err := natsHandler.Request(subject, payload, 10*time.Second)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	// Create a variable of type R to unmarshal the response into
	var result R
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	// Return the unmarshalled response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func AddProperty(w http.ResponseWriter, r *http.Request) {
	handleCreateRequest[models.CreateProperty, models.Property](w, r, "configurator.property.add")
}

func AddRecord(w http.ResponseWriter, r *http.Request) {
	handleCreateRequest[models.CreateRecord, models.Record](w, r, "configurator.record.add")
}

func GetRecords(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	message := map[string]string{
		"property_id": query.Get("property_id"),
		"type":        query.Get("type"),
		"from":        query.Get("from"),
		"to":          query.Get("to"),
		"sort":        query.Get("sort"),
		"page":        query.Get("page"),
		"limit":       query.Get("limit"),
	}
	_ = natsHandler.Publish("datalake.record.query", message)
	w.WriteHeader(http.StatusAccepted)
}

func GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	propertyID := mux.Vars(r)["property_id"]
	_ = natsHandler.Publish("datalake.balance.current", map[string]string{
		"property_id": propertyID,
	})
	w.WriteHeader(http.StatusAccepted)
}

func GetMonthlyBalance(w http.ResponseWriter, r *http.Request) {
	propertyID := mux.Vars(r)["property_id"]

	// Create request payload
	payload := map[string]string{
		"property_id": propertyID,
	}

	// Marshal payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode request", http.StatusInternalServerError)
		return
	}

	// Send request to NATS and wait for reply
	msg, err := natsHandler.Request("datalake.balance.monthly", data, 2*time.Second)
	if err != nil {
		http.Error(w, "Error waiting for datalake response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response from datalake to HTTP client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg.Data)
}
