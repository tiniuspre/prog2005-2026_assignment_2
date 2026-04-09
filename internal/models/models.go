package models

// THRESHOLD
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

// REGISTRATION
type Features struct {
	Temperature      bool     `json:"temperature"`
	Precipitation    bool     `json:"precipitation"`
	AirQuality       bool     `json:"airQuality"`
	Capital          bool     `json:"capital"`
	Coordinates      bool     `json:"coordinates"`
	Population       bool     `json:"population"`
	Area             bool     `json:"area"`
	TargetCurrencies []string `json:"targetCurrencies"`
}

type Registration struct {
	ID         string   `json:"id"`
	Country    string   `json:"country"`
	IsoCode    string   `json:"isoCode"`
	Features   Features `json:"features"`
	LastChange string   `json:"lastChange"`
}

// NOTIFICATIONS
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

// DASHBOARD
type AirQualityData struct {
	PM25  float64 `json:"pm25"`
	PM10  float64 `json:"pm10"`
	Level string  `json:"level"`
}

type DashboardFeatures struct {
	Temperature      *float64           `json:"temperature,omitempty"`
	Precipitation    *float64           `json:"precipitation,omitempty"`
	AirQuality       *AirQualityData    `json:"airQuality,omitempty"`
	Capital          *string            `json:"capital,omitempty"`
	Coordinates      *Coordinates       `json:"coordinates,omitempty"`
	Population       *int64             `json:"population,omitempty"`
	Area             *float64           `json:"area,omitempty"`
	TargetCurrencies map[string]float64 `json:"targetCurrencies,omitempty"`
}

type DashboardResponse struct {
	Country       string            `json:"country"`
	IsoCode       string            `json:"isoCode"`
	Features      DashboardFeatures `json:"features"`
	LastRetrieval string            `json:"lastRetrieval"`
}
