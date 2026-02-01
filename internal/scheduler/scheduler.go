package scheduler

import (
	"log"
	"sync"

	"intern-job-tracker/internal/model"

	"github.com/robfig/cron/v3"
)

// Scraper interface for job scraping.
type Scraper interface {
	ScrapeAll() ([]*model.Job, error)
}

// Repository interface for job storage.
type Repository interface {
	Create(job *model.Job) error
	GetByURL(url string) (*model.Job, error)
	MarkNotified(id int64) error
}

// Notifier interface for sending notifications.
type Notifier interface {
	NotifyJob(recipient string, job *model.Job) error
	Send(recipient string, message string) error
}

// Scheduler manages the job checking schedule.
type Scheduler struct {
	repo      Repository
	scraper   Scraper
	notifier  Notifier
	recipient string
	cron      *cron.Cron
	mu        sync.Mutex
}

// New creates a new Scheduler.
func New(repo Repository, scraper Scraper, notifier Notifier, recipient string) *Scheduler {
	return &Scheduler{
		repo:      repo,
		scraper:   scraper,
		notifier:  notifier,
		recipient: recipient,
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
			log.Printf("Error during scheduled check: %v", err)
		}
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	log.Printf("Scheduler started with schedule: %s", schedule)
	return nil
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cron != nil {
		s.cron.Stop()
		log.Println("Scheduler stopped")
	}
}

// RunNow performs an immediate job check.
func (s *Scheduler) RunNow() error {
	log.Println("Starting job check...")

	jobs, err := s.scraper.ScrapeAll()
	if err != nil {
		return err
	}

	newCount := 0
	for _, job := range jobs {
		// Check if job already exists
		existing, err := s.repo.GetByURL(job.URL)
		if err != nil {
			log.Printf("Error checking job %s: %v", job.URL, err)
			continue
		}

		if existing != nil {
			// Job already exists, skip
			continue
		}

		// New job found!
		if err := s.repo.Create(job); err != nil {
			log.Printf("Error creating job: %v", err)
			continue
		}

		// Send notification
		if err := s.notifier.NotifyJob(s.recipient, job); err != nil {
			log.Printf("Error sending notification: %v", err)
		} else {
			// Mark as notified
			s.repo.MarkNotified(job.ID)
			newCount++
		}
	}

	log.Printf("Job check complete. Found %d new positions.", newCount)

	// Send notification with status
	if newCount == 0 {
		msg := "📋 Intern Job Tracker Update\n\nNo new positions found this check.\nCurrently tracking: Google, Amazon, Uber, DoorDash"
		if err := s.notifier.Send(s.recipient, msg); err != nil {
			log.Printf("Error sending status notification: %v", err)
		}
	}

	return nil
}

// SetRecipient updates the notification recipient.
func (s *Scheduler) SetRecipient(recipient string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recipient = recipient
}
