package handlers

import (
	"assignment_2/internal/clients"
	"assignment_2/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var registrations = map[string]models.Registration{}

func CreateRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var reg models.Registration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if reg.Country == "" && reg.IsoCode == "" {
		writeError(w, http.StatusBadRequest, "country or isoCode are required")
		return
	}

	if reg.IsoCode != "" {
		country, err := clients.GetCountry(reg.IsoCode)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid country iso code")
			return
		}
		reg.Country = country.Name
		reg.IsoCode = country.Code
	} else {
		country, err := clients.GetCountryByName(reg.Country)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid country name")
			return
		}
		reg.IsoCode = country.Code
	}

	// temp storage (pre firebase)
	reg.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	reg.LastChange = time.Now().Format("20060102 15:04")

	// storing
	registrations[reg.ID] = reg
	DispatchEvent("REGISTER", reg.IsoCode)

	writeJSON(w, http.StatusCreated, map[string]string{
		"id":         reg.ID,
		"lastChange": reg.LastChange,
	})
}

func GetRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reg, ok := registrations[id]
	if !ok {
		writeError(w, http.StatusNotFound, "registration not found")
		return
	}

	writeJSON(w, http.StatusOK, reg)
}

func ListRegistrationsHandler(w http.ResponseWriter, _ *http.Request) {
	result := make([]models.Registration, 0, len(registrations))
	for _, reg := range registrations {
		result = append(result, reg)
	}

	writeJSON(w, http.StatusOK, result)
}

func DeleteRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reg, ok := registrations[id]
	if !ok {
		writeError(w, http.StatusNotFound, "registration not found")
		return
	}

	DispatchEvent("DELETE", reg.IsoCode)
	delete(registrations, id)

	w.WriteHeader(http.StatusNoContent) // not using writeJSON since no body
}

func UpdateRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reg, ok := registrations[id]
	if !ok {
		writeError(w, http.StatusNotFound, "registration not found")
		return
	}

	var update models.Registration
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if update.Country != "" {
		country, err := clients.GetCountryByName(update.Country)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid country name")
			return
		}
		reg.IsoCode = country.Code
	} else if update.IsoCode != "" {
		country, err := clients.GetCountry(update.IsoCode)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid country iso code")
			return
		}
		reg.Country = country.Name
		reg.IsoCode = country.Code
	}
	reg.LastChange = time.Now().Format("20060102 15:04")
	reg.Features = update.Features
	registrations[id] = reg
	DispatchEvent("CHANGE", reg.IsoCode)

	w.WriteHeader(http.StatusOK) // not using writeJSON since no body
}
