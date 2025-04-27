package models

type BalanceResponse struct {
	PropertyID string  `json:"property_id"`
	Balance    float64 `json:"balance"`
}

type BalanceParams struct {
	PropertyID string `json:"property_id"`
}

func (m BalanceParams) GetPropertyID() string {
	return m.PropertyID
}

type MonthlyBalanceParams struct {
	PropertyID string `json:"property_id"`
	YearMonth  string `json:"year_month"`
}

func (m MonthlyBalanceParams) GetPropertyID() string {
	return m.PropertyID
}
