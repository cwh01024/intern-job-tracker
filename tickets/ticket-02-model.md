# Ticket 02: Job Model & Repository

## Context
With the database schema in place (Ticket 01), we need Go structs to represent jobs and a repository layer for CRUD operations. This enables other components to interact with stored job data.

## Goals
1. Define Job struct with proper tags
2. Implement JobRepository with CRUD methods
3. Write comprehensive unit tests

## Dependencies
- Ticket 01 (DB schema must exist)

## Acceptance Criteria
- [ ] `internal/model/job.go` defines Job struct
- [ ] `internal/repository/job_repo.go` implements JobRepository
- [ ] `internal/repository/job_repo_test.go` covers all methods
- [ ] Tests pass: `go test ./internal/repository/...`

## Technical Details

### Job Struct (internal/model/job.go)
```go
type Job struct {
    ID           int64     `json:"id"`
    Company      string    `json:"company"`
    Title        string    `json:"title"`
    URL          string    `json:"url"`
    Location     string    `json:"location"`
    DiscoveredAt time.Time `json:"discovered_at"`
    Notified     bool      `json:"notified"`
}
```

### Repository Interface
```go
type JobRepository interface {
    Create(job *Job) error
    GetByURL(url string) (*Job, error)
    GetUnnotified() ([]*Job, error)
    MarkNotified(id int64) error
    GetAll() ([]*Job, error)
}
```

### Files to Create
```
internal/model/job.go
internal/repository/job_repo.go
internal/repository/job_repo_test.go
```

## TDD Steps
1. Write test for Create + GetByURL
2. Implement to pass
3. Write test for GetUnnotified
4. Implement to pass
5. Write test for MarkNotified
6. Implement to pass

## Estimated Tokens
~400 tokens
