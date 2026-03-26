package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCapitalCoordinates(t *testing.T) {
	tests := []struct {
		name       string
		capital    string
		country    string
		mockBody   interface{}
		mockStatus int
		wantLat    float64
		wantLon    float64
		wantErr    bool
	}{
		{
			name:    "valid capital",
			capital: "Oslo",
			country: "Norway",
			mockBody: []map[string]string{
				{"lat": "59.913868", "lon": "10.752245"},
			},
			mockStatus: http.StatusOK,
			wantLat:    59.913868,
			wantLon:    10.752245,
			wantErr:    false,
		},
		{
			name:       "no results found",
			capital:    "Fakecity",
			country:    "Fakeland",
			mockBody:   []map[string]string{},
			mockStatus: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Spin up a fake Nominatim server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				json.NewEncoder(w).Encode(tt.mockBody)
			}))
			defer server.Close()

			// Override the base URL to point at our fake server
			nominatimBaseURL = server.URL

			coords, err := GetCapitalCoordinates(tt.capital, tt.country)
			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr {
				if coords.Latitude != tt.wantLat || coords.Longitude != tt.wantLon {
					t.Errorf("got (%f, %f), want (%f, %f)", coords.Latitude, coords.Longitude, tt.wantLat, tt.wantLon)
				}
			}
		})
	}
}
