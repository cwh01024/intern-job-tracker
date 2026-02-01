package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"intern-job-tracker/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SchedulerRunner interface for triggering manual refresh.
type SchedulerRunner interface {
	RunNow() error
}

// Handler manages HTTP API endpoints.
type Handler struct {
	repo      *repository.JobRepository
	scheduler SchedulerRunner
}

// NewHandler creates a new API handler.
func NewHandler(repo *repository.JobRepository, scheduler SchedulerRunner) *Handler {
	return &Handler{
		repo:      repo,
		scheduler: scheduler,
	}
}

// Router returns the configured chi router.
func (h *Handler) Router() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/jobs", h.listJobs)
		r.Get("/jobs/{id}", h.getJob)
		r.Get("/stats", h.getStats)
		r.Post("/refresh", h.triggerRefresh)
	})

	// Serve static files
	r.Handle("/*", http.FileServer(http.Dir("web")))

	return r
}

func (h *Handler) listJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, jobs)
}

func (h *Handler) getJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	job, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if job == nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	respondJSON(w, job)
}

func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Count by company
	byCompany := make(map[string]int)
	for _, job := range jobs {
		byCompany[job.Company]++
	}

	notified := 0
	for _, job := range jobs {
		if job.Notified {
			notified++
		}
	}

	stats := map[string]interface{}{
		"total_jobs": len(jobs),
		"notified":   notified,
		"by_company": byCompany,
	}

	respondJSON(w, stats)
}

func (h *Handler) triggerRefresh(w http.ResponseWriter, r *http.Request) {
	if h.scheduler == nil {
		http.Error(w, "scheduler not configured", http.StatusServiceUnavailable)
		return
	}

	if err := h.scheduler.RunNow(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"status": "ok", "message": "refresh triggered"})
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
