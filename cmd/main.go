package main

import (
	"assignment_2/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// routes
	mux.HandleFunc("POST /notifications/", handlers.CreateNotificationHandler)
	mux.HandleFunc("GET /notifications/{id}", handlers.GetNotificationHandler)
	mux.HandleFunc("GET /notifications/", handlers.ListNotificationsHandler)
	mux.HandleFunc("DELETE /notifications/{id}", handlers.DeleteNotificationHandler)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
