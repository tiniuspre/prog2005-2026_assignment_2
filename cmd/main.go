package main

import (
	"assignment_2/internal/firebase"
	"assignment_2/internal/handlers"
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
	defer firestoreClient.Close()
	handlers.Init(firestoreClient)

	mux := http.NewServeMux()

	// registration routes
	mux.HandleFunc("POST /envdash/v1/registrations/", handlers.CreateRegistrationHandler)
	mux.HandleFunc("GET /envdash/v1/registrations/{id}", handlers.GetRegistrationHandler)
	mux.HandleFunc("GET /envdash/v1/registrations/", handlers.ListRegistrationsHandler)
	mux.HandleFunc("PUT /envdash/v1/registrations/{id}", handlers.UpdateRegistrationHandler)
	mux.HandleFunc("DELETE /envdash/v1/registrations/{id}", handlers.DeleteRegistrationHandler)

	// dashboard routes
	mux.HandleFunc("GET /envdash/v1/dashboards/{id}", handlers.GetDashboardHandler)

	// notification routes
	mux.HandleFunc("POST /envdash/v1/notifications/", handlers.CreateNotificationHandler)
	mux.HandleFunc("GET /envdash/v1/notifications/{id}", handlers.GetNotificationHandler)
	mux.HandleFunc("GET /envdash/v1/notifications/", handlers.ListNotificationsHandler)
	mux.HandleFunc("DELETE /envdash/v1/notifications/{id}", handlers.DeleteNotificationHandler)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
