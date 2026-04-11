package handlers

import (
	"assignment_2/internal/models"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

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

	u, err := url.Parse(reg.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		writeError(w, http.StatusBadRequest, "url must be a valid http or https URL")
		return
	}

	if reg.Event == "THRESHOLD" && reg.Threshold == nil {
		writeError(w, http.StatusBadRequest, "threshold field required for THRESHOLD events")
		return
	}

	id, err := store.CreateNotification(context.Background(), reg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create notification")
		return
	}
	reg.ID = id

	writeJSON(w, http.StatusCreated, reg)
}

func GetNotificationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reg, err := store.GetNotification(context.Background(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get notification")
		return
	}
	if reg == nil {
		writeError(w, http.StatusNotFound, "notification not found")
		return
	}

	writeJSON(w, http.StatusOK, reg)
}

func ListNotificationsHandler(w http.ResponseWriter, _ *http.Request) {
	regs, err := store.ListNotifications(context.Background())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list notifications")
		return
	}

	writeJSON(w, http.StatusOK, regs)
}

func DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reg, err := store.GetNotification(context.Background(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get notification")
		return
	}
	if reg == nil {
		writeError(w, http.StatusNotFound, "notification not found")
		return
	}

	if err := store.DeleteNotification(context.Background(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete notification")
		return
	}

	go Send(reg.URL, models.WebhookPayload{
		ID:      reg.ID,
		Country: reg.Country,
		Event:   "DELETE",
		Time:    time.Now().Format("20060102 15:04"),
	})

	w.WriteHeader(http.StatusNoContent)
}
