package handlers

import (
	"assignment_2/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetNotifications() {
	notifications = map[string]models.NotificationRegistration{}
}

func TestCreateNotificationHandler(t *testing.T) {
	threshold := &models.ThresholdConfig{Field: "temp", Operator: ">", Value: 30}

	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name:       "valid REGISTER event",
			body:       models.NotificationRegistration{URL: "http://example.com", Country: "NO", Event: "REGISTER"},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "valid THRESHOLD event",
			body:       models.NotificationRegistration{URL: "https://example.com/hook", Country: "NO", Event: "THRESHOLD", Threshold: threshold},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing url",
			body:       models.NotificationRegistration{Country: "NO", Event: "REGISTER"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing country",
			body:       models.NotificationRegistration{URL: "http://example.com", Event: "REGISTER"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing event",
			body:       models.NotificationRegistration{URL: "http://example.com", Country: "NO"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid url scheme",
			body:       models.NotificationRegistration{URL: "ftp://example.com", Country: "NO", Event: "REGISTER"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "url with no host",
			body:       models.NotificationRegistration{URL: "http://", Country: "NO", Event: "REGISTER"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "THRESHOLD event missing threshold field",
			body:       models.NotificationRegistration{URL: "http://example.com", Country: "NO", Event: "THRESHOLD"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json body",
			body:       "not json",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetNotifications()

			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(tt.body)

			req := httptest.NewRequest(http.MethodPost, "/notifications", &buf)
			w := httptest.NewRecorder()

			CreateNotificationHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status: got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestCreateNotificationHandler_ResponseBody(t *testing.T) {
	resetNotifications()

	reg := models.NotificationRegistration{
		URL:     "http://example.com",
		Country: "NO",
		Event:   "REGISTER"}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(reg)

	req := httptest.NewRequest(http.MethodPost, "/notifications", &buf)
	w := httptest.NewRecorder()

	CreateNotificationHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusCreated)
	}

	var result models.NotificationRegistration
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.ID == "" {
		t.Error("expected non-empty ID in response")
	}
	if result.URL != reg.URL || result.Country != reg.Country || result.Event != reg.Event {
		t.Errorf("response body mismatch: got %+v", result)
	}
}

func TestGetNotificationHandler(t *testing.T) {
	resetNotifications()

	notifications["abc123"] = models.NotificationRegistration{
		ID:      "abc123",
		URL:     "http://example.com",
		Country: "NO",
		Event:   "REGISTER"}

	t.Run("existing id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notifications/abc123", nil)
		req.SetPathValue("id", "abc123")
		w := httptest.NewRecorder()

		GetNotificationHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
		}

		var result models.NotificationRegistration
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if result.ID != "abc123" {
			t.Errorf("ID: got %s, want abc123", result.ID)
		}
	})

	t.Run("unknown id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notifications/nope", nil)
		req.SetPathValue("id", "nope")
		w := httptest.NewRecorder()

		GetNotificationHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestListNotificationsHandler(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		resetNotifications()

		req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
		w := httptest.NewRecorder()

		ListNotificationsHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
		}

		var result []models.NotificationRegistration
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty list, got %d entries", len(result))
		}
	})

	t.Run("returns all entries", func(t *testing.T) {
		resetNotifications()

		notifications["1"] = models.NotificationRegistration{ID: "1",
			URL:     "http://a.com",
			Country: "NO",
			Event:   "REGISTER"}

		notifications["2"] = models.NotificationRegistration{ID: "2",
			URL:     "http://b.com",
			Country: "SE",
			Event:   "DELETE"}

		req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
		w := httptest.NewRecorder()

		ListNotificationsHandler(w, req)

		var result []models.NotificationRegistration
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 entries, got %d", len(result))
		}
	})
}

func TestDeleteNotificationHandler(t *testing.T) {
	t.Run("unknown id", func(t *testing.T) {
		resetNotifications()

		req := httptest.NewRequest(http.MethodDelete, "/notifications/nope", nil)
		req.SetPathValue("id", "nope")
		w := httptest.NewRecorder()

		DeleteNotificationHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("existing id", func(t *testing.T) {
		resetNotifications()

		// Use a test server so the goroutine webhook dispatch doesn't fail noisily
		webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer webhook.Close()

		notifications["del1"] = models.NotificationRegistration{
			ID:      "del1",
			URL:     webhook.URL,
			Country: "NO",
			Event:   "REGISTER",
		}

		req := httptest.NewRequest(http.MethodDelete, "/notifications/del1", nil)
		req.SetPathValue("id", "del1")
		w := httptest.NewRecorder()

		DeleteNotificationHandler(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNoContent)
		}
		if _, ok := notifications["del1"]; ok {
			t.Error("expected entry to be removed from map")
		}
	})
}
