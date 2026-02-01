# Ticket 05: Scheduler Service

## Context
The scheduler runs the job checking process on a daily schedule. It orchestrates the scraper, compares with existing jobs, and triggers notifications for new findings.

## Goals
1. Setup cron-based scheduler
2. Implement scrape-compare-notify workflow
3. Make scheduler configurable
4. Write integration tests with mocks

## Dependencies
- Ticket 03 (Notifier)
- Ticket 04 (Scraper)

## Acceptance Criteria
- [ ] `internal/scheduler/scheduler.go` implements Scheduler
- [ ] Scheduler runs at configurable time (default 9 AM)
- [ ] Only notifies for NEW jobs not in DB
- [ ] `internal/scheduler/scheduler_test.go` covers workflow
- [ ] Logs output for debugging

## Technical Details

### Scheduler Interface
```go
type Scheduler interface {
    Start() error
    Stop()
    RunNow() error  // Manual trigger
}
```

### Workflow
```go
func (s *Scheduler) check() {
    jobs := s.scraper.ScrapeAll()
    for _, job := range jobs {
        existing := s.repo.GetByURL(job.URL)
        if existing == nil {
            s.repo.Create(job)
            s.notifier.NotifyJob(s.recipient, job)
            s.repo.MarkNotified(job.ID)
        }
    }
}
```

### Cron Expression
- Daily at 9 AM: `0 9 * * *`

### Files to Create
```
internal/scheduler/scheduler.go
internal/scheduler/scheduler_test.go
```

## TDD Steps
1. Write test with mock scraper/notifier
2. Implement basic scheduler
3. Test duplicate detection
4. Test cron scheduling

## Estimated Tokens
~400 tokens
