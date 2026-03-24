package models

type ExchangeResponse struct {
	Country        string             `json:"country"`
	BaseCurrency   string             `json:"base-currency"`
	TargetCurrencies map[string]float64 `json:"targetCurrencies"`
}
