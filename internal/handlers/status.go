package handlers

import (
	"assignment_2/internal/models"
	"context"
	"net/http"
	"sync"
	"time"
)

const version = "v1"

// healthCheck performs a HEAD request and returns the HTTP status code, or 0 on network failure.
func healthCheck(url, userAgent string) int {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return 0
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode
}

func StatusHandler(w http.ResponseWriter, _ *http.Request) {
	type result struct {
		key   string
		value int
	}

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
			code := healthCheck(url, ua)
			mu.Lock()
			results[key] = code
			mu.Unlock()
		}(c.key, c.url, c.userAgent)
	}

	// Count webhooks and test Firestore connectivity in parallel with API checks
	wg.Add(1)
	var dbStatus, webhookCount int
	go func() {
		defer wg.Done()
		regs, err := store.ListNotifications(context.Background())
		if err != nil {
			dbStatus = http.StatusServiceUnavailable
			webhookCount = 0
		} else {
			dbStatus = http.StatusOK
			webhookCount = len(regs)
		}
	}()

	wg.Wait()

	writeJSON(w, http.StatusOK, models.StatusResponse{
		CountriesAPI:   results["countries_api"],
		MeteoAPI:       results["meteo_api"],
		OpenAQAPI:      results["openaq_api"],
		NominatimAPI:   results["nominatim_api"],
		CurrencyAPI:    results["currency_api"],
		NotificationDB: dbStatus,
		Webhooks:       webhookCount,
		Version:        version,
		Uptime:         int(time.Since(startTime).Seconds()),
	})
}
