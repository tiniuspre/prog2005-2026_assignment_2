package handlers

import (
	"assignment_2/internal/models"
	"context"
)

// Store defines the interface for data storage operations used by the handlers, and the local memory store for tests.
type Store interface {
	CreateNotification(ctx context.Context, reg models.NotificationRegistration) (string, error)
	GetNotification(ctx context.Context, id string) (*models.NotificationRegistration, error)
	ListNotifications(ctx context.Context) ([]models.NotificationRegistration, error)
	DeleteNotification(ctx context.Context, id string) error

	CreateRegistration(ctx context.Context, reg models.Registration) (string, error)
	GetRegistration(ctx context.Context, id string) (*models.Registration, error)
	ListRegistrations(ctx context.Context) ([]models.Registration, error)
	UpdateRegistration(ctx context.Context, reg models.Registration) error
	DeleteRegistration(ctx context.Context, id string) error
}
