package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"intern-job-tracker/internal/model"
	"intern-job-tracker/internal/scraper"

	"github.com/robfig/cron/v3"
)

// Scraper interface for job scraping.
type Scraper interface {
	ScrapeAll() ([]*model.Job, error)
	ScrapeCompany(config scraper.CompanyConfig) ([]*model.Job, error)
}

// Repository interface for job storage.
type Repository interface {
	Create(job *model.Job) error
	GetByURL(url string) (*model.Job, error)
	MarkNotified(id int64) error
}

// CompanyRepository interface for company storage.
type CompanyRepository interface {
	GetEnabled() ([]*model.Company, error)
}

// RunLogRepository interface for run logs.
type RunLogRepository interface {
	Create(log *model.RunLog) error
}

// Notifier interface for sending notifications.
type Notifier interface {
	NotifyJob(recipient string, job *model.Job) error
	Send(recipient string, message string) error
}

// Scheduler manages the job checking schedule.
type Scheduler struct {
	repo        Repository
	companyRepo CompanyRepository
	runLogRepo  RunLogRepository
	scraper     Scraper
	notifier    Notifier
	recipient   string
	cron        *cron.Cron
	mu          sync.Mutex
}

// New creates a new Scheduler.
func New(repo Repository, companyRepo CompanyRepository, runLogRepo RunLogRepository, scr Scraper, notifier Notifier, recipient string) *Scheduler {
	return &Scheduler{
		repo:        repo,
		companyRepo: companyRepo,
		runLogRepo:  runLogRepo,
		scraper:     scr,
		notifier:    notifier,
		recipient:   recipient,
	}
}

// Start begins the scheduled job checking. Default: daily at 9 AM.
func (s *Scheduler) Start() error {
	return s.StartWithSchedule("0 9 * * *")
}

// StartWithSchedule begins job checking with a custom cron schedule.
func (s *Scheduler) StartWithSchedule(schedule string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cron = cron.New()
	_, err := s.cron.AddFunc(schedule, func() {
		if err := s.RunNow(); err != nil {
			log.Printf("❌ Error during scheduled check: %v", err)
		}
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	log.Printf("⏰ Scheduler started with schedule: %s", schedule)
	return nil
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cron != nil {
		s.cron.Stop()
		log.Println("⏹️  Scheduler stopped")
	}
}

// RunNow performs an immediate job check.
func (s *Scheduler) RunNow() error {
	startTime := time.Now()
	runLog := &model.RunLog{
		RunAt:  startTime,
		Status: "success",
	}

	log.Println("═══════════════════════════════════════════")
	log.Println("🔍 Starting job check...")
	log.Printf("⏰ Time: %s", startTime.Format("2006-01-02 15:04:05"))

	// Get enabled companies
	var companies []*model.Company
	var err error
	if s.companyRepo != nil {
		companies, err = s.companyRepo.GetEnabled()
		if err != nil {
			log.Printf("❌ Error getting companies: %v", err)
			runLog.Status = "error"
			runLog.ErrorMessage = err.Error()
			s.saveRunLog(runLog, startTime)
			return err
		}
	}

	if len(companies) == 0 {
		log.Println("⚠️  No companies configured, using defaults")
		// Fall back to default scraper
		return s.runWithDefaultScraper(runLog, startTime)
	}

	runLog.CompaniesChecked = len(companies)
	log.Printf("📋 Companies to check: %d", len(companies))
	log.Println("───────────────────────────────────────────")

	totalJobs := 0
	newCount := 0
	notificationsSent := 0

	for _, company := range companies {
		log.Printf("🏢 Checking: %s", company.Name)

		config := scraper.CompanyConfig{
			Name:       company.Name,
			CareerURL:  company.CareerURL,
			SearchTerm: company.SearchTerm,
		}

		jobs, err := s.scraper.ScrapeCompany(config)
		if err != nil {
			log.Printf("   ❌ Error scraping %s: %v", company.Name, err)
			continue
		}

		log.Printf("   📄 Found %d job listings", len(jobs))
		totalJobs += len(jobs)

		for _, job := range jobs {
			existing, err := s.repo.GetByURL(job.URL)
			if err != nil {
				log.Printf("   ❌ Error checking job: %v", err)
				continue
			}

			if existing != nil {
				continue
			}

			// New job found!
			log.Printf("   ✨ NEW: %s", job.Title)
			if err := s.repo.Create(job); err != nil {
				log.Printf("   ❌ Error saving job: %v", err)
				continue
			}

			if err := s.notifier.NotifyJob(s.recipient, job); err != nil {
				log.Printf("   ❌ Error sending notification: %v", err)
			} else {
				s.repo.MarkNotified(job.ID)
				newCount++
				notificationsSent++
			}
		}
	}

	runLog.JobsFound = totalJobs
	runLog.NewJobs = newCount
	runLog.NotificationsSent = notificationsSent

	log.Println("───────────────────────────────────────────")
	log.Printf("📊 Summary:")
	log.Printf("   • Companies checked: %d", len(companies))
	log.Printf("   • Total jobs found: %d", totalJobs)
	log.Printf("   • New positions: %d", newCount)
	log.Printf("   • Notifications sent: %d", notificationsSent)

	// Send summary notification
	if newCount == 0 {
		msg := fmt.Sprintf("📋 Intern Job Tracker Update\n\n✅ Checked %d companies\n📄 Found %d job listings\n🆕 No new positions found\n\nTracking: %s",
			len(companies), totalJobs, s.getCompanyNames(companies))
		if err := s.notifier.Send(s.recipient, msg); err != nil {
			log.Printf("   ❌ Error sending summary: %v", err)
		} else {
			notificationsSent++
		}
	}

	s.saveRunLog(runLog, startTime)

	log.Println("═══════════════════════════════════════════")
	return nil
}

func (s *Scheduler) runWithDefaultScraper(runLog *model.RunLog, startTime time.Time) error {
	jobs, err := s.scraper.ScrapeAll()
	if err != nil {
		return err
	}

	runLog.JobsFound = len(jobs)
	newCount := 0

	for _, job := range jobs {
		existing, _ := s.repo.GetByURL(job.URL)
		if existing != nil {
			continue
		}

		s.repo.Create(job)
		if err := s.notifier.NotifyJob(s.recipient, job); err == nil {
			s.repo.MarkNotified(job.ID)
			newCount++
		}
	}

	runLog.NewJobs = newCount
	runLog.NotificationsSent = newCount

	if newCount == 0 {
		s.notifier.Send(s.recipient, "📋 No new intern positions found.")
	}

	s.saveRunLog(runLog, startTime)
	return nil
}

func (s *Scheduler) saveRunLog(runLog *model.RunLog, startTime time.Time) {
	runLog.DurationMs = time.Since(startTime).Milliseconds()
	if s.runLogRepo != nil {
		if err := s.runLogRepo.Create(runLog); err != nil {
			log.Printf("❌ Error saving run log: %v", err)
		}
	}
	log.Printf("⏱️  Duration: %dms", runLog.DurationMs)
}

func (s *Scheduler) getCompanyNames(companies []*model.Company) string {
	if len(companies) == 0 {
		return ""
	}
	names := companies[0].Name
	for i := 1; i < len(companies) && i < 4; i++ {
		names += ", " + companies[i].Name
	}
	if len(companies) > 4 {
		names += fmt.Sprintf(" +%d more", len(companies)-4)
	}
	return names
}

// SetRecipient updates the notification recipient.
func (s *Scheduler) SetRecipient(recipient string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recipient = recipient
}
