package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCountry(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		mockBody    interface{}
		mockStatus  int
		wantName    string
		wantCode    string
		wantLat     float64
		wantLon     float64
		wantErr     bool
	}{
		{
			name:        "valid country",
			countryCode: "NO",
			mockBody: []map[string]interface{}{
				{
					"name":       map[string]interface{}{"common": "Norway"},
					"cca2":       "NO",
					"capital":    []interface{}{"Oslo"},
					"population": float64(5000000),
					"area":       float64(323802),
					"region":     "Europe",
					"latlng":     []interface{}{float64(62.0), float64(10.0)},
				},
			},
			mockStatus: http.StatusOK,
			wantName:   "Norway",
			wantCode:   "NO",
			wantLat:    62.0,
			wantLon:    10.0,
			wantErr:    false,
		},
		{
			name:        "country not found",
			countryCode: "XX",
			mockBody:    nil,
			mockStatus:  http.StatusNotFound,
			wantErr:     true,
		},
		{
			name:        "empty results",
			countryCode: "NO",
			mockBody:    []map[string]interface{}{},
			mockStatus:  http.StatusOK,
			wantErr:     true,
		},
		{
			name:        "API error",
			countryCode: "NO",
			mockBody:    nil,
			mockStatus:  http.StatusInternalServerError,
			wantErr:     true,
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

			restCountriesAPI = server.URL + "/"

			result, err := GetCountry(tt.countryCode)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
				return
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if result.Name != tt.wantName {
					t.Errorf("Name: got %s, want %s", result.Name, tt.wantName)
				}
				if result.Code != tt.wantCode {
					t.Errorf("Code: got %s, want %s", result.Code, tt.wantCode)
				}
				if result.Coordinates.Latitude != tt.wantLat {
					t.Errorf("Latitude: got %f, want %f", result.Coordinates.Latitude, tt.wantLat)
				}
				if result.Coordinates.Longitude != tt.wantLon {
					t.Errorf("Longitude: got %f, want %f", result.Coordinates.Longitude, tt.wantLon)
				}
			}
		})
	}
}
