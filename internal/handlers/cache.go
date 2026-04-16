package handlers

import (
	"assignment_2/internal/clients"
	"assignment_2/internal/firebase"
	"assignment_2/internal/models"
	"context"
	"log"
)

// cachedGetCountry checks Firestore cache before calling the REST Countries API.
// Cache entries expire after 24 hours as defined in firebase/cache.go.
func cachedGetCountry(isoCode string) (*models.Country, error) {
	ctx := context.Background()

	// check cache first
	country, err := firebase.GetCachedCountry(ctx, firestoreClient, isoCode)
	if err == nil && country != nil {
		return country, nil
	}

	// cache miss — fetch from API
	country, err = clients.GetCountry(isoCode)
	if err != nil {
		return nil, err
	}

	// store in cache, non-fatal if it fails
	if cacheErr := firebase.SetCachedCountry(ctx, firestoreClient, isoCode, country); cacheErr != nil {
		log.Printf("warning: failed to cache country: %v", cacheErr)
	}
	return country, nil
}
