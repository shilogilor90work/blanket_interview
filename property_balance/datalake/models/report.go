package models

type Report struct {
	PropertyID      string   `json:"property_id"`
	StartingBalance float64  `json:"starting_balance"`
	Records         []Record `json:"records"`
}
