package clients

import (
	"assignment_2/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TODO create folder/file for consts - var for testing purposes
var restCountriesAPI = "http://129.241.150.113:8080/v3.1/"

// GetCountry fetches data from the Countries API by ISO code
func GetCountry(countryCode string) (*models.Country, error) {
	countryCode = strings.ToUpper(countryCode)

	url := restCountriesAPI + "alpha/" + countryCode

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch country data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("country not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rawCountries []map[string]interface{}
	if err := json.Unmarshal(body, &rawCountries); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	if len(rawCountries) == 0 {
		return nil, errors.New("no country data found")
	}

	return mapToCountry(rawCountries[0])
}

// mapToCountry maps the raw REST Countries API response to our Country model
func mapToCountry(data map[string]interface{}) (*models.Country, error) {
	country := &models.Country{}

	if name, ok := data["name"].(map[string]interface{}); ok {
		if common, ok := name["common"].(string); ok {
			country.Name = common
		}
	}

	if code, ok := data["cca2"].(string); ok {
		country.Code = code
	}

	if capitals, ok := data["capital"].([]interface{}); ok && len(capitals) > 0 {
		if capital, ok := capitals[0].(string); ok {
			country.Capital = capital
		}
	}

	if population, ok := data["population"].(float64); ok {
		country.Population = int64(population)
	}

	if area, ok := data["area"].(float64); ok {
		country.Area = area
	}

	if region, ok := data["region"].(string); ok {
		country.Region = region
	}

	if languages, ok := data["languages"].(map[string]interface{}); ok {
		for _, lang := range languages {
			if langStr, ok := lang.(string); ok {
				country.Languages = append(country.Languages, langStr)
			}
		}
	}

	if currencies, ok := data["currencies"].(map[string]interface{}); ok {
		for code := range currencies {
			country.Currencies = append(country.Currencies, code)
		}
	}

	if borders, ok := data["borders"].([]interface{}); ok {
		for _, border := range borders {
			if borderStr, ok := border.(string); ok {
				country.Borders = append(country.Borders, borderStr)
			}
		}
	}

	// extract latlng coordinates (used by air quality client)
	if latlng, ok := data["latlng"].([]interface{}); ok && len(latlng) == 2 {
		if lat, ok := latlng[0].(float64); ok {
			country.Coordinates.Latitude = lat
		}
		if lng, ok := latlng[1].(float64); ok {
			country.Coordinates.Longitude = lng
		}
	}

	return country, nil
}
