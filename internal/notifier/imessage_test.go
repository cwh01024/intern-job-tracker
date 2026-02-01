package notifier

import (
	"strings"
	"testing"

	"intern-job-tracker/internal/model"
)

// MockCommandExecutor records executed commands for testing.
type MockCommandExecutor struct {
	ExecutedCommands []string
	ShouldFail       bool
}

func (m *MockCommandExecutor) Execute(name string, args ...string) error {
	cmd := name + " " + strings.Join(args, " ")
	m.ExecutedCommands = append(m.ExecutedCommands, cmd)
	if m.ShouldFail {
		return &mockError{"mock execution failed"}
	}
	return nil
}

type mockError struct{ msg string }

func (e *mockError) Error() string { return e.msg }

func TestIMessageNotifier_Send(t *testing.T) {
	mock := &MockCommandExecutor{}
	notifier := NewIMessageNotifier(mock)

	err := notifier.Send("+1234567890", "Test message")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.ExecutedCommands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(mock.ExecutedCommands))
	}

	cmd := mock.ExecutedCommands[0]
	if !strings.Contains(cmd, "osascript") {
		t.Error("expected osascript command")
	}
	if !strings.Contains(cmd, "+1234567890") {
		t.Error("expected recipient in command")
	}
	if !strings.Contains(cmd, "Test message") {
		t.Error("expected message in command")
	}
}

func TestIMessageNotifier_NotifyJob(t *testing.T) {
	mock := &MockCommandExecutor{}
	notifier := NewIMessageNotifier(mock)

	job := &model.Job{
		Company:  "Google",
		Title:    "Software Engineering Intern, Summer 2026",
		URL:      "https://careers.google.com/jobs/123",
		Location: "Mountain View, CA",
	}

	err := notifier.NotifyJob("+1234567890", job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.ExecutedCommands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(mock.ExecutedCommands))
	}

	cmd := mock.ExecutedCommands[0]
	if !strings.Contains(cmd, "Google") {
		t.Error("expected company name in message")
	}
	if !strings.Contains(cmd, "Software Engineering Intern") {
		t.Error("expected job title in message")
	}
	if !strings.Contains(cmd, "careers.google.com") {
		t.Error("expected URL in message")
	}
}

func TestIMessageNotifier_FormatMessage(t *testing.T) {
	job := &model.Job{
		Company:  "Amazon",
		Title:    "SDE Intern",
		URL:      "https://amazon.jobs/123",
		Location: "Seattle, WA",
	}

	msg := FormatJobMessage(job)

	expectations := []string{
		"New Intern Position",
		"Amazon",
		"SDE Intern",
		"Seattle, WA",
		"amazon.jobs",
	}

	for _, exp := range expectations {
		if !strings.Contains(msg, exp) {
			t.Errorf("expected message to contain %q", exp)
		}
	}
}

func TestIMessageNotifier_SendError(t *testing.T) {
	mock := &MockCommandExecutor{ShouldFail: true}
	notifier := NewIMessageNotifier(mock)

	err := notifier.Send("+1234567890", "Test")
	if err == nil {
		t.Error("expected error when command fails")
	}
}
