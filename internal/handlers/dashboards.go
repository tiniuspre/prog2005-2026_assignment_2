package handlers

import (
	"assignment_2/internal/clients"
	"assignment_2/internal/firebase"
	"assignment_2/internal/models"
	"context"
	"log"
	"net/http"
	"time"
)

func GetDashboardHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reg, err := store.GetRegistration(context.Background(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get registration")
		return
	}
	if reg == nil {
		writeError(w, http.StatusNotFound, "registration not found")
		return
	}

	// Check if country is cached
	country, err := firebase.GetCachedCountry(r.Context(), firestoreClient, reg.IsoCode)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read cache")
		return
	}
	if country == nil {
		country, err = clients.GetCountry(reg.IsoCode)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to fetch country data")
			return
		}
		// Cache new country
		if err := firebase.SetCachedCountry(r.Context(), firestoreClient, reg.IsoCode, country); err != nil {
			log.Printf("warning: failed to cache country data: %v", err)
		}
	}

	features := models.DashboardFeatures{}
	thresholdValues := map[string]float64{}

	if reg.Features.Temperature || reg.Features.Precipitation {
		weather, err := clients.GetWeather(country.Coordinates.Latitude,
			country.Coordinates.Longitude)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "weather fetch failed: "+err.Error())
			return
		}
		if reg.Features.Temperature {
			features.Temperature = &weather.Temperature
			thresholdValues["temperature"] = weather.Temperature
		}
		if reg.Features.Precipitation {
			features.Precipitation = &weather.Precipitation
			thresholdValues["precipitation"] = weather.Precipitation
		}
	}

	if reg.Features.AirQuality {
		aq, err := clients.GetAirQuality(country.Coordinates.Latitude,
			country.Coordinates.Longitude)
		if err == nil {
			features.AirQuality = aq // already *models.AirQualityData
			thresholdValues["pm25"] = aq.PM25
			thresholdValues["pm10"] = aq.PM10
		}
	}

	if reg.Features.Capital {
		features.Capital = &country.Capital
	}
	if reg.Features.Coordinates {
		features.Coordinates = &country.Coordinates
	}
	if reg.Features.Population {
		features.Population = &country.Population
	}
	if reg.Features.Area {
		features.Area = &country.Area
	}

	if len(reg.Features.TargetCurrencies) > 0 && len(country.Currencies) > 0 {
		rates, err := clients.GetExchangeRates(country.Currencies[0],
			reg.Features.TargetCurrencies)
		if err == nil {
			features.TargetCurrencies = rates
		}
	}

	resp := models.DashboardResponse{
		Country:       reg.Country,
		IsoCode:       reg.IsoCode,
		Features:      features,
		LastRetrieval: time.Now().Format("20060102 15:04"),
	}

	writeJSON(w, http.StatusOK, resp)
	DispatchEvent("INVOKE", reg.IsoCode)
	CheckThresholds(reg.IsoCode, thresholdValues)
}
