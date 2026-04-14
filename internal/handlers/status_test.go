package handlers

import (
	"assignment_2/internal/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type failingNotificationStore struct {
	*MemoryStore
}

func (f *failingNotificationStore) ListNotifications(_ context.Context) ([]models.NotificationRegistration, error) {
	return nil, errors.New("db unavailable")
}

func TestStatusHandler_OK(t *testing.T) {
	origProbeFn := probeFn
	origStore := store
	origStartTime := startTime
	defer func() {
		probeFn = origProbeFn
		store = origStore
		startTime = origStartTime
	}()

	store = NewMemoryStore()
	memStore().notifications["1"] = models.NotificationRegistration{ID: "1"}
	memStore().notifications["2"] = models.NotificationRegistration{ID: "2"}
	startTime = time.Now().Add(-12 * time.Second)

	probeFn = func(url, userAgent string) int {
		switch url {
		case "http://129.241.150.113:8080/v3.1/name/norge":
			return http.StatusOK
		case "https://api.open-meteo.com/v1/forecast":
			return http.StatusOK
		case "https://api.openaq.org/v3/locations":
			return http.StatusOK
		case "https://nominatim.openstreetmap.org/":
			if userAgent != "prog2005-assignment2/1.0" {
				t.Errorf("expected nominatim user agent, got %q", userAgent)
			}
			return http.StatusOK
		case "http://129.241.150.113:9090/currency/nok":
			return http.StatusOK
		default:
			t.Errorf("unexpected probe url %q", url)
			return http.StatusServiceUnavailable
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/envdash/v1/status/", nil)
	w := httptest.NewRecorder()

	StatusHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}

	var resp models.StatusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.CountriesAPI != http.StatusOK {
		t.Errorf("countries_api: got %d, want %d", resp.CountriesAPI, http.StatusOK)
	}
	if resp.MeteoAPI != http.StatusOK {
		t.Errorf("meteo_api: got %d, want %d", resp.MeteoAPI, http.StatusOK)
	}
	if resp.OpenAQAPI != http.StatusOK {
		t.Errorf("openaq_api: got %d, want %d", resp.OpenAQAPI, http.StatusOK)
	}
	if resp.NominatimAPI != http.StatusOK {
		t.Errorf("nominatim_api: got %d, want %d", resp.NominatimAPI, http.StatusOK)
	}
	if resp.CurrencyAPI != http.StatusOK {
		t.Errorf("currency_api: got %d, want %d", resp.CurrencyAPI, http.StatusOK)
	}
	if resp.NotificationDB != http.StatusOK {
		t.Errorf("notification_db: got %d, want %d", resp.NotificationDB, http.StatusOK)
	}
	if resp.Webhooks != 2 {
		t.Errorf("webhooks: got %d, want %d", resp.Webhooks, 2)
	}
	if resp.Version != "v1" {
		t.Errorf("version: got %q, want %q", resp.Version, "v1")
	}
	if resp.Uptime < 10 {
		t.Errorf("uptime: got %d, want >= 10", resp.Uptime)
	}
}

func TestStatusHandler_UpstreamFailureReturns500(t *testing.T) {
	origProbeFn := probeFn
	origStore := store
	origStartTime := startTime
	defer func() {
		probeFn = origProbeFn
		store = origStore
		startTime = origStartTime
	}()

	store = NewMemoryStore()
	startTime = time.Now().Add(-3 * time.Second)

	probeFn = func(url, userAgent string) int {
		if url == "https://api.open-meteo.com/v1/forecast" {
			return http.StatusServiceUnavailable
		}
		return http.StatusOK
	}

	req := httptest.NewRequest(http.MethodGet, "/envdash/v1/status/", nil)
	w := httptest.NewRecorder()

	StatusHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}

	var resp models.StatusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.MeteoAPI != http.StatusServiceUnavailable {
		t.Errorf("meteo_api: got %d, want %d", resp.MeteoAPI, http.StatusServiceUnavailable)
	}
}

func TestStatusHandler_DBFailure(t *testing.T) {
	origProbeFn := probeFn
	origStore := store
	origStartTime := startTime
	defer func() {
		probeFn = origProbeFn
		store = origStore
		startTime = origStartTime
	}()

	store = &failingNotificationStore{MemoryStore: NewMemoryStore()}
	startTime = time.Now().Add(-3 * time.Second)
	probeFn = func(url, userAgent string) int { return http.StatusOK }

	req := httptest.NewRequest(http.MethodGet, "/envdash/v1/status/", nil)
	w := httptest.NewRecorder()

	StatusHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}

	var resp models.StatusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.NotificationDB != http.StatusServiceUnavailable {
		t.Errorf("notification_db: got %d, want %d", resp.NotificationDB, http.StatusServiceUnavailable)
	}
	if resp.Webhooks != 0 {
		t.Errorf("webhooks: got %d, want %d", resp.Webhooks, 0)
	}
}
