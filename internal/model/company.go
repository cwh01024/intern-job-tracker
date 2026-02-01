package model

import "time"

// Company represents a company to track for job listings.
type Company struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CareerURL  string    `json:"career_url"`
	SearchTerm string    `json:"search_term"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}

// RunLog represents a record of a job check execution.
type RunLog struct {
	ID                int64     `json:"id"`
	RunAt             time.Time `json:"run_at"`
	CompaniesChecked  int       `json:"companies_checked"`
	JobsFound         int       `json:"jobs_found"`
	NewJobs           int       `json:"new_jobs"`
	NotificationsSent int       `json:"notifications_sent"`
	DurationMs        int64     `json:"duration_ms"`
	Status            string    `json:"status"`
	ErrorMessage      string    `json:"error_message,omitempty"`
}
