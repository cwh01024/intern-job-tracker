package scheduler

import (
	"testing"
	"time"

	"intern-job-tracker/internal/model"
)

// MockScraper for testing
type MockScraper struct {
	Jobs []*model.Job
	Err  error
}

func (m *MockScraper) ScrapeAll() ([]*model.Job, error) {
	return m.Jobs, m.Err
}

// MockRepository for testing
type MockRepository struct {
	Jobs       map[string]*model.Job
	Notified   map[int64]bool
	CreateErr  error
	CreatedIDs int64
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Jobs:     make(map[string]*model.Job),
		Notified: make(map[int64]bool),
	}
}

func (m *MockRepository) Create(job *model.Job) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	m.CreatedIDs++
	job.ID = m.CreatedIDs
	m.Jobs[job.URL] = job
	return nil
}

func (m *MockRepository) GetByURL(url string) (*model.Job, error) {
	return m.Jobs[url], nil
}

func (m *MockRepository) GetUnnotified() ([]*model.Job, error) {
	var jobs []*model.Job
	for _, j := range m.Jobs {
		if !m.Notified[j.ID] {
			jobs = append(jobs, j)
		}
	}
	return jobs, nil
}

func (m *MockRepository) MarkNotified(id int64) error {
	m.Notified[id] = true
	return nil
}

func (m *MockRepository) GetAll() ([]*model.Job, error) {
	var jobs []*model.Job
	for _, j := range m.Jobs {
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// MockNotifier for testing
type MockNotifier struct {
	SentMessages []string
	Err          error
}

func (m *MockNotifier) Send(recipient, message string) error {
	m.SentMessages = append(m.SentMessages, message)
	return m.Err
}

func (m *MockNotifier) NotifyJob(recipient string, job *model.Job) error {
	m.SentMessages = append(m.SentMessages, job.Title)
	return m.Err
}

func TestScheduler_RunNow_FindsNewJobs(t *testing.T) {
	repo := NewMockRepository()
	scraper := &MockScraper{
		Jobs: []*model.Job{
			{Company: "Google", Title: "Intern 1", URL: "https://google.com/1"},
			{Company: "Amazon", Title: "Intern 2", URL: "https://amazon.com/2"},
		},
	}
	notifier := &MockNotifier{}

	sched := New(repo, scraper, notifier, "+1234567890")
	err := sched.RunNow()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have created 2 jobs
	if len(repo.Jobs) != 2 {
		t.Errorf("expected 2 jobs created, got %d", len(repo.Jobs))
	}

	// Should have sent 2 notifications
	if len(notifier.SentMessages) != 2 {
		t.Errorf("expected 2 notifications, got %d", len(notifier.SentMessages))
	}
}

func TestScheduler_RunNow_SkipsDuplicates(t *testing.T) {
	repo := NewMockRepository()
	// Pre-populate with existing job
	repo.Jobs["https://google.com/1"] = &model.Job{ID: 1, URL: "https://google.com/1"}

	scraper := &MockScraper{
		Jobs: []*model.Job{
			{Company: "Google", Title: "Intern 1", URL: "https://google.com/1"},
			{Company: "Amazon", Title: "Intern 2", URL: "https://amazon.com/2"},
		},
	}
	notifier := &MockNotifier{}

	sched := New(repo, scraper, notifier, "+1234567890")
	err := sched.RunNow()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should only have 2 jobs (1 existing + 1 new)
	if len(repo.Jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(repo.Jobs))
	}

	// Should only send 1 notification (for new job)
	if len(notifier.SentMessages) != 1 {
		t.Errorf("expected 1 notification, got %d", len(notifier.SentMessages))
	}
}

func TestScheduler_RunNow_MarksNotified(t *testing.T) {
	repo := NewMockRepository()
	scraper := &MockScraper{
		Jobs: []*model.Job{
			{Company: "Uber", Title: "Intern", URL: "https://uber.com/1"},
		},
	}
	notifier := &MockNotifier{}

	sched := New(repo, scraper, notifier, "+1234567890")
	sched.RunNow()

	// Job should be marked as notified
	if !repo.Notified[1] {
		t.Error("expected job to be marked as notified")
	}
}

func TestScheduler_StartStop(t *testing.T) {
	repo := NewMockRepository()
	scraper := &MockScraper{Jobs: []*model.Job{}}
	notifier := &MockNotifier{}

	sched := New(repo, scraper, notifier, "+1234567890")

	err := sched.Start()
	if err != nil {
		t.Fatalf("failed to start scheduler: %v", err)
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	sched.Stop()
	// Should not panic
}
