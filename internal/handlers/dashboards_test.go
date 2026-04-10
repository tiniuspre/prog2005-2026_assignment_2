package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// only has one test since the other handlers require live api calls.
func TestGetDashboardHandler(t *testing.T) {
	t.Run("unknown registration id", func(t *testing.T) {
		resetRegistrations()

		req := httptest.NewRequest(http.MethodGet, "/dashboards/nope", nil)
		req.SetPathValue("id", "nope")
		w := httptest.NewRecorder()

		GetDashboardHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
