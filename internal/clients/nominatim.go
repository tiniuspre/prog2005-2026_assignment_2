package clients

import (
	"assignment_2/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type nominatimResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// GetCapitalCoordinates returns precise coordinates for a capital city using Nominatim
func GetCapitalCoordinates(capital, country string) (*models.Coordinates, error) {
	query := url.Values{}
	query.Set("q", capital+", "+country)
	query.Set("format", "json")
	query.Set("limit", "1")

	fullURL := nominatimBaseURL + "/search?" + query.Encode()

	// http.NewRequest instead of http.Get, we need to set a header which isn't possible with .Get∫
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Required by Nominatim's usage policy
	req.Header.Set("User-Agent", "prog2005-assignment2/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch coordinates: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nominatim returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var results []nominatimResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no location found for capital %s", capital)
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude value: %w", err)
	}

	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude value: %w", err)
	}

	return &models.Coordinates{
		Latitude:  lat,
		Longitude: lon,
	}, nil
}
