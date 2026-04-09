package clients

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeather(t *testing.T) {
	tests := []struct {
		name       string
		mockBody   interface{}
		mockStatus int
		wantTemp   float64
		wantPrecip float64
		wantErr    bool
	}{
		{
			name: "valid weather data",
			mockBody: map[string]interface{}{
				"hourly": map[string]interface{}{
					"temperature_2m": []float64{2.0, 4.0, 6.0},
					"precipitation":  []float64{0.0, 0.2, 0.4},
				},
			},
			mockStatus: http.StatusOK,
			wantTemp:   4.0, // mean of 2, 4, 6
			wantPrecip: 0.2, // mean of 0, 0.2, 0.4
			wantErr:    false,
		},
		{
			name: "empty hourly arrays returns zero",
			mockBody: map[string]interface{}{
				"hourly": map[string]interface{}{
					"temperature_2m": []float64{},
					"precipitation":  []float64{},
				},
			},
			mockStatus: http.StatusOK,
			wantTemp:   0.0,
			wantPrecip: 0.0,
			wantErr:    false,
		},
		{
			name:       "API error",
			mockBody:   nil,
			mockStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				if tt.mockBody != nil {
					_ = json.NewEncoder(w).Encode(tt.mockBody)
				}
			}))
			defer server.Close()

			meteoBaseURL = server.URL

			result, err := GetWeather(62.0, 10.0)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
				return
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if result.Temperature != tt.wantTemp {
					t.Errorf("Temperature: got %f, want %f", result.Temperature, tt.wantTemp)
				}
				if math.Abs(result.Precipitation-tt.wantPrecip) > 1e-9 {
					t.Errorf("Precipitation: got %f, want %f", result.Precipitation, tt.wantPrecip)
				}
			}
		})
	}
}
