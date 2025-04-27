package handler

import (
	"broker/natsHandler"
	"encoding/json"
	"net/http"
	"time"
)

type QueryParamMapper[T any] func(r *http.Request) (T, error)

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
		http.Error(w, "Failed to parse response", http.StatusInternalServerError) // should be handled and not showing to client server error
		return
	}

	// Return the unmarshalled response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Update handleGetRequest to allow custom response processing before sending it back
func handleGetRequest[T any, R any](w http.ResponseWriter, r *http.Request, subject string, mapper QueryParamMapper[T], customResponseHandler func(R)) {
	// Extract query parameters and map them to the struct
	params, err := mapper(r)
	if err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	// Send the request to NATS and wait for the response
	resp, err := natsHandler.Request(subject, params, 10*time.Second)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	// Unmarshal the response into a result of type R
	var result R
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	// Custom response handler (if provided) to modify the result before sending
	if customResponseHandler != nil {
		customResponseHandler(result)
	} else {
		// Return the response to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}

}
