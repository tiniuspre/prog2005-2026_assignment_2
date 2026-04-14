package handlers

import (
	"assignment_2/internal/firebase"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

// authRequest is the expected body at POST /auth/.
type authRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// authResponse is returned on successful registration.
type authResponse struct {
	Key       string `json:"key"`
	CreatedAt string `json:"createdAt"`
}

// generateKey produces a key like "sk-envdash-a3...".
func generateKey() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "sk-envdash-" + hex.EncodeToString(b), nil
}

// CreateAPIKeyHandler handles POST /auth/ and registers a new client.
func CreateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "name and email are required")
		return
	}

	key, err := generateKeyFn()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate key")
		return
	}

	ak := firebase.APIKey{
		Key:       key,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: formatTimestampFn(),
	}

	if err := createAPIKeyFn(r.Context(), firestoreClient, ak); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to store API key")
		return
	}

	writeJSON(w, http.StatusCreated, authResponse{
		Key:       ak.Key,
		CreatedAt: ak.CreatedAt,
	})
}

// DeleteAPIKeyHandler handles DELETE /auth/{key} and revokes the API key.
func DeleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeError(w, http.StatusBadRequest, "key is required")
		return
	}

	found, err := deleteAPIKeyFn(r.Context(), firestoreClient, key)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete API key")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "API key not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
