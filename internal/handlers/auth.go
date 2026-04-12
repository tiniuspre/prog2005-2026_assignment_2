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
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Email == "" {
		http.Error(w, "name and email are required", http.StatusBadRequest)
		return
	}

	key, err := generateKey()
	if err != nil {
		http.Error(w, "failed to generate key", http.StatusInternalServerError)
		return
	}

	ak := firebase.APIKey{
		Key:       key,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: firebase.FormatTimestamp(),
	}

	if err := firebase.CreateAPIKey(r.Context(), firestoreClient, ak); err != nil {
		http.Error(w, "failed to store API key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse{
		Key:       ak.Key,
		CreatedAt: ak.CreatedAt,
	})
}

// DeleteAPIKeyHandler handles DELETE /auth/{key} and revokes the API key.
func DeleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	found, err := firebase.DeleteAPIKey(r.Context(), firestoreClient, key)
	if err != nil {
		http.Error(w, "failed to delete API key", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "API key not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
