package router

import (
	"broker/handler"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/properties", handler.AddProperty).Methods("POST")
	r.HandleFunc("/records", handler.AddRecord).Methods("POST")
	r.HandleFunc("/records", handler.GetRecords).Methods("GET")
	r.HandleFunc("/balance/{property_id}", handler.GetCurrentBalance).Methods("GET")
	r.HandleFunc("/balance/monthly/{property_id}", handler.GetMonthlyBalance).Methods("GET")

	return r
}
