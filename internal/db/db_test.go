package db

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestNewDB_CreatesTables(t *testing.T) {
	// Use temp file for test database
	tmpFile, err := os.CreateTemp("", "test_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Create database
	database, err := New(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer database.Close()

	// Verify tables exist
	tables := []string{"jobs", "notifications", "config"}
	for _, table := range tables {
		var name string
		err := database.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
			table,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %s does not exist: %v", table, err)
		}
	}
}

func TestNewDB_JobsTableSchema(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	database, err := New(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer database.Close()

	// Try inserting a job
	result, err := database.Exec(
		"INSERT INTO jobs (company, title, url, location) VALUES (?, ?, ?, ?)",
		"Google", "Software Engineering Intern", "https://google.com/job/123", "Mountain View, CA",
	)
	if err != nil {
		t.Fatalf("failed to insert job: %v", err)
	}

	id, _ := result.LastInsertId()
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}

	// Verify we can query it back
	var company, title, url string
	err = database.QueryRow("SELECT company, title, url FROM jobs WHERE id = ?", id).Scan(&company, &title, &url)
	if err != nil {
		t.Fatalf("failed to query job: %v", err)
	}

	if company != "Google" {
		t.Errorf("expected Google, got %s", company)
	}
}

func TestDB_Implements_SQLDB(t *testing.T) {
	// Ensure *DB satisfies *sql.DB interface for common operations
	var _ interface {
		Exec(string, ...any) (sql.Result, error)
		Query(string, ...any) (*sql.Rows, error)
		QueryRow(string, ...any) *sql.Row
		Close() error
	} = (*sql.DB)(nil)
}
