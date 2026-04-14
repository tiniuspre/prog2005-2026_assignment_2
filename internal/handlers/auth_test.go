package handlers

import (
	"assignment_2/internal/firebase"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestCreateAPIKeyHandler(t *testing.T) {
	origGenerateKeyFn := generateKeyFn
	origCreateAPIKeyFn := createAPIKeyFn
	origFormatTimestampFn := formatTimestampFn
	origFirestoreClient := firestoreClient
	defer func() {
		generateKeyFn = origGenerateKeyFn
		createAPIKeyFn = origCreateAPIKeyFn
		formatTimestampFn = origFormatTimestampFn
		firestoreClient = origFirestoreClient
	}()

	firestoreClient = nil

	t.Run("success", func(t *testing.T) {
		generateKeyFn = func() (string, error) {
			return "sk-envdash-testkey", nil
		}
		formatTimestampFn = func() string {
			return "20250301 09:15"
		}

		var got firebase.APIKey
		createAPIKeyFn = func(_ context.Context, _ *firestore.Client, ak firebase.APIKey) error {
			got = ak
			return nil
		}

		body := bytes.NewBufferString(`{"name":"my-client-app","email":"user@example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/auth/", body)
		w := httptest.NewRecorder()

		CreateAPIKeyHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusCreated)
		}

		if got.Key != "sk-envdash-testkey" {
			t.Errorf("stored key: got %q", got.Key)
		}
		if got.Name != "my-client-app" {
			t.Errorf("stored name: got %q", got.Name)
		}
		if got.Email != "user@example.com" {
			t.Errorf("stored email: got %q", got.Email)
		}
		if got.CreatedAt != "20250301 09:15" {
			t.Errorf("stored createdAt: got %q", got.CreatedAt)
		}

		var resp authResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if resp.Key != "sk-envdash-testkey" {
			t.Errorf("response key: got %q", resp.Key)
		}
		if resp.CreatedAt != "20250301 09:15" {
			t.Errorf("response createdAt: got %q", resp.CreatedAt)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/", bytes.NewBufferString(`{`))
		w := httptest.NewRecorder()

		CreateAPIKeyHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/", bytes.NewBufferString(`{"name":"x"}`))
		w := httptest.NewRecorder()

		CreateAPIKeyHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("key generation failure", func(t *testing.T) {
		generateKeyFn = func() (string, error) {
			return "", errors.New("boom")
		}

		req := httptest.NewRequest(http.MethodPost, "/auth/", bytes.NewBufferString(`{"name":"x","email":"x@y.z"}`))
		w := httptest.NewRecorder()

		CreateAPIKeyHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("store failure", func(t *testing.T) {
		generateKeyFn = func() (string, error) {
			return "sk-envdash-testkey", nil
		}
		formatTimestampFn = func() string {
			return "20250301 09:15"
		}
		createAPIKeyFn = func(_ context.Context, _ *firestore.Client, _ firebase.APIKey) error {
			return errors.New("boom")
		}

		req := httptest.NewRequest(http.MethodPost, "/auth/", bytes.NewBufferString(`{"name":"x","email":"x@y.z"}`))
		w := httptest.NewRecorder()

		CreateAPIKeyHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

func TestDeleteAPIKeyHandler(t *testing.T) {
	origDeleteAPIKeyFn := deleteAPIKeyFn
	origFirestoreClient := firestoreClient
	defer func() {
		deleteAPIKeyFn = origDeleteAPIKeyFn
		firestoreClient = origFirestoreClient
	}()

	firestoreClient = nil

	t.Run("success", func(t *testing.T) {
		deleteAPIKeyFn = func(_ context.Context, _ *firestore.Client, key string) (bool, error) {
			if key != "sk-envdash-test" {
				t.Errorf("delete called with key %q", key)
			}
			return true, nil
		}

		req := httptest.NewRequest(http.MethodDelete, "/auth/sk-envdash-test", nil)
		req.SetPathValue("key", "sk-envdash-test")
		w := httptest.NewRecorder()

		DeleteAPIKeyHandler(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusNoContent)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/auth/", nil)
		w := httptest.NewRecorder()

		DeleteAPIKeyHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found", func(t *testing.T) {
		deleteAPIKeyFn = func(_ context.Context, _ *firestore.Client, _ string) (bool, error) {
			return false, nil
		}

		req := httptest.NewRequest(http.MethodDelete, "/auth/sk-missing", nil)
		req.SetPathValue("key", "sk-missing")
		w := httptest.NewRecorder()

		DeleteAPIKeyHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("delete failure", func(t *testing.T) {
		deleteAPIKeyFn = func(_ context.Context, _ *firestore.Client, _ string) (bool, error) {
			return false, errors.New("boom")
		}

		req := httptest.NewRequest(http.MethodDelete, "/auth/sk-error", nil)
		req.SetPathValue("key", "sk-error")
		w := httptest.NewRecorder()

		DeleteAPIKeyHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}
