# Ticket 06: HTTP API

## Context
The HTTP API exposes job data and configuration to the frontend dashboard. It provides REST endpoints for viewing jobs, notifications, and updating settings.

## Goals
1. Implement REST API using chi router
2. Create endpoints for jobs, notifications, config
3. Add manual refresh trigger endpoint
4. Write API tests

## Dependencies
- Ticket 02 (Repository)

## Acceptance Criteria
- [ ] `internal/api/handler.go` implements all endpoints
- [ ] `internal/api/handler_test.go` covers endpoints
- [ ] API returns proper JSON responses
- [ ] Error handling returns appropriate status codes

## Technical Details

### Endpoints
| Method | Path | Description |
|--------|------|-------------|
| GET | /api/jobs | List all jobs |
| GET | /api/jobs/:id | Get single job |
| GET | /api/notifications | Notification history |
| GET | /api/config | Get current config |
| PUT | /api/config | Update config |
| POST | /api/refresh | Trigger manual check |

### Handler Setup
```go
func NewRouter(repo repository.JobRepository, scheduler Scheduler) *chi.Mux {
    r := chi.NewRouter()
    r.Get("/api/jobs", listJobs)
    r.Get("/api/jobs/{id}", getJob)
    // ...
    return r
}
```

### Files to Create
```
internal/api/handler.go
internal/api/handler_test.go
```

## TDD Steps
1. Write test for GET /api/jobs
2. Implement handler
3. Repeat for each endpoint

## Estimated Tokens
~450 tokens
