package handlers

import (
	"assignment_2/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// temporary map (pre firebase)
var notifications = map[string]models.NotificationRegistration{}

func CreateNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var reg models.NotificationRegistration
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

func GetNotificationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reg, ok := notifications[id]
	if !ok {
		writeError(w, http.StatusNotFound, "notification not found")
		return
	}

	writeJSON(w, http.StatusOK, reg)
}

func ListNotificationsHandler(w http.ResponseWriter, _ *http.Request) {
	result := make([]models.NotificationRegistration, 0, len(notifications))
	for _, reg := range notifications {
		result = append(result, reg)
	}
	writeJSON(w, http.StatusOK, result)
}

func DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reg, ok := notifications[id]
	if !ok {
		writeError(w, http.StatusNotFound, "notification not found")
		return
	}
	delete(notifications, id)

	go Send(reg.URL, models.WebhookPayload{
		ID:      reg.ID,
		Country: reg.Country,
		Event:   "DELETE",
		Time:    time.Now().Format("20060102 15:04"),
	})

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
