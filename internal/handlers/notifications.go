package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// temporary map (pre firebase)
var notifications = map[string]NotificationRegistration{}

func createNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var reg NotificationRegistration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if reg.URL == "" || reg.Country == "" || reg.Event == "" {
		writeError(w, http.StatusBadRequest, "url, country and event are required")
		return
	}

	if reg.Event == "THRESHOLD" && reg.Threshold == nil {
		writeError(w, http.StatusBadRequest, "threshold field required for THRESHOLD events")
		return
	}

	// temp storage (pre firebase)
	reg.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	// storing
	notifications[reg.ID] = reg

	writeJSON(w, http.StatusCreated, reg)
}

func getNotificationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reg, ok := notifications[id]
	if !ok {
		writeError(w, http.StatusNotFound, "notification not found")
		return
	}

	writeJSON(w, http.StatusOK, reg)
}

func listNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	result := make([]NotificationRegistration, 0, len(notifications))
	for _, reg := range notifications {
		result = append(result, reg)
	}
	writeJSON(w, http.StatusOK, result)
}

func deleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, ok := notifications[id]; !ok {
		writeError(w, http.StatusNotFound, "notification not found")
		return
	}
	delete(notifications, id)

	writeJSON(w, http.StatusNoContent, nil)
}

// --------- HELPER FUNCTIONS ---------
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
