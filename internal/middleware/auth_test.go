package middleware

import (
	"assignment_2/internal/firebase"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestAPIKeyAuth(t *testing.T) {
	origGetAPIKeyFn := getAPIKeyFn
	defer func() {
		getAPIKeyFn = origGetAPIKeyFn
	}()

	t.Run("public route bypasses auth", func(t *testing.T) {
		called := false

		getAPIKeyFn = func(_ context.Context, _ *firestore.Client, _ string) (*firebase.APIKey, error) {
			t.Fatal("getAPIKeyFn should not be called for public route")
			return nil, nil
		}

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusNoContent)
		})

		handler := APIKeyAuth(nil, []PublicRoute{
			{Method: http.MethodGet, Path: "/envdash/v1/status/"},
		})(next)

		req := httptest.NewRequest(http.MethodGet, "/envdash/v1/status/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if !called {
			t.Fatal("expected next handler to be called")
		}
		if w.Code != http.StatusNoContent {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusNoContent)
		}
	})

	t.Run("missing key returns 401", func(t *testing.T) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler should not be called")
		})

		handler := APIKeyAuth(nil, nil)(next)

		req := httptest.NewRequest(http.MethodGet, "/envdash/v1/dashboards/1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("validator error returns 500", func(t *testing.T) {
		getAPIKeyFn = func(_ context.Context, _ *firestore.Client, _ string) (*firebase.APIKey, error) {
			return nil, errors.New("boom")
		}

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler should not be called")
		})

		handler := APIKeyAuth(nil, nil)(next)

		req := httptest.NewRequest(http.MethodGet, "/envdash/v1/dashboards/1", nil)
		req.Header.Set("X-API-Key", "sk-test")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("invalid key returns 403", func(t *testing.T) {
		getAPIKeyFn = func(_ context.Context, _ *firestore.Client, _ string) (*firebase.APIKey, error) {
			return nil, nil
		}

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler should not be called")
		})

		handler := APIKeyAuth(nil, nil)(next)

		req := httptest.NewRequest(http.MethodGet, "/envdash/v1/dashboards/1", nil)
		req.Header.Set("X-API-Key", "sk-invalid")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusForbidden)
		}
	})

	t.Run("valid key calls next", func(t *testing.T) {
		called := false

		getAPIKeyFn = func(_ context.Context, _ *firestore.Client, key string) (*firebase.APIKey, error) {
			return &firebase.APIKey{Key: key}, nil
		}

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusAccepted)
		})

		handler := APIKeyAuth(nil, nil)(next)

		req := httptest.NewRequest(http.MethodGet, "/envdash/v1/dashboards/1", nil)
		req.Header.Set("X-API-Key", "sk-valid")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if !called {
			t.Fatal("expected next handler to be called")
		}
		if w.Code != http.StatusAccepted {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusAccepted)
		}
	})
}
