package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAirQuality(t *testing.T) {
	tests := []struct {
		name       string
		mockBody   interface{}
		mockStatus int
		wantPM25   float64
		wantPM10   float64
		wantLevel  string
		wantErr    bool
	}{
		{
			name: "valid stations with good air quality",
			mockBody: map[string]interface{}{
				"results": []map[string]interface{}{
					{
						"sensors": []map[string]interface{}{
							{
								"parameter": map[string]interface{}{"id": 2},
								"latest":    map[string]interface{}{"value": 8.0},
							},
						},
					},
					{
						"sensors": []map[string]interface{}{
							{
								"parameter": map[string]interface{}{"id": 2},
								"latest":    map[string]interface{}{"value": 12.0},
							},
						},
					},
					{
						"sensors": []map[string]interface{}{
							{
								"parameter": map[string]interface{}{"id": 1},
								"latest":    map[string]interface{}{"value": 15.0},
							},
						},
					},
				},
			},
			mockStatus: http.StatusOK,
			wantPM25:   10.0, // mean of 8.0 and 12.0
			wantPM10:   15.0,
			wantLevel:  "Good",
			wantErr:    false,
		},
		{
			name: "no stations found returns unknown",
			mockBody: map[string]interface{}{
				"results": []map[string]interface{}{},
			},
			mockStatus: http.StatusOK,
			wantPM25:   -1,
			wantPM10:   -1,
			wantLevel:  "unknown",
			wantErr:    false,
		},
		{
			name: "moderate air quality level",
			mockBody: map[string]interface{}{
				"results": []map[string]interface{}{
					{
						"sensors": []map[string]interface{}{
							{
								"parameter": map[string]interface{}{"id": 2},
								"latest":    map[string]interface{}{"value": 20.0},
							},
						},
					},
				},
			},
			mockStatus: http.StatusOK,
			wantPM25:   20.0,
			wantLevel:  "Moderate",
			wantErr:    false,
		},
		{
			name:       "API error returns error",
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

			openaqBaseURL = server.URL

			result, err := GetAirQuality(59.9, 10.7)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
				return
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if result.PM25 != tt.wantPM25 {
					t.Errorf("PM25: got %f, want %f", result.PM25, tt.wantPM25)
				}
				if result.Level != tt.wantLevel {
					t.Errorf("Level: got %s, want %s", result.Level, tt.wantLevel)
				}
				if result.PM10 != tt.wantPM10 {
					t.Errorf("PM10: got %f, want %f", result.PM10, tt.wantPM10)
				}
			}
		})
	}
}

func TestAqiLevel(t *testing.T) {
	tests := []struct {
		name      string
		pm25      float64
		wantLevel string
	}{
		{"good", 5.0, "Good"},
		{"moderate", 20.0, "Moderate"},
		{"unhealthy for sensitive groups", 40.0, "Unhealthy for Sensitive Groups"},
		{"unhealthy", 100.0, "Unhealthy"},
		{"very unhealthy", 200.0, "Very Unhealthy"},
		{"hazardous", 300.0, "Hazardous"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aqiLevel(tt.pm25)
			if got != tt.wantLevel {
				t.Errorf("aqiLevel(%f): got %s, want %s", tt.pm25, got, tt.wantLevel)
			}
		})
	}
}
