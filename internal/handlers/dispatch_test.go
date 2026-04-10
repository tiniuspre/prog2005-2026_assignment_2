package handlers

import (
	"assignment_2/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// --- evaluate ---

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name      string
		measured  float64
		operator  string
		threshold float64
		want      bool
	}{
		{"greater than or equal - above", 35, ">=", 30, true},
		{"greater than or equal - equal", 30, ">=", 30, true},
		{"greater than or equal - below", 25, ">=", 30, false},
		{"less than or equal - below", 10, "<=", 20, true},
		{"less than or equal - equal", 20, "<=", 20, true},
		{"less than or equal - above", 25, "<=", 20, false},
		{"unknown operator", 10, ">", 5, false},
		{"unknown operator strict less", 3, "<", 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evaluate(tt.measured, tt.operator, tt.threshold)
			if got != tt.want {
				t.Errorf("evaluate(%v, %q, %v) = %v, want %v",
					tt.measured, tt.operator, tt.threshold, got, tt.want)
			}
		})
	}
}

// --- DispatchEvent ---

func TestDispatchEvent(t *testing.T) {
	t.Run("fires for matching event and country", func(t *testing.T) {
		resetNotifications()

		done := make(chan models.WebhookPayload, 1)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var p models.WebhookPayload
			_ = json.NewDecoder(r.Body).Decode(&p)
			done <- p
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		notifications["n1"] = models.NotificationRegistration{
			ID: "n1", URL: srv.URL, Country: "NO", Event: "REGISTER",
		}

		DispatchEvent("REGISTER", "NO")

		select {
		case p := <-done:
			if p.Event != "REGISTER" || p.Country != "NO" {
				t.Errorf("unexpected payload: %+v", p)
			}
		case <-time.After(2 * time.Second):
			t.Error("timed out waiting for webhook")
		}
	})

	t.Run("does not fire for wrong event", func(t *testing.T) {
		resetNotifications()

		done := make(chan struct{}, 1)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			done <- struct{}{}
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		notifications["n2"] = models.NotificationRegistration{
			ID: "n2", URL: srv.URL, Country: "NO", Event: "DELETE",
		}

		DispatchEvent("REGISTER", "NO")

		select {
		case <-done:
			t.Error("webhook should not have fired for non-matching event")
		case <-time.After(200 * time.Millisecond):
			// expected: no webhook fired
		}
	})

	t.Run("does not fire for wrong country", func(t *testing.T) {
		resetNotifications()

		done := make(chan struct{}, 1)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			done <- struct{}{}
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		notifications["n3"] = models.NotificationRegistration{
			ID: "n3", URL: srv.URL, Country: "SE", Event: "REGISTER",
		}

		DispatchEvent("REGISTER", "NO")

		select {
		case <-done:
			t.Error("webhook should not have fired for non-matching country")
		case <-time.After(200 * time.Millisecond):
			// expected: no webhook fired
		}
	})
}

// --- CheckThresholds ---

func TestCheckThresholds(t *testing.T) {
	t.Run("fires when threshold is met", func(t *testing.T) {
		resetNotifications()

		done := make(chan models.WebhookPayload, 1)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var p models.WebhookPayload
			_ = json.NewDecoder(r.Body).Decode(&p)
			done <- p
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		notifications["t1"] = models.NotificationRegistration{
			ID:      "t1",
			URL:     srv.URL,
			Country: "NO",
			Event:   "THRESHOLD",
			Threshold: &models.ThresholdConfig{
				Field: "temperature", Operator: ">=", Value: 30,
			},
		}

		CheckThresholds("NO", map[string]float64{"temperature": 35})

		select {
		case p := <-done:
			if p.Event != "THRESHOLD" {
				t.Errorf("event: got %s, want THRESHOLD", p.Event)
			}
			if p.Details == nil {
				t.Fatal("expected Details in payload, got nil")
			}
			if p.Details.Field != "temperature" {
				t.Errorf("details.field: got %s, want temperature", p.Details.Field)
			}
			if p.Details.MeasuredValue != 35 {
				t.Errorf("details.measuredValue: got %v, want 35", p.Details.MeasuredValue)
			}
			if p.Details.Threshold != 30 {
				t.Errorf("details.threshold: got %v, want 30", p.Details.Threshold)
			}
		case <-time.After(2 * time.Second):
			t.Error("timed out waiting for webhook")
		}
	})

	t.Run("does not fire when threshold is not met", func(t *testing.T) {
		resetNotifications()

		done := make(chan struct{}, 1)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			done <- struct{}{}
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		notifications["t2"] = models.NotificationRegistration{
			ID:      "t2",
			URL:     srv.URL,
			Country: "NO",
			Event:   "THRESHOLD",
			Threshold: &models.ThresholdConfig{
				Field: "temperature", Operator: ">=", Value: 30,
			},
		}

		CheckThresholds("NO", map[string]float64{"temperature": 25})

		select {
		case <-done:
			t.Error("webhook should not fire when threshold is not met")
		case <-time.After(200 * time.Millisecond):
			// expected: no webhook fired
		}
	})

	t.Run("does not fire when field is absent from values", func(t *testing.T) {
		resetNotifications()

		done := make(chan struct{}, 1)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			done <- struct{}{}
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		notifications["t3"] = models.NotificationRegistration{
			ID:      "t3",
			URL:     srv.URL,
			Country: "NO",
			Event:   "THRESHOLD",
			Threshold: &models.ThresholdConfig{
				Field: "pm25", Operator: ">=", Value: 50,
			},
		}

		CheckThresholds("NO", map[string]float64{"temperature": 35}) // pm25 absent

		select {
		case <-done:
			t.Error("webhook should not fire when field is absent")
		case <-time.After(200 * time.Millisecond):
			// expected: no webhook fired
		}
	})
}
