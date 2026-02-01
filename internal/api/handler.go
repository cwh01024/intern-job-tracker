package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"intern-job-tracker/internal/model"
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
	jobRepo     *repository.JobRepository
	companyRepo *repository.CompanyRepository
	runLogRepo  *repository.RunLogRepository
	scheduler   SchedulerRunner
}

// NewHandler creates a new API handler.
func NewHandler(jobRepo *repository.JobRepository, companyRepo *repository.CompanyRepository, runLogRepo *repository.RunLogRepository, scheduler SchedulerRunner) *Handler {
	return &Handler{
		jobRepo:     jobRepo,
		companyRepo: companyRepo,
		runLogRepo:  runLogRepo,
		scheduler:   scheduler,
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
		// Jobs
		r.Get("/jobs", h.listJobs)
		r.Get("/jobs/{id}", h.getJob)

		// Companies
		r.Get("/companies", h.listCompanies)
		r.Post("/companies", h.createCompany)
		r.Put("/companies/{id}", h.updateCompany)
		r.Delete("/companies/{id}", h.deleteCompany)

		// Metrics & Stats
		r.Get("/stats", h.getStats)
		r.Get("/metrics", h.getMetrics)
		r.Get("/logs", h.getRunLogs)

		// Actions
		r.Post("/refresh", h.triggerRefresh)
	})

	// Serve static files
	r.Handle("/*", http.FileServer(http.Dir("web")))

	return r
}

func (h *Handler) listJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobRepo.GetAll()
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

	job, err := h.jobRepo.GetByID(id)
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

func (h *Handler) listCompanies(w http.ResponseWriter, r *http.Request) {
	if h.companyRepo == nil {
		respondJSON(w, []interface{}{})
		return
	}
	companies, err := h.companyRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, companies)
}

func (h *Handler) createCompany(w http.ResponseWriter, r *http.Request) {
	if h.companyRepo == nil {
		http.Error(w, "company management not available", http.StatusServiceUnavailable)
		return
	}

	var company model.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if company.Name == "" || company.CareerURL == "" {
		http.Error(w, "name and career_url are required", http.StatusBadRequest)
		return
	}

	if company.SearchTerm == "" {
		company.SearchTerm = "intern"
	}
	company.Enabled = true

	if err := h.companyRepo.Create(&company); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, company)
}

func (h *Handler) updateCompany(w http.ResponseWriter, r *http.Request) {
	if h.companyRepo == nil {
		http.Error(w, "company management not available", http.StatusServiceUnavailable)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var company model.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	company.ID = id
	if err := h.companyRepo.Update(&company); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, company)
}

func (h *Handler) deleteCompany(w http.ResponseWriter, r *http.Request) {
	if h.companyRepo == nil {
		http.Error(w, "company management not available", http.StatusServiceUnavailable)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.companyRepo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	byCompany := make(map[string]int)
	notified := 0
	for _, job := range jobs {
		byCompany[job.Company]++
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

func (h *Handler) getMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := make(map[string]interface{})

	// Job stats
	jobs, _ := h.jobRepo.GetAll()
	metrics["jobs"] = map[string]interface{}{
		"total": len(jobs),
	}

	// Company stats
	if h.companyRepo != nil {
		companies, _ := h.companyRepo.GetAll()
		enabled := 0
		for _, c := range companies {
			if c.Enabled {
				enabled++
			}
		}
		metrics["companies"] = map[string]interface{}{
			"total":   len(companies),
			"enabled": enabled,
		}
	}

	// Run log stats
	if h.runLogRepo != nil {
		runStats, _ := h.runLogRepo.GetStats()
		metrics["runs"] = runStats
	}

	respondJSON(w, metrics)
}

func (h *Handler) getRunLogs(w http.ResponseWriter, r *http.Request) {
	if h.runLogRepo == nil {
		respondJSON(w, []interface{}{})
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	logs, err := h.runLogRepo.GetRecent(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, logs)
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
