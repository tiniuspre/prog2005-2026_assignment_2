package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExchangeRates(t *testing.T) {
	tests := []struct {
		name             string
		baseCurrency     string
		targetCurrencies []string
		mockBody         interface{}
		mockStatus       int
		wantRates        map[string]float64
		wantErr          bool
	}{
		{
			name:             "valid exchange rates",
			baseCurrency:     "NOK",
			targetCurrencies: []string{"EUR", "USD"},
			mockBody: map[string]interface{}{
				"rates": map[string]interface{}{
					"EUR": 0.087,
					"USD": 0.095,
					"SEK": 0.978,
				},
			},
			mockStatus: http.StatusOK,
			wantRates:  map[string]float64{"EUR": 0.087, "USD": 0.095},
			wantErr:    false,
		},
		{
			name:             "target currency not in response",
			baseCurrency:     "NOK",
			targetCurrencies: []string{"XYZ"},
			mockBody: map[string]interface{}{
				"rates": map[string]interface{}{
					"EUR": 0.087,
				},
			},
			mockStatus: http.StatusOK,
			wantRates:  map[string]float64{},
			wantErr:    false,
		},
		{
			name:             "API error",
			baseCurrency:     "NOK",
			targetCurrencies: []string{"EUR"},
			mockBody:         nil,
			mockStatus:       http.StatusInternalServerError,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				if tt.mockBody != nil {
					json.NewEncoder(w).Encode(tt.mockBody)
				}
			}))
			defer server.Close()

			currencyAPI = server.URL + "/"

			result, err := GetExchangeRates(tt.baseCurrency, tt.targetCurrencies)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
				return
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				for currency, wantRate := range tt.wantRates {
					if gotRate, ok := result[currency]; !ok {
						t.Errorf("missing currency %s in result", currency)
					} else if gotRate != wantRate {
						t.Errorf("rate for %s: got %f, want %f", currency, gotRate, wantRate)
					}
				}
			}
		})
	}
}
