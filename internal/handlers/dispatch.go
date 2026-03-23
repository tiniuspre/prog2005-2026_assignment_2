package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func Send(url string, payload WebhookPayload) {
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
