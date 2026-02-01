package notifier

import (
	"fmt"
	"os/exec"
	"strings"

	"intern-job-tracker/internal/model"
)

// CommandExecutor interface for running external commands.
type CommandExecutor interface {
	Execute(name string, args ...string) error
}

// RealCommandExecutor executes real OS commands.
type RealCommandExecutor struct{}

// Execute runs the command.
func (r *RealCommandExecutor) Execute(name string, args ...string) error {
	return exec.Command(name, args...).Run()
}

// IMessageNotifier sends notifications via macOS iMessage.
type IMessageNotifier struct {
	executor CommandExecutor
}

// NewIMessageNotifier creates a new notifier with the given executor.
func NewIMessageNotifier(executor CommandExecutor) *IMessageNotifier {
	return &IMessageNotifier{executor: executor}
}

// NewDefaultIMessageNotifier creates a notifier with real command execution.
func NewDefaultIMessageNotifier() *IMessageNotifier {
	return &IMessageNotifier{executor: &RealCommandExecutor{}}
}

// Send sends a message to the recipient via iMessage.
func (n *IMessageNotifier) Send(recipient, message string) error {
	// Escape special characters in message
	escapedMessage := escapeAppleScript(message)
	escapedRecipient := escapeAppleScript(recipient)

	script := fmt.Sprintf(`
tell application "Messages"
	set targetService to 1st service whose service type = iMessage
	set targetBuddy to buddy "%s" of targetService
	send "%s" to targetBuddy
end tell`, escapedRecipient, escapedMessage)

	return n.executor.Execute("osascript", "-e", script)
}

// NotifyJob sends a formatted job notification.
func (n *IMessageNotifier) NotifyJob(recipient string, job *model.Job) error {
	message := FormatJobMessage(job)
	return n.Send(recipient, message)
}

// FormatJobMessage formats a job into a notification message.
func FormatJobMessage(job *model.Job) string {
	var sb strings.Builder
	sb.WriteString("🚀 New Intern Position Found!\n\n")
	sb.WriteString(fmt.Sprintf("Company: %s\n", job.Company))
	sb.WriteString(fmt.Sprintf("Title: %s\n", job.Title))
	if job.Location != "" {
		sb.WriteString(fmt.Sprintf("Location: %s\n", job.Location))
	}
	sb.WriteString(fmt.Sprintf("Apply: %s\n", job.URL))
	return sb.String()
}

// escapeAppleScript escapes special characters for AppleScript strings.
func escapeAppleScript(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
