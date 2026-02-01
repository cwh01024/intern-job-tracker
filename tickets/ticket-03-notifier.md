# Ticket 03: iMessage Notifier

## Context
The core notification feature uses macOS's native Messages app via AppleScript. This component formats job alerts and sends them through the `osascript` command.

## Goals
1. Create iMessage notifier using osascript
2. Format job data into readable messages
3. Write testable code with mock support

## Dependencies
- Ticket 01 (for notification logging)

## Acceptance Criteria
- [ ] `internal/notifier/imessage.go` implements Notifier interface
- [ ] Messages format job info clearly
- [ ] `internal/notifier/imessage_test.go` uses mock executor
- [ ] Manual test sends real iMessage

## Technical Details

### Notifier Interface
```go
type Notifier interface {
    Send(recipient string, message string) error
    NotifyJob(recipient string, job *model.Job) error
}
```

### AppleScript Command
```go
script := fmt.Sprintf(`
tell application "Messages"
    set targetService to 1st service whose service type = iMessage
    set targetBuddy to buddy "%s" of targetService
    send "%s" to targetBuddy
end tell`, recipient, message)

exec.Command("osascript", "-e", script).Run()
```

### Message Format
```
🚀 New Intern Position Found!

Company: Google
Title: Software Engineering Intern, Summer 2026
Location: Mountain View, CA
Apply: https://careers.google.com/...

Found at: 2026-02-01 09:00 AM
```

### Files to Create
```
internal/notifier/imessage.go
internal/notifier/imessage_test.go
```

## TDD Steps
1. Write test with mock command executor
2. Implement notifier with real osascript
3. Add message formatting test
4. Manual test with real phone number

## Estimated Tokens
~350 tokens
