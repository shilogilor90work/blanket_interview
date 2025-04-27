package handler

import (
	"broker/cache"
	"broker/models"
	"broker/natsHandler"
	"encoding/json"
	"net/http"
	"time"
)

var cacheInstance = cache.NewCache()

type QueryParamMapper[T models.HasPropertyID] func(r *http.Request) (T, error)

func buildCacheKey(subject string, params any) string {
	paramBytes, _ := json.Marshal(params) // Best effort, ignore error
	return subject + string(paramBytes)
}

func handleCreateRequest[T models.HasPropertyID, R any](w http.ResponseWriter, r *http.Request, subject string) {
	// no need for extractor
	var payload T
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	resp, err := natsHandler.Request(subject, payload, 10*time.Second)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	var result R
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	// Use payload.GetPropertyID()
	cacheInstance.Invalidate(payload.GetPropertyID())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Update handleGetRequest to allow custom response processing before sending it back
func handleGetRequest[T models.HasPropertyID, R any](w http.ResponseWriter, r *http.Request, subject string, mapper QueryParamMapper[T], customResponseHandler func(R) []byte) {
	params, err := mapper(r)
	if err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	propertyID := params.GetPropertyID()
	cacheKey := buildCacheKey(subject, params)

	if item, ok := cacheInstance.Get(propertyID, cacheKey); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(item.Data)
		return
	}

	resp, err := natsHandler.Request(subject, params, 10*time.Second)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	var result R
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}
	cache_data := resp.Data

	if customResponseHandler != nil {
		cache_data = customResponseHandler(result)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}

	cacheInstance.Set(propertyID, cacheKey, cache_data)

}
