package handlers

import (
	"assignment_2/internal/models"
	"bytes"
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
	for _, reg := range notifications {
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
	for _, reg := range notifications {
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
