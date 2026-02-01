# Ticket 04: Job Scraper

## Context
The scraper visits career pages of target companies and extracts intern job listings. It uses HTTP requests and simple HTML parsing to find relevant positions.

## Goals
1. Define company career page configurations
2. Implement HTTP client with retry logic
3. Parse HTML to extract job listings
4. Support mock HTTP for testing

## Dependencies
- Ticket 02 (Job model)

## Acceptance Criteria
- [ ] `internal/scraper/config.go` defines company URLs
- [ ] `internal/scraper/scraper.go` fetches and parses pages
- [ ] `internal/scraper/scraper_test.go` uses mock HTTP
- [ ] Scraper handles timeouts/errors gracefully

## Technical Details

### Company Configs
```go
var Companies = []CompanyConfig{
    {
        Name:    "Google",
        URL:     "https://www.google.com/about/careers/applications/jobs/results?q=intern&location=United+States",
        Pattern: `intern|internship`,
    },
    {
        Name:    "Amazon",
        URL:     "https://www.amazon.jobs/en/search?base_query=intern+software",
        Pattern: `intern|internship`,
    },
    // Uber, DoorDash...
}
```

### Scraper Interface
```go
type Scraper interface {
    ScrapeAll() ([]*model.Job, error)
    ScrapeCompany(config CompanyConfig) ([]*model.Job, error)
}
```

### Files to Create
```
internal/scraper/config.go
internal/scraper/scraper.go
internal/scraper/scraper_test.go
```

## TDD Steps
1. Write test with mock HTML response
2. Implement basic HTTP fetching
3. Add HTML parsing logic
4. Test error handling

## Estimated Tokens
~500 tokens (more parsing logic)
