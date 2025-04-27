package models

type BalanceResponse struct {
	PropertyID string  `json:"property_id"`
	Balance    float64 `json:"balance"`
}

type BalanceParams struct {
	PropertyID string `json:"property_id"`
}

type MonthlyBalanceParams struct {
	PropertyID string `json:"property_id"`
	YearMonth  string `json:"year_month"`
}
