package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

var store Store

// Init initializes the handlers with a Firestore client.
func Init(client *firestore.Client) {
	store = NewFirestoreStore(client)
}

// Helper functions for writing JSON responses and errors
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})

}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
