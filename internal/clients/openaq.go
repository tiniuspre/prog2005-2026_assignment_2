package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var openaqBaseURL = "https://api.openaq.org/v3"

type AirQualityResult struct {
	PM25  float64 `json:"pm25"`
	Level string  `json:"level"`
}

type openaqResponse struct {
	Results []struct {
		Sensors []struct {
			Parameter struct {
				ID int `json:"id"`
			} `json:"parameter"`
			Latest struct {
				Value float64 `json:"value"`
			} `json:"latest"`
		} `json:"sensors"`
	} `json:"results"`
}

// GetAirQuality fetches mean PM2.5 readings from stations within 50km of the given coordinates
func GetAirQuality(latitude, longitude float64) (*AirQualityResult, error) {
	url := fmt.Sprintf(
		"%s/locations?coordinates=%f,%f&radius=50000&parameters_id=2&limit=100",
		openaqBaseURL, latitude, longitude,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", os.Getenv("OPENAQ_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch air quality data: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openaq API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response openaqResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Collect PM2.5 readings (parameter ID 2) across all stations
	var readings []float64
	for _, location := range response.Results {
		for _, sensor := range location.Sensors {
			if sensor.Parameter.ID == 2 && sensor.Latest.Value > 0 {
				readings = append(readings, sensor.Latest.Value)
			}
		}
	}

	// No stations found within range
	if len(readings) == 0 {
		return &AirQualityResult{
			PM25:  -1,
			Level: "unknown",
		}, nil
	}

	pm25 := mean(readings)

	return &AirQualityResult{
		PM25:  pm25,
		Level: aqiLevel(pm25),
	}, nil
}

// aqiLevel returns an EPA AQI level string based on PM2.5 concentration
func aqiLevel(pm25 float64) string {
	switch {
	case pm25 <= 12.0:
		return "Good"
	case pm25 <= 35.4:
		return "Moderate"
	case pm25 <= 55.4:
		return "Unhealthy for Sensitive Groups"
	case pm25 <= 150.4:
		return "Unhealthy"
	case pm25 <= 250.4:
		return "Very Unhealthy"
	default:
		return "Hazardous"
	}
}
