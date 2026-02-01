package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"intern-job-tracker/internal/db"
	"intern-job-tracker/internal/model"
	"intern-job-tracker/internal/repository"
)

func setupTestAPI(t *testing.T) (*Handler, func()) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "test_api_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	database, err := db.New(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to create database: %v", err)
	}

	repo := repository.NewJobRepository(database)
	handler := NewHandler(repo, nil)

	cleanup := func() {
		database.Close()
		os.Remove(tmpFile.Name())
	}

	return handler, cleanup
}

func TestAPI_ListJobs(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	// Add some jobs
	handler.repo.Create(&model.Job{Company: "Google", Title: "Intern 1", URL: "https://google.com/1"})
	handler.repo.Create(&model.Job{Company: "Amazon", Title: "Intern 2", URL: "https://amazon.com/2"})

	req := httptest.NewRequest("GET", "/api/jobs", nil)
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var jobs []*model.Job
	if err := json.NewDecoder(w.Body).Decode(&jobs); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(jobs))
	}
}

func TestAPI_GetJob(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	job := &model.Job{Company: "Uber", Title: "SDE Intern", URL: "https://uber.com/1"}
	handler.repo.Create(job)

	req := httptest.NewRequest("GET", "/api/jobs/1", nil)
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result model.Job
	json.NewDecoder(w.Body).Decode(&result)

	if result.Company != "Uber" {
		t.Errorf("expected Uber, got %s", result.Company)
	}
}

func TestAPI_GetJob_NotFound(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/jobs/999", nil)
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestAPI_GetStats(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	// Add jobs
	handler.repo.Create(&model.Job{Company: "Google", Title: "Intern", URL: "https://google.com/1", DiscoveredAt: time.Now()})
	handler.repo.Create(&model.Job{Company: "Amazon", Title: "Intern", URL: "https://amazon.com/2", DiscoveredAt: time.Now()})

	req := httptest.NewRequest("GET", "/api/stats", nil)
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var stats map[string]interface{}
	json.NewDecoder(w.Body).Decode(&stats)

	if stats["total_jobs"].(float64) != 2 {
		t.Errorf("expected 2 total jobs")
	}
}

func TestAPI_CORS(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	req := httptest.NewRequest("OPTIONS", "/api/jobs", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	// Should have CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("expected CORS header")
	}
}
