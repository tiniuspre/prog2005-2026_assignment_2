package models

type ThresholdDetails struct {
	Field         string  `json:"field"`
	Operator      string  `json:"operator"`
	Threshold     float64 `json:"threshold"`
	MeasuredValue float64 `json:"measuredValue"`
}

type ThresholdConfig struct {
	Field    string  `json:"field"`
	Operator string  `json:"operator"`
	Value    float64 `json:"value"`
}

type NotificationRegistration struct {
	ID        string           `json:"id"`
	URL       string           `json:"url"`
	Country   string           `json:"country"`
	Event     string           `json:"event"`
	Threshold *ThresholdConfig `json:"threshold,omitempty"` // only for THRESHOLD events
}

type WebhookPayload struct {
	ID      string            `json:"id"`
	Country string            `json:"country"`
	Event   string            `json:"event"`
	Time    string            `json:"time"`
	Details *ThresholdDetails `json:"details,omitempty"` // only for THRESHOLD events
}
