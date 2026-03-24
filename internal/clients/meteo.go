package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// var instead of const so tests can override with local httptest.Server URL
var meteoBaseURL = "https://api.meteo.com/v1"

// Used by daashboard, needs temperature and precipation numbers
type MeteoResult struct {
	Temperature float64
	Precipation float64
}

// GetWeather fetches mean forecast temperature and precipation for given coordinates
func GetWeather(latitude, longitude float64) (*MeteoResult, error) {
	url := fmt.Sprintf(
		"%s/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation",
		meteoBaseURL, latitude, longitude,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meteo API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		Hourly struct {
			Temperature   []float64 `json:"temperature_2m"`
			Precipitation []float64 `json:"precipitation"`
		} `json:"hourly"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &MeteoResult{
		Temperature: mean(response.Hourly.Temperature),
		Precipation: mean(response.Hourly.Precipitation),
	}, nil
}

// mean calculates the average of a float64 slice, returns 0 if empty
func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}
