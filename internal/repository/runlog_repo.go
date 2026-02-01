package repository

import (
	"database/sql"

	"intern-job-tracker/internal/model"
)

// RunLogRepository handles database operations for run logs.
type RunLogRepository struct {
	db *sql.DB
}

// NewRunLogRepository creates a new RunLogRepository.
func NewRunLogRepository(db *sql.DB) *RunLogRepository {
	return &RunLogRepository{db: db}
}

// Create adds a new run log entry.
func (r *RunLogRepository) Create(log *model.RunLog) error {
	result, err := r.db.Exec(
		`INSERT INTO run_logs (companies_checked, jobs_found, new_jobs, notifications_sent, duration_ms, status, error_message) 
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		log.CompaniesChecked, log.JobsFound, log.NewJobs, log.NotificationsSent, log.DurationMs, log.Status, log.ErrorMessage,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	log.ID = id
	return nil
}

// GetRecent returns the most recent run logs.
func (r *RunLogRepository) GetRecent(limit int) ([]*model.RunLog, error) {
	rows, err := r.db.Query(
		`SELECT id, run_at, companies_checked, jobs_found, new_jobs, notifications_sent, duration_ms, status, error_message 
		 FROM run_logs ORDER BY run_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*model.RunLog
	for rows.Next() {
		l := &model.RunLog{}
		var errMsg sql.NullString
		err := rows.Scan(&l.ID, &l.RunAt, &l.CompaniesChecked, &l.JobsFound, &l.NewJobs, &l.NotificationsSent, &l.DurationMs, &l.Status, &errMsg)
		if err != nil {
			return nil, err
		}
		if errMsg.Valid {
			l.ErrorMessage = errMsg.String
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// GetStats returns aggregated statistics.
func (r *RunLogRepository) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total runs
	var totalRuns int
	r.db.QueryRow(`SELECT COUNT(*) FROM run_logs`).Scan(&totalRuns)
	stats["total_runs"] = totalRuns

	// Successful runs
	var successRuns int
	r.db.QueryRow(`SELECT COUNT(*) FROM run_logs WHERE status = 'success'`).Scan(&successRuns)
	stats["successful_runs"] = successRuns

	// Total new jobs found
	var totalNewJobs int
	r.db.QueryRow(`SELECT COALESCE(SUM(new_jobs), 0) FROM run_logs`).Scan(&totalNewJobs)
	stats["total_new_jobs_found"] = totalNewJobs

	// Average duration
	var avgDuration float64
	r.db.QueryRow(`SELECT COALESCE(AVG(duration_ms), 0) FROM run_logs`).Scan(&avgDuration)
	stats["avg_duration_ms"] = avgDuration

	// Last run time
	var lastRun sql.NullString
	r.db.QueryRow(`SELECT MAX(run_at) FROM run_logs`).Scan(&lastRun)
	if lastRun.Valid {
		stats["last_run"] = lastRun.String
	}

	return stats, nil
}
