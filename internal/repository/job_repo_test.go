package repository

import (
	"database/sql"
	"os"
	"testing"

	"intern-job-tracker/internal/db"
	"intern-job-tracker/internal/model"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "test_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	database, err := db.New(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to create database: %v", err)
	}

	cleanup := func() {
		database.Close()
		os.Remove(tmpFile.Name())
	}

	return database, cleanup
}

func TestJobRepository_Create(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewJobRepository(database)

	job := &model.Job{
		Company:  "Google",
		Title:    "Software Engineering Intern",
		URL:      "https://google.com/job/123",
		Location: "Mountain View, CA",
	}

	err := repo.Create(job)
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	if job.ID == 0 {
		t.Error("expected job ID to be set after create")
	}
}

func TestJobRepository_GetByURL(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewJobRepository(database)

	// Create a job first
	job := &model.Job{
		Company:  "Amazon",
		Title:    "SDE Intern",
		URL:      "https://amazon.jobs/123",
		Location: "Seattle, WA",
	}
	repo.Create(job)

	// Get by URL
	found, err := repo.GetByURL("https://amazon.jobs/123")
	if err != nil {
		t.Fatalf("failed to get job by URL: %v", err)
	}

	if found == nil {
		t.Fatal("expected to find job, got nil")
	}

	if found.Company != "Amazon" {
		t.Errorf("expected Amazon, got %s", found.Company)
	}

	// Test not found
	notFound, err := repo.GetByURL("https://nonexistent.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent URL")
	}
}

func TestJobRepository_GetUnnotified(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewJobRepository(database)

	// Create jobs - one notified, one not
	job1 := &model.Job{
		Company: "Google",
		Title:   "Intern 1",
		URL:     "https://google.com/1",
	}
	job2 := &model.Job{
		Company: "Amazon",
		Title:   "Intern 2",
		URL:     "https://amazon.com/2",
	}

	repo.Create(job1)
	repo.Create(job2)
	repo.MarkNotified(job1.ID)

	// Get unnotified
	unnotified, err := repo.GetUnnotified()
	if err != nil {
		t.Fatalf("failed to get unnotified: %v", err)
	}

	if len(unnotified) != 1 {
		t.Errorf("expected 1 unnotified job, got %d", len(unnotified))
	}

	if unnotified[0].Company != "Amazon" {
		t.Errorf("expected Amazon, got %s", unnotified[0].Company)
	}
}

func TestJobRepository_MarkNotified(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewJobRepository(database)

	job := &model.Job{
		Company: "Uber",
		Title:   "Intern",
		URL:     "https://uber.com/1",
	}
	repo.Create(job)

	err := repo.MarkNotified(job.ID)
	if err != nil {
		t.Fatalf("failed to mark notified: %v", err)
	}

	// Verify it's marked
	found, _ := repo.GetByURL("https://uber.com/1")
	if !found.Notified {
		t.Error("expected job to be marked as notified")
	}
}

func TestJobRepository_GetAll(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewJobRepository(database)

	// Create multiple jobs
	for i := 0; i < 3; i++ {
		job := &model.Job{
			Company: "Company",
			Title:   "Intern",
			URL:     "https://example.com/" + string(rune('a'+i)),
		}
		repo.Create(job)
	}

	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("failed to get all: %v", err)
	}

	if len(all) != 3 {
		t.Errorf("expected 3 jobs, got %d", len(all))
	}
}
