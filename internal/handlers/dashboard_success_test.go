package handlers

import (
	"assignment_2/internal/clients"
	"assignment_2/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDashboardHandler_Success(t *testing.T) {
	origCountryByISO := countryByISO
	origWeatherFor := weatherFor
	origAirQualityFor := airQualityFor
	origExchangeRatesFor := exchangeRatesFor
	defer func() {
		countryByISO = origCountryByISO
		weatherFor = origWeatherFor
		airQualityFor = origAirQualityFor
		exchangeRatesFor = origExchangeRatesFor
	}()

	resetRegistrations()

	memStore().registrations["dash1"] = models.Registration{
		ID:      "dash1",
		Country: "Norway",
		IsoCode: "NO",
		Features: models.Features{
			Temperature:      true,
			Precipitation:    true,
			AirQuality:       true,
			Capital:          true,
			Coordinates:      true,
			Population:       true,
			Area:             true,
			TargetCurrencies: []string{"EUR", "USD"},
		},
	}

	countryByISO = func(code string) (*models.Country, error) {
		return &models.Country{
			Name:       "Norway",
			Code:       "NO",
			Capital:    "Oslo",
			Population: 5379475,
			Area:       323802.0,
			Currencies: []string{"NOK"},
			Coordinates: models.Coordinates{
				Latitude:  62.0,
				Longitude: 10.0,
			},
		}, nil
	}

	weatherFor = func(latitude, longitude float64) (*clients.MeteoResult, error) {
		return &clients.MeteoResult{
			Temperature:   4.0,
			Precipitation: 0.2,
		}, nil
	}

	airQualityFor = func(latitude, longitude float64) (*models.AirQualityData, error) {
		return &models.AirQualityData{
			PM25:  8.4,
			PM10:  14.2,
			Level: "Good",
		}, nil
	}

	exchangeRatesFor = func(baseCurrency string, targets []string) (map[string]float64, error) {
		return map[string]float64{
			"EUR": 0.087701,
			"USD": 0.095184,
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/dashboards/dash1", nil)
	req.SetPathValue("id", "dash1")
	w := httptest.NewRecorder()

	GetDashboardHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp models.DashboardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Country != "Norway" {
		t.Errorf("country: got %q, want Norway", resp.Country)
	}
	if resp.IsoCode != "NO" {
		t.Errorf("isoCode: got %q, want NO", resp.IsoCode)
	}
	if resp.Features.Temperature == nil || *resp.Features.Temperature != 4.0 {
		t.Errorf("temperature: got %+v", resp.Features.Temperature)
	}
	if resp.Features.Precipitation == nil || *resp.Features.Precipitation != 0.2 {
		t.Errorf("precipitation: got %+v", resp.Features.Precipitation)
	}
	if resp.Features.AirQuality == nil || resp.Features.AirQuality.Level != "Good" {
		t.Errorf("air quality: got %+v", resp.Features.AirQuality)
	}
	if resp.Features.Capital == nil || *resp.Features.Capital != "Oslo" {
		t.Errorf("capital: got %+v", resp.Features.Capital)
	}
	if resp.Features.Coordinates == nil || resp.Features.Coordinates.Latitude != 62.0 {
		t.Errorf("coordinates: got %+v", resp.Features.Coordinates)
	}
	if resp.Features.Population == nil || *resp.Features.Population != 5379475 {
		t.Errorf("population: got %+v", resp.Features.Population)
	}
	if resp.Features.Area == nil || *resp.Features.Area != 323802.0 {
		t.Errorf("area: got %+v", resp.Features.Area)
	}
	if len(resp.Features.TargetCurrencies) != 2 {
		t.Errorf("targetCurrencies: got %+v", resp.Features.TargetCurrencies)
	}
	if resp.LastRetrieval == "" {
		t.Error("expected lastRetrieval to be set")
	}
}

func TestGetDashboardHandler_WeatherFailure(t *testing.T) {
	origCountryByISO := countryByISO
	origWeatherFor := weatherFor
	defer func() {
		countryByISO = origCountryByISO
		weatherFor = origWeatherFor
	}()

	resetRegistrations()

	memStore().registrations["dash2"] = models.Registration{
		ID:      "dash2",
		Country: "Norway",
		IsoCode: "NO",
		Features: models.Features{
			Temperature: true,
		},
	}

	countryByISO = func(code string) (*models.Country, error) {
		return &models.Country{
			Name: "Norway",
			Code: "NO",
			Coordinates: models.Coordinates{
				Latitude:  62.0,
				Longitude: 10.0,
			},
		}, nil
	}

	weatherFor = func(latitude, longitude float64) (*clients.MeteoResult, error) {
		return nil, errors.New("weather down")
	}

	req := httptest.NewRequest(http.MethodGet, "/dashboards/dash2", nil)
	req.SetPathValue("id", "dash2")
	w := httptest.NewRecorder()

	GetDashboardHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
