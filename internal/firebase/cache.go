package firebase

import (
	"context"
	"time"

	"assignment_2/internal/models"

	"cloud.google.com/go/firestore"
)

const cacheExpiry = 24 * time.Hour
const cacheCollection = "cache"

type cachedCountry struct {
	Country  models.Country `firestore:"country"`
	CachedAt time.Time      `firestore:"cachedAt"`
}

// GetCachedCountry retrieves a country from the cache, returns nil if missing or stale
func GetCachedCountry(ctx context.Context, client *firestore.Client, isoCode string) (*models.Country, error) {
	doc, err := client.Collection(cacheCollection).Doc(isoCode).Get(ctx)
	if err != nil {
		// Document doesn't exist yet — cache miss
		return nil, nil
	}

	var cached cachedCountry
	if err := doc.DataTo(&cached); err != nil {
		return nil, err
	}

	// Check if cache is stale
	if time.Since(cached.CachedAt) > cacheExpiry {
		return nil, nil
	}

	return &cached.Country, nil
}

// SetCachedCountry stores a country in the cache with a timestamp
func SetCachedCountry(ctx context.Context, client *firestore.Client, isoCode string, country *models.Country) error {
	_, err := client.Collection(cacheCollection).Doc(isoCode).Set(ctx, cachedCountry{
		Country:  *country,
		CachedAt: time.Now(),
	})
	return err
}
