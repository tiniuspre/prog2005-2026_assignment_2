package firebase

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const apiKeysCollection = "apikeys"

// Data represented that is stored in firebase
type APIKey struct {
	Key       string `firestore:"key"       json:"key"`
	Name      string `firestore:"name"      json:"name"`
	Email     string `firestore:"email"     json:"email"`
	CreatedAt string `firestore:"createdAt" json:"createdAt"`
}

// CreateAPIKey Creates a new API key
func CreateAPIKey(ctx context.Context, client *firestore.Client, ak APIKey) error {
	_, err := client.Collection(apiKeysCollection).Doc(ak.Key).Set(ctx, ak)
	if err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}
	return nil
}

// GetAPIKey Fetches api key, gives null if it does not exist
func GetAPIKey(ctx context.Context, client *firestore.Client, key string) (*APIKey, error) {
	doc, err := client.Collection(apiKeysCollection).Doc(key).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	var ak APIKey
	if err := doc.DataTo(&ak); err != nil {
		return nil, fmt.Errorf("failed to decode API key: %w", err)
	}
	return &ak, nil
}

// DeleteAPIKey Deltes the API key.
func DeleteAPIKey(ctx context.Context, client *firestore.Client, key string) (bool, error) {
	doc, err := client.Collection(apiKeysCollection).Doc(key).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check API key: %w", err)
	}

	if _, err := doc.Ref.Delete(ctx); err != nil {
		return false, fmt.Errorf("failed to delete API key: %w", err)
	}
	return true, nil
}

// FormatTimestamp returns the current time formatted as "20060102 15:04".
func FormatTimestamp() string {
	return time.Now().UTC().Format("20060102 15:04")
}
