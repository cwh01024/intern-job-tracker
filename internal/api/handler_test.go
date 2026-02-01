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

	jobRepo := repository.NewJobRepository(database)
	companyRepo := repository.NewCompanyRepository(database)
	runLogRepo := repository.NewRunLogRepository(database)
	handler := NewHandler(jobRepo, companyRepo, runLogRepo, nil)

	cleanup := func() {
		database.Close()
		os.Remove(tmpFile.Name())
	}

	return handler, cleanup
}

func TestAPI_ListJobs(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	handler.jobRepo.Create(&model.Job{Company: "Google", Title: "Intern 1", URL: "https://google.com/1"})
	handler.jobRepo.Create(&model.Job{Company: "Amazon", Title: "Intern 2", URL: "https://amazon.com/2"})

	req := httptest.NewRequest("GET", "/api/jobs", nil)
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var jobs []*model.Job
	json.NewDecoder(w.Body).Decode(&jobs)

	if len(jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(jobs))
	}
}

func TestAPI_GetJob(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	job := &model.Job{Company: "Uber", Title: "SDE Intern", URL: "https://uber.com/1"}
	handler.jobRepo.Create(job)

	req := httptest.NewRequest("GET", "/api/jobs/1", nil)
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestAPI_ListCompanies(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/companies", nil)
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestAPI_GetMetrics(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	handler.jobRepo.Create(&model.Job{Company: "Google", Title: "Intern", URL: "https://google.com/1", DiscoveredAt: time.Now()})

	req := httptest.NewRequest("GET", "/api/metrics", nil)
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var metrics map[string]interface{}
	json.NewDecoder(w.Body).Decode(&metrics)

	if metrics["jobs"] == nil {
		t.Error("expected jobs in metrics")
	}
}

func TestAPI_GetLogs(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/logs", nil)
	w := httptest.NewRecorder()

	handler.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
