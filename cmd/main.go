package main

import (
	"assignment_2/internal/firebase"
	"assignment_2/internal/handlers"
	"assignment_2/internal/middleware"
	"context"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	firestoreClient, err := firebase.NewFirestoreClient(ctx)
	if err != nil {
		log.Fatalf("failed to initialise Firestore: %v", err)
	}
	defer func() {
		if err := firestoreClient.Close(); err != nil {
			log.Printf("failed to close Firestore client: %v", err)
		}
	}()

	handlers.Init(firestoreClient)

	mux := http.NewServeMux()

	// auth routes
	mux.HandleFunc("POST /auth/", handlers.CreateAPIKeyHandler)
	mux.HandleFunc("DELETE /auth/{key}", handlers.DeleteAPIKeyHandler)

	// registration routes
	mux.HandleFunc("POST /envdash/v1/registrations/", handlers.CreateRegistrationHandler)
	mux.HandleFunc("GET /envdash/v1/registrations/{id}", handlers.GetRegistrationHandler)
	mux.HandleFunc("GET /envdash/v1/registrations/", handlers.ListRegistrationsHandler)
	mux.HandleFunc("PUT /envdash/v1/registrations/{id}", handlers.UpdateRegistrationHandler)
	mux.HandleFunc("DELETE /envdash/v1/registrations/{id}", handlers.DeleteRegistrationHandler)

	// status route
	mux.HandleFunc("GET /envdash/v1/status/", handlers.StatusHandler)

	// dashboard routes
	mux.HandleFunc("GET /envdash/v1/dashboards/{id}", handlers.GetDashboardHandler)

	// notification routes
	mux.HandleFunc("POST /envdash/v1/notifications/", handlers.CreateNotificationHandler)
	mux.HandleFunc("GET /envdash/v1/notifications/{id}", handlers.GetNotificationHandler)
	mux.HandleFunc("GET /envdash/v1/notifications/", handlers.ListNotificationsHandler)
	mux.HandleFunc("DELETE /envdash/v1/notifications/{id}", handlers.DeleteNotificationHandler)

	// Wraps routes to handle auth. Only auth /env../status in is public without auth.
	authMiddleware := middleware.APIKeyAuth(firestoreClient, []middleware.PublicRoute{
		{Method: http.MethodPost, Path: "/auth/"},
		{Method: http.MethodGet, Path: "/envdash/v1/status/"},
	})

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", authMiddleware(mux)))
}
