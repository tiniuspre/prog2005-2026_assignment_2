package handlers

import (
	"assignment_2/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// resetRegistrations clears the store to ensure test isolation.
func resetRegistrations() {
	store = NewMemoryStore()
}

func TestCreateRegistrationHandler_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name:       "invalid json body",
			body:       "not json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing both country and isoCode",
			body:       models.Registration{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRegistrations()

			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(tt.body)

			req := httptest.NewRequest(http.MethodPost, "/registrations", &buf)
			w := httptest.NewRecorder()

			CreateRegistrationHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status: got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetRegistrationHandler(t *testing.T) {
	t.Run("existing id", func(t *testing.T) {
		resetRegistrations()

		memStore().registrations["reg1"] = models.Registration{
			ID:      "reg1",
			Country: "Norway",
			IsoCode: "NO",
		}

		req := httptest.NewRequest(http.MethodGet, "/registrations/reg1", nil)
		req.SetPathValue("id", "reg1")
		w := httptest.NewRecorder()

		GetRegistrationHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
		}

		var result models.Registration
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if result.ID != "reg1" {
			t.Errorf("ID: got %s, want reg1", result.ID)
		}
		if result.IsoCode != "NO" {
			t.Errorf("IsoCode: got %s, want NO", result.IsoCode)
		}
	})

	t.Run("unknown id", func(t *testing.T) {
		resetRegistrations()

		req := httptest.NewRequest(http.MethodGet, "/registrations/nope", nil)
		req.SetPathValue("id", "nope")
		w := httptest.NewRecorder()

		GetRegistrationHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestListRegistrationsHandler(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		resetRegistrations()

		req := httptest.NewRequest(http.MethodGet, "/registrations", nil)
		w := httptest.NewRecorder()

		ListRegistrationsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
		}

		var result []models.Registration
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty list, got %d entries", len(result))
		}
	})

	t.Run("returns all entries", func(t *testing.T) {
		resetRegistrations()

		memStore().registrations["r1"] = models.Registration{ID: "r1", Country: "Norway", IsoCode: "NO"}
		memStore().registrations["r2"] = models.Registration{ID: "r2", Country: "Sweden", IsoCode: "SE"}

		req := httptest.NewRequest(http.MethodGet, "/registrations", nil)
		w := httptest.NewRecorder()

		ListRegistrationsHandler(w, req)

		var result []models.Registration
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 entries, got %d", len(result))
		}
	})
}

func TestDeleteRegistrationHandler(t *testing.T) {
	t.Run("unknown id", func(t *testing.T) {
		resetRegistrations()

		req := httptest.NewRequest(http.MethodDelete, "/registrations/nope", nil)
		req.SetPathValue("id", "nope")
		w := httptest.NewRecorder()

		DeleteRegistrationHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("existing id", func(t *testing.T) {
		resetRegistrations()

		memStore().registrations["del1"] = models.Registration{ID: "del1", Country: "Norway", IsoCode: "NO"}

		req := httptest.NewRequest(http.MethodDelete, "/registrations/del1", nil)
		req.SetPathValue("id", "del1")
		w := httptest.NewRecorder()

		DeleteRegistrationHandler(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNoContent)
		}
		if _, ok := memStore().registrations["del1"]; ok {
			t.Error("expected entry to be removed from store")
		}
	})
}

func TestUpdateRegistrationHandler_Validation(t *testing.T) {
	t.Run("unknown id", func(t *testing.T) {
		resetRegistrations()

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(models.Registration{Country: "Norway"})

		req := httptest.NewRequest(http.MethodPut, "/registrations/nope", &buf)
		req.SetPathValue("id", "nope")
		w := httptest.NewRecorder()

		UpdateRegistrationHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("invalid json body", func(t *testing.T) {
		resetRegistrations()

		memStore().registrations["upd1"] = models.Registration{ID: "upd1", Country: "Norway", IsoCode: "NO"}

		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode("not json")

		req := httptest.NewRequest(http.MethodPut, "/registrations/upd1", &buf)
		req.SetPathValue("id", "upd1")
		w := httptest.NewRecorder()

		UpdateRegistrationHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}
