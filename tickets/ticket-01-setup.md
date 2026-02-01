# Ticket 01: Project Setup & DB Schema

## Context
This is the foundation ticket for the Intern Job Tracker system. It establishes the Go project structure, dependencies, and database schema that all other components will build upon.

## Goals
1. Initialize Go module with required dependencies
2. Create SQLite database schema
3. Setup project directory structure
4. Write migration test to verify schema

## Dependencies
- None (first ticket)

## Acceptance Criteria
- [ ] `go.mod` and `go.sum` exist with dependencies
- [ ] `migrations/001_init.sql` creates tables
- [ ] `internal/db/db.go` has connection logic
- [ ] `internal/db/db_test.go` verifies schema creation
- [ ] All tests pass: `go test ./internal/db/...`

## Technical Details

### Dependencies to install
```bash
go mod init intern-job-tracker
go get modernc.org/sqlite
go get github.com/go-chi/chi/v5
go get github.com/robfig/cron/v3
```

### Schema (migrations/001_init.sql)
```sql
CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    location TEXT,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    notified BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id INTEGER REFERENCES jobs(id),
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT
);

CREATE TABLE IF NOT EXISTS config (
    key TEXT PRIMARY KEY,
    value TEXT
);
```

### Files to Create
```
cmd/server/main.go        (stub only)
internal/db/db.go
internal/db/db_test.go
migrations/001_init.sql
```

## TDD Steps
1. Write test that opens DB and runs migration
2. Implement db.go to make test pass
3. Verify with `go test -v`

## Estimated Tokens
~500 tokens (small, focused scope)
