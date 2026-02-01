package repository

import (
	"database/sql"
	"time"

	"intern-job-tracker/internal/model"
)

// CompanyRepository handles database operations for companies.
type CompanyRepository struct {
	db *sql.DB
}

// NewCompanyRepository creates a new CompanyRepository.
func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// GetAll returns all companies.
func (r *CompanyRepository) GetAll() ([]*model.Company, error) {
	rows, err := r.db.Query(`SELECT id, name, career_url, search_term, enabled, created_at FROM companies ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*model.Company
	for rows.Next() {
		c := &model.Company{}
		err := rows.Scan(&c.ID, &c.Name, &c.CareerURL, &c.SearchTerm, &c.Enabled, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, rows.Err()
}

// GetEnabled returns only enabled companies.
func (r *CompanyRepository) GetEnabled() ([]*model.Company, error) {
	rows, err := r.db.Query(`SELECT id, name, career_url, search_term, enabled, created_at FROM companies WHERE enabled = TRUE ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*model.Company
	for rows.Next() {
		c := &model.Company{}
		err := rows.Scan(&c.ID, &c.Name, &c.CareerURL, &c.SearchTerm, &c.Enabled, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, rows.Err()
}

// Create adds a new company.
func (r *CompanyRepository) Create(c *model.Company) error {
	result, err := r.db.Exec(
		`INSERT INTO companies (name, career_url, search_term, enabled) VALUES (?, ?, ?, ?)`,
		c.Name, c.CareerURL, c.SearchTerm, c.Enabled,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	c.ID = id
	c.CreatedAt = time.Now()
	return nil
}

// Update modifies an existing company.
func (r *CompanyRepository) Update(c *model.Company) error {
	_, err := r.db.Exec(
		`UPDATE companies SET name = ?, career_url = ?, search_term = ?, enabled = ? WHERE id = ?`,
		c.Name, c.CareerURL, c.SearchTerm, c.Enabled, c.ID,
	)
	return err
}

// Delete removes a company.
func (r *CompanyRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM companies WHERE id = ?`, id)
	return err
}

// GetByID retrieves a company by ID.
func (r *CompanyRepository) GetByID(id int64) (*model.Company, error) {
	c := &model.Company{}
	err := r.db.QueryRow(
		`SELECT id, name, career_url, search_term, enabled, created_at FROM companies WHERE id = ?`,
		id,
	).Scan(&c.ID, &c.Name, &c.CareerURL, &c.SearchTerm, &c.Enabled, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}
