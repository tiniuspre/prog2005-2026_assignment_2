package handlers

import (
	"assignment_2/internal/firebase"
	"assignment_2/internal/models"
	"context"

	"cloud.google.com/go/firestore"
)

// Actual Firestore storage implementation to be used in prod.
type FirestoreStore struct {
	client *firestore.Client
}

func NewFirestoreStore(client *firestore.Client) *FirestoreStore {
	return &FirestoreStore{client: client}
}

func (s *FirestoreStore) CreateNotification(ctx context.Context, reg models.NotificationRegistration) (string, error) {
	return firebase.CreateNotification(ctx, s.client, reg)
}

func (s *FirestoreStore) GetNotification(ctx context.Context, id string) (*models.NotificationRegistration, error) {
	return firebase.GetNotification(ctx, s.client, id)
}

func (s *FirestoreStore) ListNotifications(ctx context.Context) ([]models.NotificationRegistration, error) {
	return firebase.ListNotifications(ctx, s.client)
}

func (s *FirestoreStore) DeleteNotification(ctx context.Context, id string) error {
	return firebase.DeleteNotification(ctx, s.client, id)
}

func (s *FirestoreStore) CreateRegistration(ctx context.Context, reg models.Registration) (string, error) {
	return firebase.CreateRegistration(ctx, s.client, reg)
}

func (s *FirestoreStore) GetRegistration(ctx context.Context, id string) (*models.Registration, error) {
	return firebase.GetRegistration(ctx, s.client, id)
}

func (s *FirestoreStore) ListRegistrations(ctx context.Context) ([]models.Registration, error) {
	return firebase.ListRegistrations(ctx, s.client)
}

func (s *FirestoreStore) UpdateRegistration(ctx context.Context, reg models.Registration) error {
	return firebase.UpdateRegistration(ctx, s.client, reg)
}

func (s *FirestoreStore) DeleteRegistration(ctx context.Context, id string) error {
	return firebase.DeleteRegistration(ctx, s.client, id)
}
