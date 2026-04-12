package firebase

import (
	"assignment_2/internal/models"
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const notificationsCollection = "notifications"

func CreateNotification(ctx context.Context, client *firestore.Client, reg models.NotificationRegistration) (string, error) {
	doc := client.Collection(notificationsCollection).NewDoc()
	reg.ID = doc.ID
	_, err := doc.Set(ctx, reg)
	if err != nil {
		return "", fmt.Errorf("failed to create notification: %w", err)
	}
	return doc.ID, nil
}

func GetNotification(ctx context.Context, client *firestore.Client, id string) (*models.NotificationRegistration, error) {
	doc, err := client.Collection(notificationsCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	var reg models.NotificationRegistration
	if err := doc.DataTo(&reg); err != nil {
		return nil, fmt.Errorf("failed to decode notification: %w", err)
	}
	return &reg, nil
}

func ListNotifications(ctx context.Context, client *firestore.Client) ([]models.NotificationRegistration, error) {
	iter := client.Collection(notificationsCollection).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	result := make([]models.NotificationRegistration, 0, len(docs))
	for _, doc := range docs {
		var reg models.NotificationRegistration
		if err := doc.DataTo(&reg); err != nil {
			return nil, fmt.Errorf("failed to decode notification: %w", err)
		}
		result = append(result, reg)
	}
	return result, nil
}

func DeleteNotification(ctx context.Context, client *firestore.Client, id string) error {
	_, err := client.Collection(notificationsCollection).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}
