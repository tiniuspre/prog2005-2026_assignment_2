package firebase

import (
	"assignment_2/internal/models"
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const registrationsCollection = "registrations"

func CreateRegistration(ctx context.Context, client *firestore.Client, reg models.Registration) (string, error) {
	doc := client.Collection(registrationsCollection).NewDoc()
	reg.ID = doc.ID
	_, err := doc.Set(ctx, reg)
	if err != nil {
		return "", fmt.Errorf("failed to create registration: %w", err)
	}
	return doc.ID, nil
}

func GetRegistration(ctx context.Context, client *firestore.Client, id string) (*models.Registration, error) {
	doc, err := client.Collection(registrationsCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get registration: %w", err)
	}

	var reg models.Registration
	if err := doc.DataTo(&reg); err != nil {
		return nil, fmt.Errorf("failed to decode registration: %w", err)
	}
	return &reg, nil
}

func ListRegistrations(ctx context.Context, client *firestore.Client) ([]models.Registration, error) {
	iter := client.Collection(registrationsCollection).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list registrations: %w", err)
	}

	result := make([]models.Registration, 0, len(docs))
	for _, doc := range docs {
		var reg models.Registration
		if err := doc.DataTo(&reg); err != nil {
			return nil, fmt.Errorf("failed to decode registration: %w", err)
		}
		result = append(result, reg)
	}
	return result, nil
}

func UpdateRegistration(ctx context.Context, client *firestore.Client, reg models.Registration) error {
	_, err := client.Collection(registrationsCollection).Doc(reg.ID).Set(ctx, reg)
	if err != nil {
		return fmt.Errorf("failed to update registration: %w", err)
	}
	return nil
}

func DeleteRegistration(ctx context.Context, client *firestore.Client, id string) error {
	_, err := client.Collection(registrationsCollection).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete registration: %w", err)
	}
	return nil
}
