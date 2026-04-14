package handlers

import (
	"assignment_2/internal/models"
	"context"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const version = "v1"

// Assuming these exist elsewhere in your package:
// var probeFn = healthCheck
// var startTime = time.Now()
// var store NotificationStore

// healthCheck performs a request and returns the HTTP status code.
// If the service is unreachable, it returns 503.
func healthCheck(url, userAgent string) int {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return http.StatusServiceUnavailable
	}

	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	// OpenAQ and the course-hosted services do not reliably support HEAD.
	if strings.Contains(url, "api.openaq.org") {
		req.Header.Set("X-API-Key", os.Getenv("OPENAQ_API_KEY"))
		req.Method = http.MethodGet
	}
	if strings.Contains(url, "129.241.150.113") {
		req.Method = http.MethodGet
	}

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusServiceUnavailable
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode
}

func StatusHandler(w http.ResponseWriter, _ *http.Request) {
	checks := []struct {
		key       string
		url       string
		userAgent string
	}{
		{"countries_api", "http://129.241.150.113:8080/v3.1/", ""},
		{"meteo_api", "https://api.open-meteo.com/v1/forecast", ""},
		{"openaq_api", "https://api.openaq.org/v3/", ""},
		{"nominatim_api", "https://nominatim.openstreetmap.org/", "prog2005-assignment2/1.0"},
		{"currency_api", "http://129.241.150.113:9090/currency/", ""},
	}

	results := make(map[string]int, len(checks))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, c := range checks {
		wg.Add(1)
		go func(key, url, ua string) {
			defer wg.Done()
			code := probeFn(url, ua)

			mu.Lock()
			results[key] = code
			mu.Unlock()
		}(c.key, c.url, c.userAgent)
	}

	var dbStatus, webhookCount int
	wg.Add(1)
	go func() {
		defer wg.Done()

		regs, err := store.ListNotifications(context.Background())
		if err != nil {
			dbStatus = http.StatusServiceUnavailable
			webhookCount = 0
			return
		}

		dbStatus = http.StatusOK
		webhookCount = len(regs)
	}()

	wg.Wait()

	resp := models.StatusResponse{
		CountriesAPI:   results["countries_api"],
		MeteoAPI:       results["meteo_api"],
		OpenAQAPI:      results["openaq_api"],
		NominatimAPI:   results["nominatim_api"],
		CurrencyAPI:    results["currency_api"],
		NotificationDB: dbStatus,
		Webhooks:       webhookCount,
		Version:        version,
		Uptime:         int(time.Since(startTime).Seconds()),
	}

	overallStatus := http.StatusOK
	if resp.CountriesAPI != http.StatusOK ||
		resp.MeteoAPI != http.StatusOK ||
		resp.OpenAQAPI != http.StatusOK ||
		resp.NominatimAPI != http.StatusOK ||
		resp.CurrencyAPI != http.StatusOK ||
		resp.NotificationDB != http.StatusOK {
		overallStatus = http.StatusInternalServerError
	}

	writeJSON(w, overallStatus, resp)
}
