package scheduler

import (
	"testing"
	"time"

	"intern-job-tracker/internal/model"
	"intern-job-tracker/internal/scraper"
)

// MockScraper for testing
type MockScraper struct {
	Jobs []*model.Job
	Err  error
}

func (m *MockScraper) ScrapeAll() ([]*model.Job, error) {
	return m.Jobs, m.Err
}

func (m *MockScraper) ScrapeCompany(config scraper.CompanyConfig) ([]*model.Job, error) {
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

func (m *MockRepository) MarkNotified(id int64) error {
	m.Notified[id] = true
	return nil
}

// MockCompanyRepository for testing
type MockCompanyRepository struct {
	Companies []*model.Company
}

func (m *MockCompanyRepository) GetEnabled() ([]*model.Company, error) {
	return m.Companies, nil
}

// MockRunLogRepository for testing
type MockRunLogRepository struct {
	Logs []*model.RunLog
}

func (m *MockRunLogRepository) Create(log *model.RunLog) error {
	m.Logs = append(m.Logs, log)
	return nil
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
	companyRepo := &MockCompanyRepository{}
	runLogRepo := &MockRunLogRepository{}
	scr := &MockScraper{
		Jobs: []*model.Job{
			{Company: "Google", Title: "Intern 1", URL: "https://google.com/1"},
			{Company: "Amazon", Title: "Intern 2", URL: "https://amazon.com/2"},
		},
	}
	notifier := &MockNotifier{}

	sched := New(repo, companyRepo, runLogRepo, scr, notifier, "+1234567890")
	err := sched.RunNow()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have created 2 jobs
	if len(repo.Jobs) != 2 {
		t.Errorf("expected 2 jobs created, got %d", len(repo.Jobs))
	}

	// Should have sent 2 notifications + 1 run log
	if len(runLogRepo.Logs) != 1 {
		t.Errorf("expected 1 run log, got %d", len(runLogRepo.Logs))
	}
}

func TestScheduler_RunNow_WithCompanies(t *testing.T) {
	repo := NewMockRepository()
	companyRepo := &MockCompanyRepository{
		Companies: []*model.Company{
			{ID: 1, Name: "Google", CareerURL: "https://google.com/careers", SearchTerm: "intern"},
		},
	}
	runLogRepo := &MockRunLogRepository{}
	scr := &MockScraper{
		Jobs: []*model.Job{
			{Company: "Google", Title: "SDE Intern", URL: "https://google.com/job/1"},
		},
	}
	notifier := &MockNotifier{}

	sched := New(repo, companyRepo, runLogRepo, scr, notifier, "+1234567890")
	err := sched.RunNow()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.Jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(repo.Jobs))
	}
}

func TestScheduler_StartStop(t *testing.T) {
	repo := NewMockRepository()
	companyRepo := &MockCompanyRepository{}
	runLogRepo := &MockRunLogRepository{}
	scr := &MockScraper{Jobs: []*model.Job{}}
	notifier := &MockNotifier{}

	sched := New(repo, companyRepo, runLogRepo, scr, notifier, "+1234567890")

	err := sched.Start()
	if err != nil {
		t.Fatalf("failed to start scheduler: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	sched.Stop()
}
