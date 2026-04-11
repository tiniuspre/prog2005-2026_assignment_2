package handlers

import (
	"assignment_2/internal/models"
	"context"
	"errors"
	"fmt"
	"time"
)

type MemoryStore struct {
	notifications map[string]models.NotificationRegistration
	registrations map[string]models.Registration
	counter       int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		notifications: map[string]models.NotificationRegistration{},
		registrations: map[string]models.Registration{},
	}
}

func (m *MemoryStore) nextID() string {
	m.counter++
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), m.counter)
}

func (m *MemoryStore) CreateNotification(_ context.Context, reg models.NotificationRegistration) (string, error) {
	reg.ID = m.nextID()
	m.notifications[reg.ID] = reg
	return reg.ID, nil
}

func (m *MemoryStore) GetNotification(_ context.Context, id string) (*models.NotificationRegistration, error) {
	reg, ok := m.notifications[id]
	if !ok {
		return nil, nil
	}
	return &reg, nil
}

func (m *MemoryStore) ListNotifications(_ context.Context) ([]models.NotificationRegistration, error) {
	result := make([]models.NotificationRegistration, 0, len(m.notifications))
	for _, reg := range m.notifications {
		result = append(result, reg)
	}
	return result, nil
}

func (m *MemoryStore) DeleteNotification(_ context.Context, id string) error {
	if _, ok := m.notifications[id]; !ok {
		return errors.New("not found")
	}
	delete(m.notifications, id)
	return nil
}

func (m *MemoryStore) CreateRegistration(_ context.Context, reg models.Registration) (string, error) {
	reg.ID = m.nextID()
	m.registrations[reg.ID] = reg
	return reg.ID, nil
}

func (m *MemoryStore) GetRegistration(_ context.Context, id string) (*models.Registration, error) {
	reg, ok := m.registrations[id]
	if !ok {
		return nil, nil
	}
	return &reg, nil
}

func (m *MemoryStore) ListRegistrations(_ context.Context) ([]models.Registration, error) {
	result := make([]models.Registration, 0, len(m.registrations))
	for _, reg := range m.registrations {
		result = append(result, reg)
	}
	return result, nil
}

func (m *MemoryStore) UpdateRegistration(_ context.Context, reg models.Registration) error {
	if _, ok := m.registrations[reg.ID]; !ok {
		return errors.New("not found")
	}
	m.registrations[reg.ID] = reg
	return nil
}

func (m *MemoryStore) DeleteRegistration(_ context.Context, id string) error {
	if _, ok := m.registrations[id]; !ok {
		return errors.New("not found")
	}
	delete(m.registrations, id)
	return nil
}
