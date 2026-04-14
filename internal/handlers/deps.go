package handlers

import (
	"assignment_2/internal/clients"
	"assignment_2/internal/firebase"
)

var (
	countryByISO     = clients.GetCountry
	countryByName    = clients.GetCountryByName
	weatherFor       = clients.GetWeather
	capitalCoordsFor = clients.GetCapitalCoordinates
	airQualityFor    = clients.GetAirQuality
	exchangeRatesFor = clients.GetExchangeRates

	generateKeyFn     = generateKey
	createAPIKeyFn    = firebase.CreateAPIKey
	deleteAPIKeyFn    = firebase.DeleteAPIKey
	formatTimestampFn = firebase.FormatTimestamp

	probeFn = healthCheck
)
