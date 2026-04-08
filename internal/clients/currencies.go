package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TODO create folder/file for consts - var for testing purposes
var currencyAPI = "http://129.241.150.113:9090/currency/"

// GetExchangeRates fetches exchange rates from a base currency to a list of target currencies
func GetExchangeRates(baseCurrency string, targetCurrencies []string) (map[string]float64, error) {
	url := currencyAPI + strings.ToLower(baseCurrency)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("currency API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rateResponse map[string]interface{}
	if err := json.Unmarshal(body, &rateResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	rates, ok := rateResponse["rates"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from currency API")
	}

	result := make(map[string]float64)
	for _, target := range targetCurrencies {
		key := strings.ToUpper(target)
		if rate, ok := rates[key].(float64); ok {
			result[key] = rate
		}
	}

	return result, nil
}
