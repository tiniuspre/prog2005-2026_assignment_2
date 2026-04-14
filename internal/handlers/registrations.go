package handlers

import (
	"assignment_2/internal/models"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

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
		country, err := countryByISO(reg.IsoCode)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid country iso code")
			return
		}
		reg.Country = country.Name
		reg.IsoCode = country.Code
	} else {
		country, err := countryByName(reg.Country)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid country name")
			return
		}
		reg.IsoCode = country.Code
	}

	reg.LastChange = time.Now().Format("20060102 15:04")

	id, err := store.CreateRegistration(context.Background(), reg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create registration")
		return
	}

	DispatchEvent("REGISTER", reg.IsoCode)

	writeJSON(w, http.StatusCreated, map[string]string{
		"id":         id,
		"lastChange": reg.LastChange,
	})
}

func GetRegistrationHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusOK, reg)
}

func ListRegistrationsHandler(w http.ResponseWriter, _ *http.Request) {
	regs, err := store.ListRegistrations(context.Background())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list registrations")
		return
	}

	writeJSON(w, http.StatusOK, regs)
}

func DeleteRegistrationHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := store.DeleteRegistration(context.Background(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete registration")
		return
	}

	DispatchEvent("DELETE", reg.IsoCode)

	w.WriteHeader(http.StatusNoContent)
}

func UpdateRegistrationHandler(w http.ResponseWriter, r *http.Request) {
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

	var update models.Registration
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if update.Country != "" {
		country, err := countryByName(update.Country)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid country name")
			return
		}
		reg.Country = country.Name
		reg.IsoCode = country.Code
	} else if update.IsoCode != "" {
		country, err := countryByISO(update.IsoCode)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid country iso code")
			return
		}
		reg.Country = country.Name
		reg.IsoCode = country.Code
	}
	reg.LastChange = time.Now().Format("20060102 15:04")
	reg.Features = update.Features

	if err := store.UpdateRegistration(context.Background(), *reg); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update registration")
		return
	}

	DispatchEvent("CHANGE", reg.IsoCode)

	w.WriteHeader(http.StatusOK)
}
