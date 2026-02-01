package repository

import (
	"database/sql"
	"time"

	"intern-job-tracker/internal/model"
)

// JobRepository handles database operations for jobs.
type JobRepository struct {
	db *sql.DB
}

// NewJobRepository creates a new JobRepository.
func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Create inserts a new job into the database.
func (r *JobRepository) Create(job *model.Job) error {
	result, err := r.db.Exec(
		`INSERT INTO jobs (company, title, url, location, notified) VALUES (?, ?, ?, ?, ?)`,
		job.Company, job.Title, job.URL, job.Location, false,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	job.ID = id
	job.DiscoveredAt = time.Now()
	job.Notified = false
	return nil
}

// GetByURL retrieves a job by its URL. Returns nil if not found.
func (r *JobRepository) GetByURL(url string) (*model.Job, error) {
	job := &model.Job{}
	err := r.db.QueryRow(
		`SELECT id, company, title, url, location, discovered_at, notified FROM jobs WHERE url = ?`,
		url,
	).Scan(&job.ID, &job.Company, &job.Title, &job.URL, &job.Location, &job.DiscoveredAt, &job.Notified)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return job, nil
}

// GetUnnotified returns all jobs that haven't been notified yet.
func (r *JobRepository) GetUnnotified() ([]*model.Job, error) {
	rows, err := r.db.Query(
		`SELECT id, company, title, url, location, discovered_at, notified FROM jobs WHERE notified = FALSE`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJobs(rows)
}

// MarkNotified marks a job as notified.
func (r *JobRepository) MarkNotified(id int64) error {
	_, err := r.db.Exec(`UPDATE jobs SET notified = TRUE WHERE id = ?`, id)
	return err
}

// GetAll returns all jobs.
func (r *JobRepository) GetAll() ([]*model.Job, error) {
	rows, err := r.db.Query(
		`SELECT id, company, title, url, location, discovered_at, notified FROM jobs ORDER BY discovered_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJobs(rows)
}

// GetByID retrieves a job by its ID.
func (r *JobRepository) GetByID(id int64) (*model.Job, error) {
	job := &model.Job{}
	err := r.db.QueryRow(
		`SELECT id, company, title, url, location, discovered_at, notified FROM jobs WHERE id = ?`,
		id,
	).Scan(&job.ID, &job.Company, &job.Title, &job.URL, &job.Location, &job.DiscoveredAt, &job.Notified)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return job, nil
}

func scanJobs(rows *sql.Rows) ([]*model.Job, error) {
	var jobs []*model.Job
	for rows.Next() {
		job := &model.Job{}
		err := rows.Scan(&job.ID, &job.Company, &job.Title, &job.URL, &job.Location, &job.DiscoveredAt, &job.Notified)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}
