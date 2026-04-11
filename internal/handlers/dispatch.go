package handlers

import (
	"assignment_2/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func Send(url string, payload models.WebhookPayload) {
	body, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("webhook dispatch failed to %s: %v", url, err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()
	log.Printf("webhook delivered to %s: %d", url, resp.StatusCode)
}

func DispatchEvent(event, country string) {
	regs, err := store.ListNotifications(context.Background())
	if err != nil {
		log.Printf("DispatchEvent: failed to list notifications: %v", err)
		return
	}

	for _, reg := range regs {
		if reg.Event != event {
			continue
		}
		if reg.Country != "" && reg.Country != country {
			continue
		}
		go Send(reg.URL, models.WebhookPayload{
			ID:      reg.ID,
			Country: country,
			Event:   event,
			Time:    time.Now().Format("20060102 15:04"),
		})
	}
}

func CheckThresholds(country string, values map[string]float64) {
	regs, err := store.ListNotifications(context.Background())
	if err != nil {
		log.Printf("CheckThresholds: failed to list notifications: %v", err)
		return
	}

	for _, reg := range regs {
		if reg.Event != "THRESHOLD" || reg.Threshold == nil {
			continue
		}
		if reg.Country != "" && reg.Country != country {
			continue
		}
		val, ok := values[reg.Threshold.Field]
		if !ok {
			continue
		}
		if evaluate(val, reg.Threshold.Operator, reg.Threshold.Value) {
			go Send(reg.URL, models.WebhookPayload{
				ID:      reg.ID,
				Country: country,
				Event:   "THRESHOLD",
				Time:    time.Now().Format("20060102 15:04"),
				Details: &models.ThresholdDetails{
					Field:         reg.Threshold.Field,
					Operator:      reg.Threshold.Operator,
					Threshold:     reg.Threshold.Value,
					MeasuredValue: val,
				},
			})
		}
	}
}

func evaluate(measured float64, operator string, threshold float64) bool {
	switch operator {
	case ">=":
		return measured >= threshold
	case "<=":
		return measured <= threshold
	default:
		return false
	}
}
