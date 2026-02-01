# Ticket 08: Integration & GitHub CI

## Context
Final ticket to wire everything together, add end-to-end tests, setup GitHub Actions, and write documentation.

## Goals
1. Complete main.go to wire all components
2. Add GitHub Actions CI workflow
3. Write README with setup instructions
4. Create end-to-end test

## Dependencies
- All previous tickets

## Acceptance Criteria
- [ ] `cmd/server/main.go` starts full application
- [ ] `.github/workflows/ci.yml` runs tests on push
- [ ] `README.md` has complete documentation
- [ ] E2E test verifies full workflow
- [ ] All tests pass in CI

## Technical Details

### main.go Structure
```go
func main() {
    db := db.New("jobs.db")
    repo := repository.NewJobRepository(db)
    notifier := notifier.NewIMessageNotifier()
    scraper := scraper.New()
    scheduler := scheduler.New(repo, scraper, notifier)
    
    scheduler.Start()
    
    router := api.NewRouter(repo, scheduler)
    http.ListenAndServe(":8080", router)
}
```

### GitHub Actions
```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - run: go test ./... -v
```

### Files to Create/Update
```
cmd/server/main.go (complete)
.github/workflows/ci.yml
README.md
internal/e2e_test.go
```

## TDD Steps
1. Write E2E test for full workflow
2. Complete main.go wiring
3. Push to GitHub, verify CI passes

## Estimated Tokens
~450 tokens
