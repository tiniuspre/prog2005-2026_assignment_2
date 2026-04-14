package handlers

import (
	"assignment_2/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRegistrationHandler_Success(t *testing.T) {
	origCountryByISO := countryByISO
	origCountryByName := countryByName
	defer func() {
		countryByISO = origCountryByISO
		countryByName = origCountryByName
	}()

	t.Run("create by iso code", func(t *testing.T) {
		resetRegistrations()

		countryByISO = func(code string) (*models.Country, error) {
			if code != "NO" {
				t.Fatalf("unexpected iso code %q", code)
			}
			return &models.Country{Name: "Norway", Code: "NO"}, nil
		}

		body := models.Registration{
			IsoCode: "NO",
			Features: models.Features{
				Temperature: true,
			},
		}

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(body)

		req := httptest.NewRequest(http.MethodPost, "/registrations", &buf)
		w := httptest.NewRecorder()

		CreateRegistrationHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusCreated)
		}

		var resp map[string]string
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if resp["id"] == "" {
			t.Fatal("expected id in response")
		}
		if resp["lastChange"] == "" {
			t.Fatal("expected lastChange in response")
		}

		if len(memStore().registrations) != 1 {
			t.Fatalf("expected 1 stored registration, got %d", len(memStore().registrations))
		}

		for _, reg := range memStore().registrations {
			if reg.Country != "Norway" {
				t.Errorf("stored country: got %q, want Norway", reg.Country)
			}
			if reg.IsoCode != "NO" {
				t.Errorf("stored isoCode: got %q, want NO", reg.IsoCode)
			}
			if !reg.Features.Temperature {
				t.Error("stored feature temperature should be true")
			}
			if reg.LastChange == "" {
				t.Error("stored lastChange should be set")
			}
		}
	})

	t.Run("create by country name", func(t *testing.T) {
		resetRegistrations()

		countryByName = func(name string) (*models.Country, error) {
			if name != "Norway" {
				t.Fatalf("unexpected country %q", name)
			}
			return &models.Country{Name: "Norway", Code: "NO"}, nil
		}

		body := models.Registration{
			Country: "Norway",
		}

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(body)

		req := httptest.NewRequest(http.MethodPost, "/registrations", &buf)
		w := httptest.NewRecorder()

		CreateRegistrationHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusCreated)
		}

		if len(memStore().registrations) != 1 {
			t.Fatalf("expected 1 stored registration, got %d", len(memStore().registrations))
		}

		for _, reg := range memStore().registrations {
			if reg.IsoCode != "NO" {
				t.Errorf("stored isoCode: got %q, want NO", reg.IsoCode)
			}
		}
	})

	t.Run("invalid iso code", func(t *testing.T) {
		resetRegistrations()

		countryByISO = func(code string) (*models.Country, error) {
			return nil, errors.New("not found")
		}

		body := models.Registration{IsoCode: "XX"}
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(body)

		req := httptest.NewRequest(http.MethodPost, "/registrations", &buf)
		w := httptest.NewRecorder()

		CreateRegistrationHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestUpdateRegistrationHandler_Success(t *testing.T) {
	origCountryByISO := countryByISO
	defer func() {
		countryByISO = origCountryByISO
	}()

	resetRegistrations()

	memStore().registrations["reg1"] = models.Registration{
		ID:      "reg1",
		Country: "Norway",
		IsoCode: "NO",
		Features: models.Features{
			Temperature: true,
		},
	}

	countryByISO = func(code string) (*models.Country, error) {
		if code != "SE" {
			t.Fatalf("unexpected iso code %q", code)
		}
		return &models.Country{Name: "Sweden", Code: "SE"}, nil
	}

	update := models.Registration{
		IsoCode: "SE",
		Features: models.Features{
			Area: true,
		},
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(update)

	req := httptest.NewRequest(http.MethodPut, "/registrations/reg1", &buf)
	req.SetPathValue("id", "reg1")
	w := httptest.NewRecorder()

	UpdateRegistrationHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}

	reg, ok := memStore().registrations["reg1"]
	if !ok {
		t.Fatal("updated registration not found in store")
	}
	if reg.Country != "Sweden" {
		t.Errorf("country: got %q, want Sweden", reg.Country)
	}
	if reg.IsoCode != "SE" {
		t.Errorf("isoCode: got %q, want SE", reg.IsoCode)
	}
	if !reg.Features.Area {
		t.Error("expected area feature to be true")
	}
	if reg.LastChange == "" {
		t.Error("expected lastChange to be updated")
	}
}
