package main

import (
	"log"
	"net/http"

	"intern-job-tracker/internal/db"
)

func main() {
	// Initialize database
	database, err := db.New("jobs.db")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("Database initialized successfully")

	// TODO: Initialize other components (repository, scraper, notifier, scheduler, API)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
